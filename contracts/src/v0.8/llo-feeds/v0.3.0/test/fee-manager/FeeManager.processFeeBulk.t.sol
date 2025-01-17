// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import "./BaseFeeManager.t.sol";
import {IRewardManager} from "../../interfaces/IRewardManager.sol";

/**
 * @title BaseFeeManagerTest
 * @author Michael Fletcher
 * @notice This contract will test the functionality of the feeManager processFee
 */
contract FeeManagerProcessFeeTest is BaseFeeManagerTest {
  uint256 internal constant NUMBER_OF_REPORTS = 5;

  function setUp() public override {
    super.setUp();
  }

  function test_processMultiplePliReports() public {
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    bytes[] memory payloads = new bytes[](NUMBER_OF_REPORTS);
    for (uint256 i = 0; i < NUMBER_OF_REPORTS; ++i) {
      payloads[i] = payload;
    }

    approvePli(address(rewardManager), DEFAULT_REPORT_PLI_FEE * NUMBER_OF_REPORTS, USER);

    processFee(payloads, USER, address(pli), DEFAULT_NATIVE_MINT_QUANTITY);

    assertEq(getPliBalance(address(rewardManager)), DEFAULT_REPORT_PLI_FEE * NUMBER_OF_REPORTS);
    assertEq(getPliBalance(address(feeManager)), 0);
    assertEq(getPliBalance(USER), DEFAULT_PLI_MINT_QUANTITY - DEFAULT_REPORT_PLI_FEE * NUMBER_OF_REPORTS);

    //the subscriber (user) should receive funds back and not the proxy, although when live the proxy will forward the funds sent and not cover it seen here
    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY);
    assertEq(PROXY.balance, DEFAULT_NATIVE_MINT_QUANTITY);
  }

  function test_processMultipleWrappedNativeReports() public {
    mintPli(address(feeManager), DEFAULT_REPORT_PLI_FEE * NUMBER_OF_REPORTS + 1);

    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    bytes[] memory payloads = new bytes[](NUMBER_OF_REPORTS);
    for (uint256 i; i < NUMBER_OF_REPORTS; ++i) {
      payloads[i] = payload;
    }

    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS, USER);

    processFee(payloads, USER, address(native), 0);

    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS);
    assertEq(getPliBalance(address(rewardManager)), DEFAULT_REPORT_PLI_FEE * NUMBER_OF_REPORTS);
    assertEq(getPliBalance(address(feeManager)), 1);
    assertEq(getNativeBalance(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS);
  }

  function test_processMultipleUnwrappedNativeReports() public {
    mintPli(address(feeManager), DEFAULT_REPORT_PLI_FEE * NUMBER_OF_REPORTS + 1);

    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    bytes[] memory payloads = new bytes[](NUMBER_OF_REPORTS);
    for (uint256 i; i < NUMBER_OF_REPORTS; ++i) {
      payloads[i] = payload;
    }

    processFee(payloads, USER, address(native), DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS * 2);

    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS);
    assertEq(getPliBalance(address(rewardManager)), DEFAULT_REPORT_PLI_FEE * NUMBER_OF_REPORTS);
    assertEq(getPliBalance(address(feeManager)), 1);

    assertEq(PROXY.balance, DEFAULT_NATIVE_MINT_QUANTITY);
    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS);
  }

  function test_processV1V2V3Reports() public {
    mintPli(address(feeManager), 1);

    bytes memory payloadV1 = abi.encode(
      [DEFAULT_CONFIG_DIGEST, 0, 0],
      getV1Report(DEFAULT_FEED_1_V1),
      new bytes32[](1),
      new bytes32[](1),
      bytes32("")
    );

    bytes memory pliPayloadV2 = getPayload(getV2Report(DEFAULT_FEED_1_V2));
    bytes memory pliPayloadV3 = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    bytes[] memory payloads = new bytes[](5);
    payloads[0] = payloadV1;
    payloads[1] = pliPayloadV2;
    payloads[2] = pliPayloadV2;
    payloads[3] = pliPayloadV3;
    payloads[4] = pliPayloadV3;

    approvePli(address(rewardManager), DEFAULT_REPORT_PLI_FEE * 4, USER);

    processFee(payloads, USER, address(pli), 0);

    assertEq(getNativeBalance(address(feeManager)), 0);
    assertEq(getPliBalance(address(rewardManager)), DEFAULT_REPORT_PLI_FEE * 4);
    assertEq(getPliBalance(address(feeManager)), 1);

    assertEq(getPliBalance(USER), DEFAULT_PLI_MINT_QUANTITY - DEFAULT_REPORT_PLI_FEE * 4);
    assertEq(getNativeBalance(USER), DEFAULT_NATIVE_MINT_QUANTITY - 0);
  }

  function test_processV1V2V3ReportsWithUnwrapped() public {
    mintPli(address(feeManager), DEFAULT_REPORT_PLI_FEE * 4 + 1);

    bytes memory payloadV1 = abi.encode(
      [DEFAULT_CONFIG_DIGEST, 0, 0],
      getV1Report(DEFAULT_FEED_1_V1),
      new bytes32[](1),
      new bytes32[](1),
      bytes32("")
    );

    bytes memory nativePayloadV2 = getPayload(getV2Report(DEFAULT_FEED_1_V2));
    bytes memory nativePayloadV3 = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    bytes[] memory payloads = new bytes[](5);
    payloads[0] = payloadV1;
    payloads[1] = nativePayloadV2;
    payloads[2] = nativePayloadV2;
    payloads[3] = nativePayloadV3;
    payloads[4] = nativePayloadV3;

    processFee(payloads, USER, address(native), DEFAULT_REPORT_NATIVE_FEE * 4);

    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE * 4);
    assertEq(getPliBalance(address(rewardManager)), DEFAULT_REPORT_PLI_FEE * 4);
    assertEq(getPliBalance(address(feeManager)), 1);

    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * 4);
    assertEq(PROXY.balance, DEFAULT_NATIVE_MINT_QUANTITY);
  }

  function test_processMultipleV1Reports() public {
    bytes memory payload = abi.encode(
      [DEFAULT_CONFIG_DIGEST, 0, 0],
      getV1Report(DEFAULT_FEED_1_V1),
      new bytes32[](1),
      new bytes32[](1),
      bytes32("")
    );

    bytes[] memory payloads = new bytes[](NUMBER_OF_REPORTS);
    for (uint256 i = 0; i < NUMBER_OF_REPORTS; ++i) {
      payloads[i] = payload;
    }

    processFee(payloads, USER, address(native), DEFAULT_REPORT_NATIVE_FEE * 5);

    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY);
    assertEq(PROXY.balance, DEFAULT_NATIVE_MINT_QUANTITY);
  }

  function test_eventIsEmittedIfNotEnoughPli() public {
    bytes memory nativePayload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    bytes[] memory payloads = new bytes[](5);
    payloads[0] = nativePayload;
    payloads[1] = nativePayload;
    payloads[2] = nativePayload;
    payloads[3] = nativePayload;
    payloads[4] = nativePayload;

    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE * 5, USER);

    IRewardManager.FeePayment[] memory payments = new IRewardManager.FeePayment[](5);
    payments[0] = IRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_PLI_FEE));
    payments[1] = IRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_PLI_FEE));
    payments[2] = IRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_PLI_FEE));
    payments[3] = IRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_PLI_FEE));
    payments[4] = IRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_PLI_FEE));

    vm.expectEmit();

    emit InsufficientPli(payments);

    processFee(payloads, USER, address(native), 0);

    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE * 5);
    assertEq(getNativeBalance(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * 5);
    assertEq(getPliBalance(USER), DEFAULT_PLI_MINT_QUANTITY);
  }
}
