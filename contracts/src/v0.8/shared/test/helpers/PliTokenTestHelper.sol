// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {PliToken} from "../../token/ERC677/PliToken.sol";

// This contract exists to mirror the functionality of the old token, which
// always deployed with 1b tokens sent to the deployer.
contract PliTokenTestHelper is PliToken {
  constructor() {
    _mint(msg.sender, 1e27);
  }
}
