// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {AggregatorV3Interface} from "../../shared/interfaces/AggregatorV3Interface.sol";

// Ideally this contract should inherit from VRFCoordinatorV2 and delegate calls to VRFCoordinatorV2
// However, due to exceeding contract size limit, the logic from VRFCoordinatorV2 is ported over to this contract
contract VRFCoordinatorV2TestHelper {
  uint96 internal s_paymentAmount;

  AggregatorV3Interface public immutable PLI_ETH_FEED;

  struct Config {
    uint16 minimumRequestConfirmations;
    uint32 maxGasLimit;
    // Reentrancy protection.
    bool reentrancyLock;
    // stalenessSeconds is how long before we consider the feed price to be stale
    // and fallback to fallbackWeiPerUnitPli.
    uint32 stalenessSeconds;
    // Gas to cover oracle payment after we calculate the payment.
    // We make it configurable in case those operations are repriced.
    uint32 gasAfterPaymentCalculation;
  }
  int256 private s_fallbackWeiPerUnitPli;
  Config private s_config;

  constructor(
    address pliEthFeed // solhint-disable-next-line no-empty-blocks
  ) {
    PLI_ETH_FEED = AggregatorV3Interface(pliEthFeed);
  }

  function calculatePaymentAmountTest(
    uint256 gasAfterPaymentCalculation,
    uint32 fulfillmentFlatFeePliPPM,
    uint256 weiPerUnitGas
  ) external {
    s_paymentAmount = calculatePaymentAmount(
      gasleft(),
      gasAfterPaymentCalculation,
      fulfillmentFlatFeePliPPM,
      weiPerUnitGas
    );
  }

  error InvalidPliWeiPrice(int256 pliWei);
  error PaymentTooLarge();

  function getFeedData() private view returns (int256) {
    uint32 stalenessSeconds = s_config.stalenessSeconds;
    bool staleFallback = stalenessSeconds > 0;
    uint256 timestamp;
    int256 weiPerUnitPli;
    (, weiPerUnitPli, , timestamp, ) = PLI_ETH_FEED.latestRoundData();
    // solhint-disable-next-line not-rely-on-time
    if (staleFallback && stalenessSeconds < block.timestamp - timestamp) {
      weiPerUnitPli = s_fallbackWeiPerUnitPli;
    }
    return weiPerUnitPli;
  }

  // Get the amount of gas used for fulfillment
  function calculatePaymentAmount(
    uint256 startGas,
    uint256 gasAfterPaymentCalculation,
    uint32 fulfillmentFlatFeePliPPM,
    uint256 weiPerUnitGas
  ) internal view returns (uint96) {
    int256 weiPerUnitPli;
    weiPerUnitPli = getFeedData();
    if (weiPerUnitPli <= 0) {
      revert InvalidPliWeiPrice(weiPerUnitPli);
    }
    // (1e18 juels/pli) (wei/gas * gas) / (wei/pli) = juels
    uint256 paymentNoFee = (1e18 * weiPerUnitGas * (gasAfterPaymentCalculation + startGas - gasleft())) /
      uint256(weiPerUnitPli);
    uint256 fee = 1e12 * uint256(fulfillmentFlatFeePliPPM);
    if (paymentNoFee > (1e27 - fee)) {
      revert PaymentTooLarge(); // Payment + fee cannot be more than all of the pli in existence.
    }
    return uint96(paymentNoFee + fee);
  }

  function getPaymentAmount() public view returns (uint96) {
    return s_paymentAmount;
  }
}
