// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {Test} from "forge-std/Test.sol";
import {DestinationFeeManager} from "../../DestinationFeeManager.sol";
import {DestinationRewardManager} from "../../DestinationRewardManager.sol";
import {Common} from "../../../libraries/Common.sol";
import {ERC20Mock} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/mocks/ERC20Mock.sol";
import {WERC20Mock} from "../../../../shared/mocks/WERC20Mock.sol";
import {IDestinationRewardManager} from "../../interfaces/IDestinationRewardManager.sol";
import {DestinationFeeManagerProxy} from "../mocks/DestinationFeeManagerProxy.sol";

/**
 * @title BaseDestinationFeeManagerTest
 * @author Michael Fletcher
 * @notice Base class for all feeManager tests
 * @dev This contract is intended to be inherited from and not used directly. It contains functionality to setup the feeManager
 */
contract BaseDestinationFeeManagerTest is Test {
  //contracts
  DestinationFeeManager internal feeManager;
  DestinationRewardManager internal rewardManager;
  DestinationFeeManagerProxy internal feeManagerProxy;

  ERC20Mock internal pli;
  WERC20Mock internal native;

  //erc20 config
  uint256 internal constant DEFAULT_PLI_MINT_QUANTITY = 100 ether;
  uint256 internal constant DEFAULT_NATIVE_MINT_QUANTITY = 100 ether;

  //contract owner
  address internal constant INVALID_ADDRESS = address(0);
  address internal constant ADMIN = address(uint160(uint256(keccak256("ADMIN"))));
  address internal constant USER = address(uint160(uint256(keccak256("USER"))));
  address internal constant PROXY = address(uint160(uint256(keccak256("PROXY"))));

  //version masks
  bytes32 internal constant V_MASK = 0x0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff;
  bytes32 internal constant V1_BITMASK = 0x0001000000000000000000000000000000000000000000000000000000000000;
  bytes32 internal constant V2_BITMASK = 0x0002000000000000000000000000000000000000000000000000000000000000;
  bytes32 internal constant V3_BITMASK = 0x0003000000000000000000000000000000000000000000000000000000000000;

  //feed ids & config digests
  bytes32 internal constant DEFAULT_FEED_1_V1 = (keccak256("ETH-USD") & V_MASK) | V1_BITMASK;
  bytes32 internal constant DEFAULT_FEED_1_V2 = (keccak256("ETH-USD") & V_MASK) | V2_BITMASK;
  bytes32 internal constant DEFAULT_FEED_1_V3 = (keccak256("ETH-USD") & V_MASK) | V3_BITMASK;

  bytes32 internal constant DEFAULT_FEED_2_V3 = (keccak256("PLI-USD") & V_MASK) | V3_BITMASK;
  bytes32 internal constant DEFAULT_CONFIG_DIGEST = keccak256("DEFAULT_CONFIG_DIGEST");
  bytes32 internal constant DEFAULT_CONFIG_DIGEST2 = keccak256("DEFAULT_CONFIG_DIGEST2");

  //report
  uint256 internal constant DEFAULT_REPORT_PLI_FEE = 1e10;
  uint256 internal constant DEFAULT_REPORT_NATIVE_FEE = 1e12;

  //rewards
  uint64 internal constant FEE_SCALAR = 1e18;

  address internal constant NATIVE_WITHDRAW_ADDRESS = address(0);

  //the selector for each error
  bytes4 internal immutable INVALID_DISCOUNT_ERROR = DestinationFeeManager.InvalidDiscount.selector;
  bytes4 internal immutable INVALID_ADDRESS_ERROR = DestinationFeeManager.InvalidAddress.selector;
  bytes4 internal immutable INVALID_SURCHARGE_ERROR = DestinationFeeManager.InvalidSurcharge.selector;
  bytes4 internal immutable EXPIRED_REPORT_ERROR = DestinationFeeManager.ExpiredReport.selector;
  bytes4 internal immutable INVALID_DEPOSIT_ERROR = DestinationFeeManager.InvalidDeposit.selector;
  bytes4 internal immutable INVALID_QUOTE_ERROR = DestinationFeeManager.InvalidQuote.selector;
  bytes4 internal immutable UNAUTHORIZED_ERROR = DestinationFeeManager.Unauthorized.selector;
  bytes4 internal immutable POOLID_MISMATCH_ERROR = DestinationFeeManager.PoolIdMismatch.selector;
  bytes internal constant ONLY_CALLABLE_BY_OWNER_ERROR = "Only callable by owner";
  bytes internal constant INSUFFICIENT_ALLOWANCE_ERROR = "ERC20: insufficient allowance";
  bytes4 internal immutable ZERO_DEFICIT = DestinationFeeManager.ZeroDeficit.selector;

  //events emitted
  event SubscriberDiscountUpdated(address indexed subscriber, bytes32 indexed feedId, address token, uint64 discount);
  event NativeSurchargeUpdated(uint64 newSurcharge);
  event InsufficientPli(IDestinationRewardManager.FeePayment[] feesAndRewards);
  event Withdraw(address adminAddress, address recipient, address assetAddress, uint192 quantity);
  event PliDeficitCleared(bytes32 indexed configDigest, uint256 pliQuantity);
  event DiscountApplied(
    bytes32 indexed configDigest,
    address indexed subscriber,
    Common.Asset fee,
    Common.Asset reward,
    uint256 appliedDiscountQuantity
  );

  function setUp() public virtual {
    //change to admin user
    vm.startPrank(ADMIN);

    //init required contracts
    _initializeContracts();
  }

  function _initializeContracts() internal {
    pli = new ERC20Mock("PLI", "PLI", ADMIN, 0);
    native = new WERC20Mock();

    feeManagerProxy = new DestinationFeeManagerProxy();
    rewardManager = new DestinationRewardManager(address(pli));
    feeManager = new DestinationFeeManager(
      address(pli),
      address(native),
      address(feeManagerProxy),
      address(rewardManager)
    );

    //pli the feeManager to the proxy
    feeManagerProxy.setDestinationFeeManager(address(feeManager));

    //pli the feeManager to the reward manager
    rewardManager.addFeeManager(address(feeManager));

    //mint some tokens to the admin
    pli.mint(ADMIN, DEFAULT_PLI_MINT_QUANTITY);
    native.mint(ADMIN, DEFAULT_NATIVE_MINT_QUANTITY);
    vm.deal(ADMIN, DEFAULT_NATIVE_MINT_QUANTITY);

    //mint some tokens to the user
    pli.mint(USER, DEFAULT_PLI_MINT_QUANTITY);
    native.mint(USER, DEFAULT_NATIVE_MINT_QUANTITY);
    vm.deal(USER, DEFAULT_NATIVE_MINT_QUANTITY);

    //mint some tokens to the proxy
    pli.mint(PROXY, DEFAULT_PLI_MINT_QUANTITY);
    native.mint(PROXY, DEFAULT_NATIVE_MINT_QUANTITY);
    vm.deal(PROXY, DEFAULT_NATIVE_MINT_QUANTITY);
  }

  function setSubscriberDiscount(
    address subscriber,
    bytes32 feedId,
    address token,
    uint256 discount,
    address sender
  ) internal {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //set the discount
    feeManager.updateSubscriberDiscount(subscriber, feedId, token, uint64(discount));

    //change back to the original address
    changePrank(originalAddr);
  }

  function setSubscriberGlobalDiscount(address subscriber, address token, uint256 discount, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //set the discount
    feeManager.updateSubscriberGlobalDiscount(subscriber, token, uint64(discount));

    //change back to the original address
    changePrank(originalAddr);
  }

  function setNativeSurcharge(uint256 surcharge, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //set the surcharge
    feeManager.setNativeSurcharge(uint64(surcharge));

    //change back to the original address
    changePrank(originalAddr);
  }

  // solium-disable-next-line no-unused-vars
  function getFee(bytes memory report, address quote, address subscriber) public view returns (Common.Asset memory) {
    //get the fee
    (Common.Asset memory fee, , ) = feeManager.getFeeAndReward(subscriber, report, quote);

    return fee;
  }

  function getReward(bytes memory report, address quote, address subscriber) public view returns (Common.Asset memory) {
    //get the reward
    (, Common.Asset memory reward, ) = feeManager.getFeeAndReward(subscriber, report, quote);

    return reward;
  }

  function getAppliedDiscount(bytes memory report, address quote, address subscriber) public view returns (uint256) {
    //get the reward
    (, , uint256 appliedDiscount) = feeManager.getFeeAndReward(subscriber, report, quote);

    return appliedDiscount;
  }

  function getV1Report(bytes32 feedId) public pure returns (bytes memory) {
    return abi.encode(feedId, uint32(0), int192(0), int192(0), int192(0), uint64(0), bytes32(0), uint64(0), uint64(0));
  }

  function getV2Report(bytes32 feedId) public view returns (bytes memory) {
    return
      abi.encode(
        feedId,
        uint32(0),
        uint32(0),
        uint192(DEFAULT_REPORT_NATIVE_FEE),
        uint192(DEFAULT_REPORT_PLI_FEE),
        uint32(block.timestamp),
        int192(0)
      );
  }

  function getV3Report(bytes32 feedId) public view returns (bytes memory) {
    return
      abi.encode(
        feedId,
        uint32(0),
        uint32(0),
        uint192(DEFAULT_REPORT_NATIVE_FEE),
        uint192(DEFAULT_REPORT_PLI_FEE),
        uint32(block.timestamp),
        int192(0),
        int192(0),
        int192(0)
      );
  }

  function getV3ReportWithCustomExpiryAndFee(
    bytes32 feedId,
    uint256 expiry,
    uint256 pliFee,
    uint256 nativeFee
  ) public pure returns (bytes memory) {
    return
      abi.encode(
        feedId,
        uint32(0),
        uint32(0),
        uint192(nativeFee),
        uint192(pliFee),
        uint32(expiry),
        int192(0),
        int192(0),
        int192(0)
      );
  }

  function getPliQuote() public view returns (address) {
    return address(pli);
  }

  function getNativeQuote() public view returns (address) {
    return address(native);
  }

  function withdraw(address assetAddress, address recipient, uint256 amount, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //set the surcharge
    feeManager.withdraw(assetAddress, recipient, uint192(amount));

    //change back to the original address
    changePrank(originalAddr);
  }

  function getPliBalance(address balanceAddress) public view returns (uint256) {
    return pli.balanceOf(balanceAddress);
  }

  function getNativeBalance(address balanceAddress) public view returns (uint256) {
    return native.balanceOf(balanceAddress);
  }

  function getNativeUnwrappedBalance(address balanceAddress) public view returns (uint256) {
    return balanceAddress.balance;
  }

  function mintPli(address recipient, uint256 amount) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(ADMIN);

    //mint the pli to the recipient
    pli.mint(recipient, amount);

    //change back to the original address
    changePrank(originalAddr);
  }

  function mintNative(address recipient, uint256 amount, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //mint the native to the recipient
    native.mint(recipient, amount);

    //change back to the original address
    changePrank(originalAddr);
  }

  function issueUnwrappedNative(address recipient, uint256 quantity) public {
    vm.deal(recipient, quantity);
  }

  function ProcessFeeAsUser(
    bytes32 poolId,
    bytes memory payload,
    address subscriber,
    address tokenAddress,
    uint256 wrappedNativeValue,
    address sender
  ) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //process the fee
    feeManager.processFee{value: wrappedNativeValue}(poolId, payload, abi.encode(tokenAddress), subscriber);

    //change ProcessFeeAsUserback to the original address
    changePrank(originalAddr);
  }

  function processFee(
    bytes32 poolId,
    bytes memory payload,
    address subscriber,
    address feeAddress,
    uint256 wrappedNativeValue
  ) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(subscriber);

    //process the fee
    feeManagerProxy.processFee{value: wrappedNativeValue}(poolId, payload, abi.encode(feeAddress));

    //change back to the original address
    changePrank(originalAddr);
  }

  function processFee(
    bytes32[] memory poolIds,
    bytes[] memory payloads,
    address subscriber,
    address feeAddress,
    uint256 wrappedNativeValue
  ) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(subscriber);

    //process the fee
    feeManagerProxy.processFeeBulk{value: wrappedNativeValue}(poolIds, payloads, abi.encode(feeAddress));

    //change back to the original address
    changePrank(originalAddr);
  }

  function getPayload(bytes memory reportPayload) public pure returns (bytes memory) {
    return abi.encode([DEFAULT_CONFIG_DIGEST, 0, 0], reportPayload, new bytes32[](1), new bytes32[](1), bytes32(""));
  }

  function approvePli(address spender, uint256 quantity, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //approve the pli to be transferred
    pli.approve(spender, quantity);

    //change back to the original address
    changePrank(originalAddr);
  }

  function approveNative(address spender, uint256 quantity, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //approve the pli to be transferred
    native.approve(spender, quantity);

    //change back to the original address
    changePrank(originalAddr);
  }

  function payPliDeficit(bytes32 configDigest, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //approve the pli to be transferred
    feeManager.payPliDeficit(configDigest);

    //change back to the original address
    changePrank(originalAddr);
  }

  function getPliDeficit(bytes32 configDigest) public view returns (uint256) {
    return feeManager.s_pliDeficit(configDigest);
  }
}
