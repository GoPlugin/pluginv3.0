// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {BurnMintERC677} from "./BurnMintERC677.sol";

contract LinkToken is BurnMintERC677 {
  constructor() BurnMintERC677("Plugin Token", "PLI", 18, 1e27) {}
}
