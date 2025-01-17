// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {TypeAndVersionInterface} from "../../interfaces/TypeAndVersionInterface.sol";
import {VRFConsumerBaseV2Plus} from "./VRFConsumerBaseV2Plus.sol";
import {PliTokenInterface} from "../../shared/interfaces/PliTokenInterface.sol";
import {AggregatorV3Interface} from "../../shared/interfaces/AggregatorV3Interface.sol";
import {VRFV2PlusClient} from "./libraries/VRFV2PlusClient.sol";
import {IVRFV2PlusWrapper} from "./interfaces/IVRFV2PlusWrapper.sol";
import {VRFV2PlusWrapperConsumerBase} from "./VRFV2PlusWrapperConsumerBase.sol";

/**
 * @notice A wrapper for VRFCoordinatorV2 that provides an interface better suited to one-off
 * @notice requests for randomness.
 */
// solhint-disable-next-line max-states-count
contract VRFV2PlusWrapper is ConfirmedOwner, TypeAndVersionInterface, VRFConsumerBaseV2Plus, IVRFV2PlusWrapper {
  event WrapperFulfillmentFailed(uint256 indexed requestId, address indexed consumer);

  // upper bound limit for premium percentages to make sure fee calculations don't overflow
  uint8 private constant PREMIUM_PERCENTAGE_MAX = 155;

  // 5k is plenty for an EXTCODESIZE call (2600) + warm CALL (100)
  // and some arithmetic operations.
  uint256 private constant GAS_FOR_CALL_EXACT_CHECK = 5_000;
  uint16 private constant EXPECTED_MIN_LENGTH = 36;

  // solhint-disable-next-line plugin-solidity/prefix-immutable-variables-with-i
  uint256 public immutable SUBSCRIPTION_ID;
  PliTokenInterface internal immutable i_pli;
  AggregatorV3Interface internal immutable i_pli_native_feed;

  event FulfillmentTxSizeSet(uint32 size);
  event ConfigSet(
    uint32 wrapperGasOverhead,
    uint32 coordinatorGasOverheadNative,
    uint32 coordinatorGasOverheadPli,
    uint16 coordinatorGasOverheadPerWord,
    uint8 coordinatorNativePremiumPercentage,
    uint8 coordinatorPliPremiumPercentage,
    bytes32 keyHash,
    uint8 maxNumWords,
    uint32 stalenessSeconds,
    int256 fallbackWeiPerUnitPli,
    uint32 fulfillmentFlatFeeNativePPM,
    uint32 fulfillmentFlatFeePliDiscountPPM
  );
  event FallbackWeiPerUnitPliUsed(uint256 requestId, int256 fallbackWeiPerUnitPli);
  event Withdrawn(address indexed to, uint256 amount);
  event NativeWithdrawn(address indexed to, uint256 amount);
  event Enabled();
  event Disabled();

  error PliAlreadySet();
  error PliDiscountTooHigh(uint32 flatFeePliDiscountPPM, uint32 flatFeeNativePPM);
  error InvalidPremiumPercentage(uint8 premiumPercentage, uint8 max);
  error FailedToTransferPli();
  error IncorrectExtraArgsLength(uint16 expectedMinimumLength, uint16 actualLength);
  error NativePaymentInOnTokenTransfer();
  error PLIPaymentInRequestRandomWordsInNative();
  error SubscriptionIdMissing();

  /* Storage Slot 1: BEGIN */
  // 20 bytes used by VRFConsumerBaseV2Plus.s_vrfCoordinator

  // s_configured tracks whether this contract has been configured. If not configured, randomness
  // requests cannot be made.
  bool public s_configured;

  // s_disabled disables the contract when true. When disabled, new VRF requests cannot be made
  // but existing ones can still be fulfilled.
  bool public s_disabled;

  // s_maxNumWords is the max number of words that can be requested in a single wrapped VRF request.
  uint8 internal s_maxNumWords;

  // 9 bytes left
  /* Storage Slot 1: END */

  /* Storage Slot 2: BEGIN */
  // s_keyHash is the key hash to use when requesting randomness. Fees are paid based on current gas
  // fees, so this should be set to the highest gas lane on the network.
  bytes32 internal s_keyHash;
  /* Storage Slot 2: END */

  /* Storage Slot 3: BEGIN */
  // lastRequestId is the request ID of the most recent VRF V2 request made by this wrapper. This
  // should only be relied on within the same transaction the request was made.
  uint256 public override lastRequestId;
  /* Storage Slot 3: END */

  /* Storage Slot 4: BEGIN */
  // s_fallbackWeiPerUnitPli is the backup PLI exchange rate used when the PLI/NATIVE feed is
  // stale.
  int256 private s_fallbackWeiPerUnitPli;
  /* Storage Slot 4: END */

  /* Storage Slot 5: BEGIN */
  // s_stalenessSeconds is the number of seconds before we consider the feed price to be stale and
  // fallback to fallbackWeiPerUnitPli.
  uint32 private s_stalenessSeconds;

  // s_wrapperGasOverhead reflects the gas overhead of the wrapper's fulfillRandomWords
  // function. The cost for this gas is passed to the user.
  uint32 private s_wrapperGasOverhead;

  // Configuration fetched from VRFCoordinatorV2_5

  /// @dev this is the size of a VRF v2plus fulfillment's calldata abi-encoded in bytes.
  /// @dev proofSize = 13 words = 13 * 256 = 3328 bits
  /// @dev commitmentSize = 10 words = 10 * 256 = 2560 bits
  /// @dev onlyPremiumParameterSize = 256 bits
  /// @dev dataSize = proofSize + commitmentSize + onlyPremiumParameterSize = 6144 bits
  /// @dev function selector = 32 bits
  /// @dev total data size = 6144 bits + 32 bits = 6176 bits = 772 bytes
  uint32 public s_fulfillmentTxSizeBytes = 772;

  // s_coordinatorGasOverheadNative reflects the gas overhead of the coordinator's fulfillRandomWords
  // function for native payment. The cost for this gas is billed to the subscription, and must therefor be included
  // in the pricing for wrapped requests. This includes the gas costs of proof verification and
  // payment calculation in the coordinator.
  uint32 private s_coordinatorGasOverheadNative;

  // s_coordinatorGasOverheadPli reflects the gas overhead of the coordinator's fulfillRandomWords
  // function for pli payment. The cost for this gas is billed to the subscription, and must therefor be included
  // in the pricing for wrapped requests. This includes the gas costs of proof verification and
  // payment calculation in the coordinator.
  uint32 private s_coordinatorGasOverheadPli;

  uint16 private s_coordinatorGasOverheadPerWord;

  // s_fulfillmentFlatFeePliPPM is the flat fee in millionths of native that VRFCoordinatorV2
  // charges for native payment.
  uint32 private s_fulfillmentFlatFeeNativePPM;

  // s_fulfillmentFlatFeePliDiscountPPM is the flat fee discount in millionths of native that VRFCoordinatorV2
  // charges for pli payment.
  uint32 private s_fulfillmentFlatFeePliDiscountPPM;

  // s_coordinatorNativePremiumPercentage is the coordinator's premium ratio in percentage for native payment.
  // For example, a value of 0 indicates no premium. A value of 15 indicates a 15 percent premium.
  // Wrapper has no premium. This premium is for VRFCoordinator.
  uint8 private s_coordinatorNativePremiumPercentage;

  // s_coordinatorPliPremiumPercentage is the premium ratio in percentage for pli payment. For example, a
  // value of 0 indicates no premium. A value of 15 indicates a 15 percent premium.
  // Wrapper has no premium. This premium is for VRFCoordinator.
  uint8 private s_coordinatorPliPremiumPercentage;
  /* Storage Slot 5: END */

  struct Callback {
    address callbackAddress;
    uint32 callbackGasLimit;
    // Reducing requestGasPrice from uint256 to uint64 slots Callback struct
    // into a single word, thus saving an entire SSTORE and leading to 21K
    // gas cost saving. 18 ETH would be the max gas price we can process.
    // GasPrice is unlikely to be more than 14 ETH on most chains
    uint64 requestGasPrice;
  }
  /* Storage Slot 6: BEGIN */
  mapping(uint256 => Callback) /* requestID */ /* callback */ public s_callbacks;
  /* Storage Slot 6: END */

  constructor(
    address _pli,
    address _pliNativeFeed,
    address _coordinator,
    uint256 _subId
  ) VRFConsumerBaseV2Plus(_coordinator) {
    i_pli = PliTokenInterface(_pli);
    i_pli_native_feed = AggregatorV3Interface(_pliNativeFeed);

    if (_subId == 0) {
      revert SubscriptionIdMissing();
    }

    // Sanity check: should revert if the subscription does not exist
    s_vrfCoordinator.getSubscription(_subId);

    // Subscription for the wrapper is created and managed by an external account.
    // Expectation is that wrapper contract address will be added as a consumer
    // to this subscription by the external account (owner of the subscription).
    // Migration of the wrapper's subscription to the new coordinator has to be
    // handled by the external account (owner of the subscription).
    SUBSCRIPTION_ID = _subId;
  }

  /**
   * @notice setFulfillmentTxSize sets the size of the fulfillment transaction in bytes.
   * @param _size is the size of the fulfillment transaction in bytes.
   */
  function setFulfillmentTxSize(uint32 _size) external onlyOwner {
    s_fulfillmentTxSizeBytes = _size;

    emit FulfillmentTxSizeSet(_size);
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
   * @param _coordinatorGasOverheadNative reflects the gas overhead of the coordinator's
   *        fulfillRandomWords function for native payment.
   *
   * @param _coordinatorGasOverheadPli reflects the gas overhead of the coordinator's
   *        fulfillRandomWords function for pli payment.
   *
   * @param _coordinatorGasOverheadPerWord reflects the gas overhead per word of the coordinator's
   *        fulfillRandomWords function.
   *
   * @param _coordinatorNativePremiumPercentage is the coordinator's premium ratio in percentage for requests paid in native.
   *
   * @param _coordinatorPliPremiumPercentage is the coordinator's premium ratio in percentage for requests paid in pli.
   *
   * @param _keyHash to use for requesting randomness.
   * @param _maxNumWords is the max number of words that can be requested in a single wrapped VRF request
   * @param _stalenessSeconds is the number of seconds before we consider the feed price to be stale
   *        and fallback to fallbackWeiPerUnitPli.
   *
   * @param _fallbackWeiPerUnitPli is the backup PLI exchange rate used when the PLI/NATIVE feed
   *        is stale.
   *
   * @param _fulfillmentFlatFeeNativePPM is the flat fee in millionths of native that VRFCoordinatorV2Plus
   *        charges for native payment.
   *
   * @param _fulfillmentFlatFeePliDiscountPPM is the flat fee discount in millionths of native that VRFCoordinatorV2Plus
   *        charges for pli payment.
   */
  /// @dev This function while having only 12 parameters is causing a Stack too deep error when running forge coverage.
  function setConfig(
    uint32 _wrapperGasOverhead,
    uint32 _coordinatorGasOverheadNative,
    uint32 _coordinatorGasOverheadPli,
    uint16 _coordinatorGasOverheadPerWord,
    uint8 _coordinatorNativePremiumPercentage,
    uint8 _coordinatorPliPremiumPercentage,
    bytes32 _keyHash,
    uint8 _maxNumWords,
    uint32 _stalenessSeconds,
    int256 _fallbackWeiPerUnitPli,
    uint32 _fulfillmentFlatFeeNativePPM,
    uint32 _fulfillmentFlatFeePliDiscountPPM
  ) external onlyOwner {
    if (_fulfillmentFlatFeePliDiscountPPM > _fulfillmentFlatFeeNativePPM) {
      revert PliDiscountTooHigh(_fulfillmentFlatFeePliDiscountPPM, _fulfillmentFlatFeeNativePPM);
    }
    if (_coordinatorNativePremiumPercentage > PREMIUM_PERCENTAGE_MAX) {
      revert InvalidPremiumPercentage(_coordinatorNativePremiumPercentage, PREMIUM_PERCENTAGE_MAX);
    }
    if (_coordinatorPliPremiumPercentage > PREMIUM_PERCENTAGE_MAX) {
      revert InvalidPremiumPercentage(_coordinatorPliPremiumPercentage, PREMIUM_PERCENTAGE_MAX);
    }

    s_wrapperGasOverhead = _wrapperGasOverhead;
    s_coordinatorGasOverheadNative = _coordinatorGasOverheadNative;
    s_coordinatorGasOverheadPli = _coordinatorGasOverheadPli;
    s_coordinatorGasOverheadPerWord = _coordinatorGasOverheadPerWord;
    s_coordinatorNativePremiumPercentage = _coordinatorNativePremiumPercentage;
    s_coordinatorPliPremiumPercentage = _coordinatorPliPremiumPercentage;
    s_keyHash = _keyHash;
    s_maxNumWords = _maxNumWords;
    s_configured = true;

    // Get other configuration from coordinator
    s_stalenessSeconds = _stalenessSeconds;
    s_fallbackWeiPerUnitPli = _fallbackWeiPerUnitPli;
    s_fulfillmentFlatFeeNativePPM = _fulfillmentFlatFeeNativePPM;
    s_fulfillmentFlatFeePliDiscountPPM = _fulfillmentFlatFeePliDiscountPPM;

    emit ConfigSet(
      _wrapperGasOverhead,
      _coordinatorGasOverheadNative,
      _coordinatorGasOverheadPli,
      _coordinatorGasOverheadPerWord,
      _coordinatorNativePremiumPercentage,
      _coordinatorPliPremiumPercentage,
      _keyHash,
      _maxNumWords,
      _stalenessSeconds,
      _fallbackWeiPerUnitPli,
      _fulfillmentFlatFeeNativePPM,
      s_fulfillmentFlatFeePliDiscountPPM
    );
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
   * @return fulfillmentFlatFeeNativePPM is the flat fee in millionths of native that VRFCoordinatorV2Plus
   *         charges for native payment.
   *
   * @return fulfillmentFlatFeePliDiscountPPM is the flat fee discount in millionths of native that VRFCoordinatorV2Plus
   *         charges for pli payment.
   *
   * @return wrapperGasOverhead reflects the gas overhead of the wrapper's fulfillRandomWords
   *         function. The cost for this gas is passed to the user.
   *
   * @return coordinatorGasOverheadNative reflects the gas overhead of the coordinator's
   *         fulfillRandomWords function for native payment.
   *
   * @return coordinatorGasOverheadPli reflects the gas overhead of the coordinator's
   *         fulfillRandomWords function for pli payment.
   *
   * @return coordinatorGasOverheadPerWord reflects the gas overhead per word of the coordinator's
   *         fulfillRandomWords function.
   *
   * @return wrapperNativePremiumPercentage is the premium ratio in percentage for native payment. For example, a value of 0
   *         indicates no premium. A value of 15 indicates a 15 percent premium.
   *
   * @return wrapperPliPremiumPercentage is the premium ratio in percentage for pli payment. For example, a value of 0
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
      uint32 fulfillmentFlatFeeNativePPM,
      uint32 fulfillmentFlatFeePliDiscountPPM,
      uint32 wrapperGasOverhead,
      uint32 coordinatorGasOverheadNative,
      uint32 coordinatorGasOverheadPli,
      uint16 coordinatorGasOverheadPerWord,
      uint8 wrapperNativePremiumPercentage,
      uint8 wrapperPliPremiumPercentage,
      bytes32 keyHash,
      uint8 maxNumWords
    )
  {
    return (
      s_fallbackWeiPerUnitPli,
      s_stalenessSeconds,
      s_fulfillmentFlatFeeNativePPM,
      s_fulfillmentFlatFeePliDiscountPPM,
      s_wrapperGasOverhead,
      s_coordinatorGasOverheadNative,
      s_coordinatorGasOverheadPli,
      s_coordinatorGasOverheadPerWord,
      s_coordinatorNativePremiumPercentage,
      s_coordinatorPliPremiumPercentage,
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
    uint32 _callbackGasLimit,
    uint32 _numWords
  ) external view override onlyConfiguredNotDisabled returns (uint256) {
    (int256 weiPerUnitPli, ) = _getFeedData();
    return _calculateRequestPrice(_callbackGasLimit, _numWords, tx.gasprice, weiPerUnitPli);
  }

  function calculateRequestPriceNative(
    uint32 _callbackGasLimit,
    uint32 _numWords
  ) external view override onlyConfiguredNotDisabled returns (uint256) {
    return _calculateRequestPriceNative(_callbackGasLimit, _numWords, tx.gasprice);
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
    uint32 _numWords,
    uint256 _requestGasPriceWei
  ) external view override onlyConfiguredNotDisabled returns (uint256) {
    (int256 weiPerUnitPli, ) = _getFeedData();
    return _calculateRequestPrice(_callbackGasLimit, _numWords, _requestGasPriceWei, weiPerUnitPli);
  }

  function estimateRequestPriceNative(
    uint32 _callbackGasLimit,
    uint32 _numWords,
    uint256 _requestGasPriceWei
  ) external view override onlyConfiguredNotDisabled returns (uint256) {
    return _calculateRequestPriceNative(_callbackGasLimit, _numWords, _requestGasPriceWei);
  }

  /**
   * @notice Returns the L1 fee for the fulfillment calldata payload (always return 0 on L1 chains).
   * @notice Override this function in chain specific way for L2 chains.
   */
  function _getL1CostWei() internal view virtual returns (uint256) {
    return 0;
  }

  function _calculateRequestPriceNative(
    uint256 _gas,
    uint32 _numWords,
    uint256 _requestGasPrice
  ) internal view returns (uint256) {
    // costWei is the base fee denominated in wei (native)
    // (wei/gas) * gas
    uint256 wrapperCostWei = _requestGasPrice * s_wrapperGasOverhead;

    // coordinatorCostWei takes into account the L1 posting costs of the VRF fulfillment transaction, if we are on an L2.
    // (wei/gas) * gas + l1wei
    uint256 coordinatorCostWei = _requestGasPrice *
      (_gas + _getCoordinatorGasOverhead(_numWords, true)) +
      _getL1CostWei();

    // coordinatorCostWithPremiumAndFlatFeeWei is the coordinator cost with the percentage premium and flat fee applied
    // coordinator cost * premium multiplier + flat fee
    uint256 coordinatorCostWithPremiumAndFlatFeeWei = ((coordinatorCostWei *
      (s_coordinatorNativePremiumPercentage + 100)) / 100) + (1e12 * uint256(s_fulfillmentFlatFeeNativePPM));

    return wrapperCostWei + coordinatorCostWithPremiumAndFlatFeeWei;
  }

  function _calculateRequestPrice(
    uint256 _gas,
    uint32 _numWords,
    uint256 _requestGasPrice,
    int256 _weiPerUnitPli
  ) internal view returns (uint256) {
    // costWei is the base fee denominated in wei (native)
    // (wei/gas) * gas
    uint256 wrapperCostWei = _requestGasPrice * s_wrapperGasOverhead;

    // coordinatorCostWei takes into account the L1 posting costs of the VRF fulfillment transaction, if we are on an L2.
    // (wei/gas) * gas + l1wei
    uint256 coordinatorCostWei = _requestGasPrice *
      (_gas + _getCoordinatorGasOverhead(_numWords, false)) +
      _getL1CostWei();

    // coordinatorCostWithPremiumAndFlatFeeWei is the coordinator cost with the percentage premium and flat fee applied
    // coordinator cost * premium multiplier + flat fee
    uint256 coordinatorCostWithPremiumAndFlatFeeWei = ((coordinatorCostWei *
      (s_coordinatorPliPremiumPercentage + 100)) / 100) +
      (1e12 * uint256(s_fulfillmentFlatFeeNativePPM - s_fulfillmentFlatFeePliDiscountPPM));

    // requestPrice is denominated in juels (pli)
    // (1e18 juels/pli) * wei / (wei/pli) = juels
    return (1e18 * (wrapperCostWei + coordinatorCostWithPremiumAndFlatFeeWei)) / uint256(_weiPerUnitPli);
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
    require(msg.sender == address(i_pli), "only callable from PLI");

    (uint32 callbackGasLimit, uint16 requestConfirmations, uint32 numWords, bytes memory extraArgs) = abi.decode(
      _data,
      (uint32, uint16, uint32, bytes)
    );
    checkPaymentMode(extraArgs, true);
    uint32 eip150Overhead = _getEIP150Overhead(callbackGasLimit);
    (int256 weiPerUnitPli, bool isFeedStale) = _getFeedData();
    uint256 price = _calculateRequestPrice(callbackGasLimit, numWords, tx.gasprice, weiPerUnitPli);
    // solhint-disable-next-line gas-custom-errors
    require(_amount >= price, "fee too low");
    // solhint-disable-next-line gas-custom-errors
    require(numWords <= s_maxNumWords, "numWords too high");
    VRFV2PlusClient.RandomWordsRequest memory req = VRFV2PlusClient.RandomWordsRequest({
      keyHash: s_keyHash,
      subId: SUBSCRIPTION_ID,
      requestConfirmations: requestConfirmations,
      callbackGasLimit: callbackGasLimit + eip150Overhead + s_wrapperGasOverhead,
      numWords: numWords,
      extraArgs: extraArgs // empty extraArgs defaults to pli payment
    });
    uint256 requestId = s_vrfCoordinator.requestRandomWords(req);
    s_callbacks[requestId] = Callback({
      callbackAddress: _sender,
      callbackGasLimit: callbackGasLimit,
      requestGasPrice: uint64(tx.gasprice)
    });
    lastRequestId = requestId;

    if (isFeedStale) {
      emit FallbackWeiPerUnitPliUsed(requestId, s_fallbackWeiPerUnitPli);
    }
  }

  function checkPaymentMode(bytes memory extraArgs, bool isPliMode) public pure {
    // If extraArgs is empty, payment mode is PLI by default
    if (extraArgs.length == 0) {
      if (!isPliMode) {
        revert PLIPaymentInRequestRandomWordsInNative();
      }
      return;
    }
    if (extraArgs.length < EXPECTED_MIN_LENGTH) {
      revert IncorrectExtraArgsLength(EXPECTED_MIN_LENGTH, uint16(extraArgs.length));
    }
    // ExtraArgsV1 only has struct {bool nativePayment} as of now
    // The following condition checks if nativePayment in abi.encode of
    // ExtraArgsV1 matches the appropriate function call (onTokenTransfer
    // for PLI and requestRandomWordsInNative for Native payment)
    bool nativePayment = extraArgs[35] == hex"01";
    if (nativePayment && isPliMode) {
      revert NativePaymentInOnTokenTransfer();
    }
    if (!nativePayment && !isPliMode) {
      revert PLIPaymentInRequestRandomWordsInNative();
    }
  }

  function requestRandomWordsInNative(
    uint32 _callbackGasLimit,
    uint16 _requestConfirmations,
    uint32 _numWords,
    bytes calldata extraArgs
  ) external payable override onlyConfiguredNotDisabled returns (uint256 requestId) {
    checkPaymentMode(extraArgs, false);

    uint32 eip150Overhead = _getEIP150Overhead(_callbackGasLimit);
    uint256 price = _calculateRequestPriceNative(_callbackGasLimit, _numWords, tx.gasprice);
    // solhint-disable-next-line gas-custom-errors
    require(msg.value >= price, "fee too low");
    // solhint-disable-next-line gas-custom-errors
    require(_numWords <= s_maxNumWords, "numWords too high");
    VRFV2PlusClient.RandomWordsRequest memory req = VRFV2PlusClient.RandomWordsRequest({
      keyHash: s_keyHash,
      subId: SUBSCRIPTION_ID,
      requestConfirmations: _requestConfirmations,
      callbackGasLimit: _callbackGasLimit + eip150Overhead + s_wrapperGasOverhead,
      numWords: _numWords,
      extraArgs: extraArgs
    });
    requestId = s_vrfCoordinator.requestRandomWords(req);
    s_callbacks[requestId] = Callback({
      callbackAddress: msg.sender,
      callbackGasLimit: _callbackGasLimit,
      requestGasPrice: uint64(tx.gasprice)
    });

    return requestId;
  }

  /**
   * @notice withdraw is used by the VRFV2Wrapper's owner to withdraw PLI revenue.
   *
   * @param _recipient is the address that should receive the PLI funds.
   */
  function withdraw(address _recipient) external onlyOwner {
    uint256 amount = i_pli.balanceOf(address(this));
    if (!i_pli.transfer(_recipient, amount)) {
      revert FailedToTransferPli();
    }

    emit Withdrawn(_recipient, amount);
  }

  /**
   * @notice withdraw is used by the VRFV2Wrapper's owner to withdraw native revenue.
   *
   * @param _recipient is the address that should receive the native funds.
   */
  function withdrawNative(address _recipient) external onlyOwner {
    uint256 amount = address(this).balance;
    (bool success, ) = payable(_recipient).call{value: amount}("");
    // solhint-disable-next-line gas-custom-errors
    require(success, "failed to withdraw native");

    emit NativeWithdrawn(_recipient, amount);
  }

  /**
   * @notice enable this contract so that new requests can be accepted.
   */
  function enable() external onlyOwner {
    s_disabled = false;

    emit Enabled();
  }

  /**
   * @notice disable this contract so that new requests will be rejected. When disabled, new requests
   * @notice will revert but existing requests can still be fulfilled.
   */
  function disable() external onlyOwner {
    s_disabled = true;

    emit Disabled();
  }

  // solhint-disable-next-line plugin-solidity/prefix-internal-functions-with-underscore
  function fulfillRandomWords(uint256 _requestId, uint256[] calldata _randomWords) internal override {
    Callback memory callback = s_callbacks[_requestId];
    delete s_callbacks[_requestId];

    address callbackAddress = callback.callbackAddress;
    // solhint-disable-next-line gas-custom-errors
    require(callbackAddress != address(0), "request not found"); // This should never happen

    VRFV2PlusWrapperConsumerBase c;
    bytes memory resp = abi.encodeWithSelector(c.rawFulfillRandomWords.selector, _requestId, _randomWords);

    bool success = _callWithExactGas(callback.callbackGasLimit, callbackAddress, resp);
    if (!success) {
      emit WrapperFulfillmentFailed(_requestId, callbackAddress);
    }
  }

  function pli() external view override returns (address) {
    return address(i_pli);
  }

  function pliNativeFeed() external view override returns (address) {
    return address(i_pli_native_feed);
  }

  function _getFeedData() private view returns (int256 weiPerUnitPli, bool isFeedStale) {
    uint32 stalenessSeconds = s_stalenessSeconds;
    uint256 timestamp;
    (, weiPerUnitPli, , timestamp, ) = i_pli_native_feed.latestRoundData();
    // solhint-disable-next-line not-rely-on-time
    isFeedStale = stalenessSeconds > 0 && stalenessSeconds < block.timestamp - timestamp;
    if (isFeedStale) {
      weiPerUnitPli = s_fallbackWeiPerUnitPli;
    }
    // solhint-disable-next-line gas-custom-errors
    require(weiPerUnitPli >= 0, "Invalid PLI wei price");
    return (weiPerUnitPli, isFeedStale);
  }

  /**
   * @dev Calculates extra amount of gas required for running an assembly call() post-EIP150.
   */
  function _getEIP150Overhead(uint32 gas) private pure returns (uint32) {
    return gas / 63 + 1;
  }

  function _getCoordinatorGasOverhead(uint32 numWords, bool nativePayment) internal view returns (uint32) {
    if (nativePayment) {
      return s_coordinatorGasOverheadNative + numWords * s_coordinatorGasOverheadPerWord;
    } else {
      return s_coordinatorGasOverheadPli + numWords * s_coordinatorGasOverheadPerWord;
    }
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
    return "VRFV2PlusWrapper 1.0.0";
  }

  modifier onlyConfiguredNotDisabled() {
    // solhint-disable-next-line gas-custom-errors
    require(s_configured, "wrapper is not configured");
    // solhint-disable-next-line gas-custom-errors
    require(!s_disabled, "wrapper is disabled");
    _;
  }
}
