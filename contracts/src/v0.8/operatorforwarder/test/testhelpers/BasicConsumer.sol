// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Consumer} from "./Consumer.sol";

contract BasicConsumer is Consumer {
  constructor(address _pli, address _oracle, bytes32 _specId) {
    _setPluginToken(_pli);
    _setPluginOracle(_oracle);
    s_specId = _specId;
  }
}
