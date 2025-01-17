// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import "./BaseFeeManager.t.sol";

/**
 * @title BaseFeeManagerTest
 * @author Michael Fletcher
 * @notice This contract will test the setup functionality of the feemanager
 */
contract FeeManagerProcessFeeTest is BaseFeeManagerTest {
  function setUp() public override {
    super.setUp();
  }

  function test_WithdrawERC20() public {
    //simulate a fee
    mintPli(address(feeManager), DEFAULT_PLI_MINT_QUANTITY);

    //get the balances to ne used for comparison
    uint256 contractBalance = getPliBalance(address(feeManager));
    uint256 adminBalance = getPliBalance(ADMIN);

    //the amount to withdraw
    uint256 withdrawAmount = contractBalance / 2;

    //withdraw some balance
    withdraw(address(pli), ADMIN, withdrawAmount, ADMIN);

    //check the balance has been reduced
    uint256 newContractBalance = getPliBalance(address(feeManager));
    uint256 newAdminBalance = getPliBalance(ADMIN);

    //check the balance is greater than zero
    assertGt(newContractBalance, 0);
    //check the balance has been reduced by the correct amount
    assertEq(newContractBalance, contractBalance - withdrawAmount);
    //check the admin balance has increased by the correct amount
    assertEq(newAdminBalance, adminBalance + withdrawAmount);
  }

  function test_WithdrawUnwrappedNative() public {
    //issue funds straight to the contract to bypass the lack of fallback function
    issueUnwrappedNative(address(feeManager), DEFAULT_NATIVE_MINT_QUANTITY);

    //get the balances to be used for comparison
    uint256 contractBalance = getNativeUnwrappedBalance(address(feeManager));
    uint256 adminBalance = getNativeUnwrappedBalance(ADMIN);

    //the amount to withdraw
    uint256 withdrawAmount = contractBalance / 2;

    //withdraw some balance
    withdraw(NATIVE_WITHDRAW_ADDRESS, ADMIN, withdrawAmount, ADMIN);

    //check the balance has been reduced
    uint256 newContractBalance = getNativeUnwrappedBalance(address(feeManager));
    uint256 newAdminBalance = getNativeUnwrappedBalance(ADMIN);

    //check the balance is greater than zero
    assertGt(newContractBalance, 0);
    //check the balance has been reduced by the correct amount
    assertEq(newContractBalance, contractBalance - withdrawAmount);
    //check the admin balance has increased by the correct amount
    assertEq(newAdminBalance, adminBalance + withdrawAmount);
  }

  function test_WithdrawNonAdminAddr() public {
    //simulate a fee
    mintPli(address(feeManager), DEFAULT_PLI_MINT_QUANTITY);

    //should revert if not admin
    vm.expectRevert(ONLY_CALLABLE_BY_OWNER_ERROR);

    //withdraw some balance
    withdraw(address(pli), ADMIN, DEFAULT_PLI_MINT_QUANTITY, USER);
  }

  function test_eventIsEmittedAfterSurchargeIsSet() public {
    //native surcharge
    uint64 nativeSurcharge = FEE_SCALAR / 5;

    //expect an emit
    vm.expectEmit();

    //emit the event that is expected to be emitted
    emit NativeSurchargeUpdated(nativeSurcharge);

    //set the surcharge
    setNativeSurcharge(nativeSurcharge, ADMIN);
  }

  function test_subscriberDiscountEventIsEmittedOnUpdate() public {
    //native surcharge
    uint64 discount = FEE_SCALAR / 3;

    //an event should be emitted
    vm.expectEmit();

    //emit the event that is expected to be emitted
    emit SubscriberDiscountUpdated(USER, DEFAULT_FEED_1_V3, address(native), discount);

    //set the surcharge
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), discount, ADMIN);
  }

  function test_eventIsEmittedUponWithdraw() public {
    //simulate a fee
    mintPli(address(feeManager), DEFAULT_PLI_MINT_QUANTITY);

    //the amount to withdraw
    uint192 withdrawAmount = 1;

    //expect an emit
    vm.expectEmit();

    //the event to be emitted
    emit Withdraw(ADMIN, ADMIN, address(pli), withdrawAmount);

    //withdraw some balance
    withdraw(address(pli), ADMIN, withdrawAmount, ADMIN);
  }

  function test_pliAvailableForPaymentReturnsPliBalance() public {
    //simulate a deposit of pli for the conversion pool
    mintPli(address(feeManager), DEFAULT_REPORT_PLI_FEE);

    //check there's a balance
    assertGt(getPliBalance(address(feeManager)), 0);

    //check the pli available for payment is the pli balance
    assertEq(feeManager.pliAvailableForPayment(), getPliBalance(address(feeManager)));
  }

  function test_payPliDeficit() public {
    //get the default payload
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V3));

    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE, USER);

    //not enough funds in the reward pool should trigger an insufficient pli event
    vm.expectEmit();

    IRewardManager.FeePayment[] memory contractFees = new IRewardManager.FeePayment[](1);
    contractFees[0] = IRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_PLI_FEE));

    emit InsufficientPli(contractFees);

    //process the fee
    processFee(payload, USER, address(native), 0);

    //double check the rewardManager balance is 0
    assertEq(getPliBalance(address(rewardManager)), 0);

    //simulate a deposit of pli to cover the deficit
    mintPli(address(feeManager), DEFAULT_REPORT_PLI_FEE);

    vm.expectEmit();
    emit PliDeficitCleared(DEFAULT_CONFIG_DIGEST, DEFAULT_REPORT_PLI_FEE);

    //pay the deficit which will transfer pli from the rewardManager to the rewardManager
    payPliDeficit(DEFAULT_CONFIG_DIGEST, ADMIN);

    //check the rewardManager received the pli
    assertEq(getPliBalance(address(rewardManager)), DEFAULT_REPORT_PLI_FEE);
  }

  function test_payPliDeficitTwice() public {
    //get the default payload
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V3));

    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE, USER);

    //not enough funds in the reward pool should trigger an insufficient pli event
    vm.expectEmit();

    IRewardManager.FeePayment[] memory contractFees = new IRewardManager.FeePayment[](1);
    contractFees[0] = IRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_PLI_FEE));

    //emit the event that is expected to be emitted
    emit InsufficientPli(contractFees);

    //process the fee
    processFee(payload, USER, address(native), 0);

    //double check the rewardManager balance is 0
    assertEq(getPliBalance(address(rewardManager)), 0);

    //simulate a deposit of pli to cover the deficit
    mintPli(address(feeManager), DEFAULT_REPORT_PLI_FEE);

    vm.expectEmit();
    emit PliDeficitCleared(DEFAULT_CONFIG_DIGEST, DEFAULT_REPORT_PLI_FEE);

    //pay the deficit which will transfer pli from the rewardManager to the rewardManager
    payPliDeficit(DEFAULT_CONFIG_DIGEST, ADMIN);

    //check the rewardManager received the pli
    assertEq(getPliBalance(address(rewardManager)), DEFAULT_REPORT_PLI_FEE);

    //paying again should revert with 0
    vm.expectRevert(ZERO_DEFICIT);

    payPliDeficit(DEFAULT_CONFIG_DIGEST, ADMIN);
  }

  function test_payPliDeficitPaysAllFeesProcessed() public {
    //get the default payload
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V3));

    //approve the native to be transferred from the user
    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE * 2, USER);

    //processing the fee will transfer the native from the user to the feeManager
    processFee(payload, USER, address(native), 0);
    processFee(payload, USER, address(native), 0);

    //check the deficit has been increased twice
    assertEq(getPliDeficit(DEFAULT_CONFIG_DIGEST), DEFAULT_REPORT_PLI_FEE * 2);

    //double check the rewardManager balance is 0
    assertEq(getPliBalance(address(rewardManager)), 0);

    //simulate a deposit of pli to cover the deficit
    mintPli(address(feeManager), DEFAULT_REPORT_PLI_FEE * 2);

    vm.expectEmit();
    emit PliDeficitCleared(DEFAULT_CONFIG_DIGEST, DEFAULT_REPORT_PLI_FEE * 2);

    //pay the deficit which will transfer pli from the rewardManager to the rewardManager
    payPliDeficit(DEFAULT_CONFIG_DIGEST, ADMIN);

    //check the rewardManager received the pli
    assertEq(getPliBalance(address(rewardManager)), DEFAULT_REPORT_PLI_FEE * 2);
  }

  function test_payPliDeficitOnlyCallableByAdmin() public {
    vm.expectRevert(ONLY_CALLABLE_BY_OWNER_ERROR);

    payPliDeficit(DEFAULT_CONFIG_DIGEST, USER);
  }
}
