pragma solidity 0.8.19;

import "./BaseTest.t.sol";
import {MockPliToken} from "../../mocks/MockPliToken.sol";
import {MockV3Aggregator} from "../../tests/MockV3Aggregator.sol";
import {ExposedVRFCoordinatorV2_5_Arbitrum} from "../dev/testhelpers/ExposedVRFCoordinatorV2_5_Arbitrum.sol";
import {BlockhashStore} from "../dev/BlockhashStore.sol";
import {ArbGasInfo} from "../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol";
import {ArbSys} from "../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol";
import {VmSafe} from "forge-std/Vm.sol";

contract VRFV2CoordinatorV2_5_Arbitrum is BaseTest {
  /// @dev ARBGAS_ADDR is the address of the ArbGasInfo precompile on Arbitrum.
  /// @dev reference: https://github.com/OffchainLabs/nitro/blob/v2.0.14/contracts/src/precompiles/ArbGasInfo.sol#L10
  address private constant ARBGAS_ADDR = address(0x000000000000000000000000000000000000006C);
  ArbGasInfo private constant ARBGAS = ArbGasInfo(ARBGAS_ADDR);

  /// @dev ARBSYS_ADDR is the address of the ArbSys precompile on Arbitrum.
  /// @dev reference: https://github.com/OffchainLabs/nitro/blob/v2.0.14/contracts/src/precompiles/ArbSys.sol#L10
  address private constant ARBSYS_ADDR = address(0x0000000000000000000000000000000000000064);
  ArbSys private constant ARBSYS = ArbSys(ARBSYS_ADDR);

  address internal constant DEPLOYER = 0xD883a6A1C22fC4AbFE938a5aDF9B2Cc31b1BF18B;

  BlockhashStore s_bhs;
  ExposedVRFCoordinatorV2_5_Arbitrum s_testCoordinator;
  MockPliToken s_pliToken;
  MockV3Aggregator s_pliNativeFeed;

  uint256 s_startGas = 0.0038 gwei;
  uint256 s_weiPerUnitGas = 0.003 gwei;

  function setUp() public override {
    BaseTest.setUp();

    // Fund our users.
    vm.roll(1);
    vm.deal(DEPLOYER, 10_000 ether);
    changePrank(DEPLOYER);

    vm.txGasPrice(100 gwei);

    // Instantiate BHS.
    s_bhs = new BlockhashStore();

    // Deploy coordinator, PLI token and PLI/Native feed.
    s_testCoordinator = new ExposedVRFCoordinatorV2_5_Arbitrum(address(s_bhs));
    s_pliToken = new MockPliToken();
    s_pliNativeFeed = new MockV3Aggregator(18, 500000000000000000); // .5 ETH (good for testing)

    // Configure the coordinator.
    s_testCoordinator.setPLIAndPLINativeFeed(address(s_pliToken), address(s_pliNativeFeed));
    s_testCoordinator.setConfig(
      0, // minRequestConfirmations
      2_500_000, // maxGasLimit
      1, // stalenessSeconds
      50_000, // gasAfterPaymentCalculation
      50000000000000000, // fallbackWeiPerUnitPli
      500_000, // fulfillmentFlatFeeNativePPM
      100_000, // fulfillmentFlatFeePliDiscountPPM
      15, // nativePremiumPercentage
      10 // pliPremiumPercentage
    );
  }

  function _encodeCalculatePaymentAmountNativeExternal(
    uint256 startGas,
    uint256 weiPerUnitGas,
    bool onlyPremium
  ) internal pure returns (bytes memory) {
    return
      abi.encodeWithSelector(
        ExposedVRFCoordinatorV2_5_Arbitrum.calculatePaymentAmountNativeExternal.selector,
        startGas,
        weiPerUnitGas,
        onlyPremium
      );
  }

  function _encodeCalculatePaymentAmountPliExternal(
    uint256 startGas,
    uint256 weiPerUnitGas,
    bool onlyPremium
  ) internal pure returns (bytes memory) {
    return
      abi.encodeWithSelector(
        ExposedVRFCoordinatorV2_5_Arbitrum.calculatePaymentAmountPliExternal.selector,
        startGas,
        weiPerUnitGas,
        onlyPremium
      );
  }

  function _checkL1GasFeeEmittedLogs(uint256 expectedL1GasFee) internal {
    VmSafe.Log[] memory entries = vm.getRecordedLogs();
    assertEq(entries.length, 1);
    assertEq(entries[0].topics.length, 1);
    assertEq(entries[0].topics[0], keccak256("L1GasFee(uint256)"));
    // 1e15 is less than 1 percent discrepancy
    uint256 actualL1GasFee = abi.decode(entries[0].data, (uint256));
    assertApproxEqAbs(expectedL1GasFee, actualL1GasFee, 1e15);
  }

  function test_getBlockNumber() public {
    // sanity check that Arbitrum will use ArbSys to get the block number
    vm.mockCall(ARBSYS_ADDR, abi.encodeWithSelector(ARBSYS.arbBlockNumber.selector), abi.encode(33691));
    assertEq(33691, s_testCoordinator.getBlockNumberExternal());
  }

  function test_getBlockhash() public {
    // for blocks within 256 blocks from the current block return the blockhash using ArbSys
    bytes32 testBlockHash = bytes32(keccak256("testBlock"));
    vm.mockCall(ARBSYS_ADDR, abi.encodeWithSelector(ARBSYS.arbBlockNumber.selector), abi.encode(45830));
    vm.mockCall(ARBSYS_ADDR, abi.encodeWithSelector(ARBSYS.arbBlockHash.selector, 45825), abi.encode(testBlockHash));
    assertEq(testBlockHash, s_testCoordinator.getBlockhashExternal(45825));
    // for blocks outside 256 blocks from the current block return nothing
    assertEq("", s_testCoordinator.getBlockhashExternal(33830));
    // for blocks greater than the current block return nothing
    assertEq("", s_testCoordinator.getBlockhashExternal(50550));
  }

  function test_calculatePaymentAmountNative() public {
    // first we test premium and flat fee payment combined
    bool onlyPremium = false;
    bytes memory txMsgData = _encodeCalculatePaymentAmountNativeExternal(s_startGas, s_weiPerUnitGas, onlyPremium);
    vm.mockCall(ARBGAS_ADDR, abi.encodeWithSelector(ARBGAS.getCurrentTxL1GasFees.selector), abi.encode(10 gwei));
    vm.recordLogs();

    uint256 gasLimit = 0.0001 gwei; // needed because gasleft() is used in the payment calculation
    (bool success, bytes memory returnData) = address(s_testCoordinator).call{gas: gasLimit}(txMsgData);
    assertTrue(success);
    _checkL1GasFeeEmittedLogs(uint256(10 gwei));

    uint96 payment = abi.decode(returnData, (uint96));
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 5.000129 * 1e17, 1e15);

    // now we test only premium payment
    onlyPremium = true;
    txMsgData = _encodeCalculatePaymentAmountNativeExternal(s_startGas, s_weiPerUnitGas, onlyPremium);

    (success, returnData) = address(s_testCoordinator).call{gas: gasLimit}(txMsgData);
    assertTrue(success);
    _checkL1GasFeeEmittedLogs(uint256(10 gwei));

    payment = abi.decode(returnData, (uint96));
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 5.000017 * 1e17, 1e15);
  }

  function test_calculatePaymentAmountPli() public {
    // first we test premium and flat fee payment combined
    bool onlyPremium = false;
    bytes memory txMsgData = _encodeCalculatePaymentAmountPliExternal(s_startGas, s_weiPerUnitGas, onlyPremium);
    vm.mockCall(ARBGAS_ADDR, abi.encodeWithSelector(ARBGAS.getCurrentTxL1GasFees.selector), abi.encode(10 gwei));
    vm.recordLogs();

    uint256 gasLimit = 0.0001 gwei; // needed because gasleft() is used in the payment calculation
    (bool success, bytes memory returnData) = address(s_testCoordinator).call{gas: gasLimit}(txMsgData);
    assertTrue(success);
    _checkL1GasFeeEmittedLogs(uint256(10 gwei));

    uint96 payment = abi.decode(returnData, (uint96));
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 8.00024772 * 1e17, 1e15);

    // now we test only premium payment
    onlyPremium = true;
    txMsgData = _encodeCalculatePaymentAmountPliExternal(s_startGas, s_weiPerUnitGas, onlyPremium);

    (success, returnData) = address(s_testCoordinator).call{gas: gasLimit}(txMsgData);
    assertTrue(success);
    _checkL1GasFeeEmittedLogs(uint256(10 gwei));

    payment = abi.decode(returnData, (uint96));
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 8.00002252 * 1e17, 1e15);
  }
}
