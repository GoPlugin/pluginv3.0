// SPDX-License-Identifier: MIT
// solhint-disable-next-line one-contract-per-file
pragma solidity ^0.8.6;

import {ConfirmedOwner} from "../shared/access/ConfirmedOwner.sol";
import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {VRFConsumerBaseV2} from "./VRFConsumerBaseV2.sol";
import {PliTokenInterface} from "../shared/interfaces/PliTokenInterface.sol";
import {AggregatorV3Interface} from "../shared/interfaces/AggregatorV3Interface.sol";
import {VRFCoordinatorV2Interface} from "./interfaces/VRFCoordinatorV2Interface.sol";
import {VRFV2WrapperInterface} from "./interfaces/VRFV2WrapperInterface.sol";
import {VRFV2WrapperConsumerBase} from "./VRFV2WrapperConsumerBase.sol";
import {ChainSpecificUtil} from "../ChainSpecificUtil_v0_8_6.sol";

/**
 * @notice A wrapper for VRFCoordinatorV2 that provides an interface better suited to one-off
 * @notice requests for randomness.
 */
contract VRFV2Wrapper is ConfirmedOwner, TypeAndVersionInterface, VRFConsumerBaseV2, VRFV2WrapperInterface {
  event WrapperFulfillmentFailed(uint256 indexed requestId, address indexed consumer);

  // solhint-disable-next-line plugin-solidity/prefix-immutable-variables-with-i
  PliTokenInterface public immutable PLI;
  // solhint-disable-next-line plugin-solidity/prefix-immutable-variables-with-i
  AggregatorV3Interface public immutable PLI_ETH_FEED;
  // solhint-disable-next-line plugin-solidity/prefix-immutable-variables-with-i
  ExtendedVRFCoordinatorV2Interface public immutable COORDINATOR;
  // solhint-disable-next-line plugin-solidity/prefix-immutable-variables-with-i
  uint64 public immutable SUBSCRIPTION_ID;
  /// @dev this is the size of a VRF v2 fulfillment's calldata abi-encoded in bytes.
  /// @dev proofSize = 13 words = 13 * 256 = 3328 bits
  /// @dev commitmentSize = 5 words = 5 * 256 = 1280 bits
  /// @dev dataSize = proofSize + commitmentSize = 4608 bits
  /// @dev selector = 32 bits
  /// @dev total data size = 4608 bits + 32 bits = 4640 bits = 580 bytes
  uint32 public s_fulfillmentTxSizeBytes = 580;

  // 5k is plenty for an EXTCODESIZE call (2600) + warm CALL (100)
  // and some arithmetic operations.
  uint256 private constant GAS_FOR_CALL_EXACT_CHECK = 5_000;

  // lastRequestId is the request ID of the most recent VRF V2 request made by this wrapper. This
  // should only be relied on within the same transaction the request was made.
  uint256 public override lastRequestId;

  // Configuration fetched from VRFCoordinatorV2

  // s_configured tracks whether this contract has been configured. If not configured, randomness
  // requests cannot be made.
  bool public s_configured;

  // s_disabled disables the contract when true. When disabled, new VRF requests cannot be made
  // but existing ones can still be fulfilled.
  bool public s_disabled;

  // s_fallbackWeiPerUnitPli is the backup PLI exchange rate used when the PLI/NATIVE feed is
  // stale.
  int256 private s_fallbackWeiPerUnitPli;

  // s_stalenessSeconds is the number of seconds before we consider the feed price to be stale and
  // fallback to fallbackWeiPerUnitPli.
  uint32 private s_stalenessSeconds;

  // s_fulfillmentFlatFeePliPPM is the flat fee in millionths of PLI that VRFCoordinatorV2
  // charges.
  uint32 private s_fulfillmentFlatFeePliPPM;

  // Other configuration

  // s_wrapperGasOverhead reflects the gas overhead of the wrapper's fulfillRandomWords
  // function. The cost for this gas is passed to the user.
  uint32 private s_wrapperGasOverhead;

  // s_coordinatorGasOverhead reflects the gas overhead of the coordinator's fulfillRandomWords
  // function. The cost for this gas is billed to the subscription, and must therefor be included
  // in the pricing for wrapped requests. This includes the gas costs of proof verification and
  // payment calculation in the coordinator.
  uint32 private s_coordinatorGasOverhead;

  // s_wrapperPremiumPercentage is the premium ratio in percentage. For example, a value of 0
  // indicates no premium. A value of 15 indicates a 15 percent premium.
  uint8 private s_wrapperPremiumPercentage;

  // s_keyHash is the key hash to use when requesting randomness. Fees are paid based on current gas
  // fees, so this should be set to the highest gas lane on the network.
  bytes32 internal s_keyHash;

  // s_maxNumWords is the max number of words that can be requested in a single wrapped VRF request.
  uint8 internal s_maxNumWords;

  struct Callback {
    address callbackAddress;
    uint32 callbackGasLimit;
    uint256 requestGasPrice;
    int256 requestWeiPerUnitPli;
    uint256 juelsPaid;
  }
  mapping(uint256 => Callback) /* requestID */ /* callback */ public s_callbacks;

  constructor(
    address _pli,
    address _pliEthFeed,
    address _coordinator
  ) ConfirmedOwner(msg.sender) VRFConsumerBaseV2(_coordinator) {
    PLI = PliTokenInterface(_pli);
    PLI_ETH_FEED = AggregatorV3Interface(_pliEthFeed);
    COORDINATOR = ExtendedVRFCoordinatorV2Interface(_coordinator);

    // Create this wrapper's subscription and add itself as a consumer.
    uint64 subId = ExtendedVRFCoordinatorV2Interface(_coordinator).createSubscription();
    SUBSCRIPTION_ID = subId;
    ExtendedVRFCoordinatorV2Interface(_coordinator).addConsumer(subId, address(this));
  }

  /**
   * @notice setFulfillmentTxSize sets the size of the fulfillment transaction in bytes.
   * @param size is the size of the fulfillment transaction in bytes.
   */
  function setFulfillmentTxSize(uint32 size) external onlyOwner {
    s_fulfillmentTxSizeBytes = size;
  }

  /**
   * @notice setConfig configures VRFV2Wrapper.
   *
   * @dev Sets wrapper-specific configuration based on the given parameters, and fetches any needed
   * @dev VRFCoordinatorV2 configuration from the coordinator.
   *
   * @param _wrapperGasOverhead reflects the gas overhead of the wrapper's fulfillRandomWords
   *        function.
   *
   * @param _coordinatorGasOverhead reflects the gas overhead of the coordinator's
   *        fulfillRandomWords function.
   *
   * @param _wrapperPremiumPercentage is the premium ratio in percentage for wrapper requests.
   *
   * @param _keyHash to use for requesting randomness.
   */
  function setConfig(
    uint32 _wrapperGasOverhead,
    uint32 _coordinatorGasOverhead,
    uint8 _wrapperPremiumPercentage,
    bytes32 _keyHash,
    uint8 _maxNumWords
  ) external onlyOwner {
    s_wrapperGasOverhead = _wrapperGasOverhead;
    s_coordinatorGasOverhead = _coordinatorGasOverhead;
    s_wrapperPremiumPercentage = _wrapperPremiumPercentage;
    s_keyHash = _keyHash;
    s_maxNumWords = _maxNumWords;
    s_configured = true;

    // Get other configuration from coordinator
    (, , s_stalenessSeconds, ) = COORDINATOR.getConfig();
    s_fallbackWeiPerUnitPli = COORDINATOR.getFallbackWeiPerUnitPli();
    (s_fulfillmentFlatFeePliPPM, , , , , , , , ) = COORDINATOR.getFeeConfig();
  }

  /**
   * @notice getConfig returns the current VRFV2Wrapper configuration.
   *
   * @return fallbackWeiPerUnitPli is the backup PLI exchange rate used when the PLI/NATIVE feed
   *         is stale.
   *
   * @return stalenessSeconds is the number of seconds before we consider the feed price to be stale
   *         and fallback to fallbackWeiPerUnitPli.
   *
   * @return fulfillmentFlatFeePliPPM is the flat fee in millionths of PLI that VRFCoordinatorV2
   *         charges.
   *
   * @return wrapperGasOverhead reflects the gas overhead of the wrapper's fulfillRandomWords
   *         function. The cost for this gas is passed to the user.
   *
   * @return coordinatorGasOverhead reflects the gas overhead of the coordinator's
   *         fulfillRandomWords function.
   *
   * @return wrapperPremiumPercentage is the premium ratio in percentage. For example, a value of 0
   *         indicates no premium. A value of 15 indicates a 15 percent premium.
   *
   * @return keyHash is the key hash to use when requesting randomness. Fees are paid based on
   *         current gas fees, so this should be set to the highest gas lane on the network.
   *
   * @return maxNumWords is the max number of words that can be requested in a single wrapped VRF
   *         request.
   */
  function getConfig()
    external
    view
    returns (
      int256 fallbackWeiPerUnitPli,
      uint32 stalenessSeconds,
      uint32 fulfillmentFlatFeePliPPM,
      uint32 wrapperGasOverhead,
      uint32 coordinatorGasOverhead,
      uint8 wrapperPremiumPercentage,
      bytes32 keyHash,
      uint8 maxNumWords
    )
  {
    return (
      s_fallbackWeiPerUnitPli,
      s_stalenessSeconds,
      s_fulfillmentFlatFeePliPPM,
      s_wrapperGasOverhead,
      s_coordinatorGasOverhead,
      s_wrapperPremiumPercentage,
      s_keyHash,
      s_maxNumWords
    );
  }

  /**
   * @notice Calculates the price of a VRF request with the given callbackGasLimit at the current
   * @notice block.
   *
   * @dev This function relies on the transaction gas price which is not automatically set during
   * @dev simulation. To estimate the price at a specific gas price, use the estimatePrice function.
   *
   * @param _callbackGasLimit is the gas limit used to estimate the price.
   */
  function calculateRequestPrice(
    uint32 _callbackGasLimit
  ) external view override onlyConfiguredNotDisabled returns (uint256) {
    int256 weiPerUnitPli = _getFeedData();
    return _calculateRequestPrice(_callbackGasLimit, tx.gasprice, weiPerUnitPli);
  }

  /**
   * @notice Estimates the price of a VRF request with a specific gas limit and gas price.
   *
   * @dev This is a convenience function that can be called in simulation to better understand
   * @dev pricing.
   *
   * @param _callbackGasLimit is the gas limit used to estimate the price.
   * @param _requestGasPriceWei is the gas price in wei used for the estimation.
   */
  function estimateRequestPrice(
    uint32 _callbackGasLimit,
    uint256 _requestGasPriceWei
  ) external view override onlyConfiguredNotDisabled returns (uint256) {
    int256 weiPerUnitPli = _getFeedData();
    return _calculateRequestPrice(_callbackGasLimit, _requestGasPriceWei, weiPerUnitPli);
  }

  function _calculateRequestPrice(
    uint256 _gas,
    uint256 _requestGasPrice,
    int256 _weiPerUnitPli
  ) internal view returns (uint256) {
    // costWei is the base fee denominated in wei (native)
    // costWei takes into account the L1 posting costs of the VRF fulfillment
    // transaction, if we are on an L2.
    uint256 costWei = (_requestGasPrice *
      (_gas + s_wrapperGasOverhead + s_coordinatorGasOverhead) +
      ChainSpecificUtil._getL1CalldataGasCost(s_fulfillmentTxSizeBytes));
    // (1e18 juels/pli) * ((wei/gas * (gas)) + l1wei) / (wei/pli) == 1e18 juels * wei/pli / (wei/pli) == 1e18 juels * wei/pli * pli/wei == juels
    // baseFee is the base fee denominated in juels (pli)
    uint256 baseFee = (1e18 * costWei) / uint256(_weiPerUnitPli);
    // feeWithPremium is the fee after the percentage premium is applied
    uint256 feeWithPremium = (baseFee * (s_wrapperPremiumPercentage + 100)) / 100;
    // feeWithFlatFee is the fee after the flat fee is applied on top of the premium
    uint256 feeWithFlatFee = feeWithPremium + (1e12 * uint256(s_fulfillmentFlatFeePliPPM));

    return feeWithFlatFee;
  }

  /**
   * @notice onTokenTransfer is called by PliToken upon payment for a VRF request.
   *
   * @dev Reverts if payment is too low.
   *
   * @param _sender is the sender of the payment, and the address that will receive a VRF callback
   *        upon fulfillment.
   *
   * @param _amount is the amount of PLI paid in Juels.
   *
   * @param _data is the abi-encoded VRF request parameters: uint32 callbackGasLimit,
   *        uint16 requestConfirmations, and uint32 numWords.
   */
  function onTokenTransfer(address _sender, uint256 _amount, bytes calldata _data) external onlyConfiguredNotDisabled {
    // solhint-disable-next-line gas-custom-errors
    require(msg.sender == address(PLI), "only callable from PLI");

    (uint32 callbackGasLimit, uint16 requestConfirmations, uint32 numWords) = abi.decode(
      _data,
      (uint32, uint16, uint32)
    );
    uint32 eip150Overhead = _getEIP150Overhead(callbackGasLimit);
    int256 weiPerUnitPli = _getFeedData();
    uint256 price = _calculateRequestPrice(callbackGasLimit, tx.gasprice, weiPerUnitPli);
    // solhint-disable-next-line gas-custom-errors
    require(_amount >= price, "fee too low");
    // solhint-disable-next-line gas-custom-errors
    require(numWords <= s_maxNumWords, "numWords too high");

    uint256 requestId = COORDINATOR.requestRandomWords(
      s_keyHash,
      SUBSCRIPTION_ID,
      requestConfirmations,
      callbackGasLimit + eip150Overhead + s_wrapperGasOverhead,
      numWords
    );
    s_callbacks[requestId] = Callback({
      callbackAddress: _sender,
      callbackGasLimit: callbackGasLimit,
      requestGasPrice: tx.gasprice,
      requestWeiPerUnitPli: weiPerUnitPli,
      juelsPaid: _amount
    });
    lastRequestId = requestId;
  }

  /**
   * @notice withdraw is used by the VRFV2Wrapper's owner to withdraw PLI revenue.
   *
   * @param _recipient is the address that should receive the PLI funds.
   *
   * @param _amount is the amount of PLI in Juels that should be withdrawn.
   */
  function withdraw(address _recipient, uint256 _amount) external onlyOwner {
    PLI.transfer(_recipient, _amount);
  }

  /**
   * @notice enable this contract so that new requests can be accepted.
   */
  function enable() external onlyOwner {
    s_disabled = false;
  }

  /**
   * @notice disable this contract so that new requests will be rejected. When disabled, new requests
   * @notice will revert but existing requests can still be fulfilled.
   */
  function disable() external onlyOwner {
    s_disabled = true;
  }

  // solhint-disable-next-line plugin-solidity/prefix-internal-functions-with-underscore
  function fulfillRandomWords(uint256 _requestId, uint256[] memory _randomWords) internal override {
    Callback memory callback = s_callbacks[_requestId];
    delete s_callbacks[_requestId];
    // solhint-disable-next-line gas-custom-errors
    require(callback.callbackAddress != address(0), "request not found"); // This should never happen

    VRFV2WrapperConsumerBase c;
    bytes memory resp = abi.encodeWithSelector(c.rawFulfillRandomWords.selector, _requestId, _randomWords);

    bool success = _callWithExactGas(callback.callbackGasLimit, callback.callbackAddress, resp);
    if (!success) {
      emit WrapperFulfillmentFailed(_requestId, callback.callbackAddress);
    }
  }

  function _getFeedData() private view returns (int256) {
    bool staleFallback = s_stalenessSeconds > 0;
    uint256 timestamp;
    int256 weiPerUnitPli;
    (, weiPerUnitPli, , timestamp, ) = PLI_ETH_FEED.latestRoundData();
    // solhint-disable-next-line not-rely-on-time
    if (staleFallback && s_stalenessSeconds < block.timestamp - timestamp) {
      weiPerUnitPli = s_fallbackWeiPerUnitPli;
    }
    // solhint-disable-next-line gas-custom-errors
    require(weiPerUnitPli >= 0, "Invalid PLI wei price");
    return weiPerUnitPli;
  }

  /**
   * @dev Calculates extra amount of gas required for running an assembly call() post-EIP150.
   */
  function _getEIP150Overhead(uint32 gas) private pure returns (uint32) {
    return gas / 63 + 1;
  }

  /**
   * @dev calls target address with exactly gasAmount gas and data as calldata
   * or reverts if at least gasAmount gas is not available.
   */
  function _callWithExactGas(uint256 gasAmount, address target, bytes memory data) private returns (bool success) {
    assembly {
      let g := gas()
      // Compute g -= GAS_FOR_CALL_EXACT_CHECK and check for underflow
      // The gas actually passed to the callee is min(gasAmount, 63//64*gas available).
      // We want to ensure that we revert if gasAmount >  63//64*gas available
      // as we do not want to provide them with less, however that check itself costs
      // gas.  GAS_FOR_CALL_EXACT_CHECK ensures we have at least enough gas to be able
      // to revert if gasAmount >  63//64*gas available.
      if lt(g, GAS_FOR_CALL_EXACT_CHECK) {
        revert(0, 0)
      }
      g := sub(g, GAS_FOR_CALL_EXACT_CHECK)
      // if g - g//64 <= gasAmount, revert
      // (we subtract g//64 because of EIP-150)
      if iszero(gt(sub(g, div(g, 64)), gasAmount)) {
        revert(0, 0)
      }
      // solidity calls check that a contract actually exists at the destination, so we do the same
      if iszero(extcodesize(target)) {
        revert(0, 0)
      }
      // call and return whether we succeeded. ignore return data
      // call(gas,addr,value,argsOffset,argsLength,retOffset,retLength)
      success := call(gasAmount, target, 0, add(data, 0x20), mload(data), 0, 0)
    }
    return success;
  }

  function typeAndVersion() external pure virtual override returns (string memory) {
    return "VRFV2Wrapper 1.0.0";
  }

  modifier onlyConfiguredNotDisabled() {
    // solhint-disable-next-line gas-custom-errors
    require(s_configured, "wrapper is not configured");
    // solhint-disable-next-line gas-custom-errors
    require(!s_disabled, "wrapper is disabled");
    _;
  }
}

// solhint-disable-next-line interface-starts-with-i
interface ExtendedVRFCoordinatorV2Interface is VRFCoordinatorV2Interface {
  function getConfig()
    external
    view
    returns (
      uint16 minimumRequestConfirmations,
      uint32 maxGasLimit,
      uint32 stalenessSeconds,
      uint32 gasAfterPaymentCalculation
    );

  function getFallbackWeiPerUnitPli() external view returns (int256);

  function getFeeConfig()
    external
    view
    returns (
      uint32 fulfillmentFlatFeePliPPMTier1,
      uint32 fulfillmentFlatFeePliPPMTier2,
      uint32 fulfillmentFlatFeePliPPMTier3,
      uint32 fulfillmentFlatFeePliPPMTier4,
      uint32 fulfillmentFlatFeePliPPMTier5,
      uint24 reqsForTier2,
      uint24 reqsForTier3,
      uint24 reqsForTier4,
      uint24 reqsForTier5
    );
}
