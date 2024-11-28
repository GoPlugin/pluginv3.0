// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @notice Implement this contract so that a keeper-compatible contract can monitor
/// and fund the implementation contract with PLI if it falls below a defined threshold.
interface ILinkAvailable {
  function linkAvailableForPayment() external view returns (int256 availableBalance);
}
