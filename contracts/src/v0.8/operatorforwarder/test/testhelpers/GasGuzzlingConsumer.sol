// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Consumer} from "./Consumer.sol";
import {Plugin} from "../../../Plugin.sol";

contract GasGuzzlingConsumer is Consumer {
  using Plugin for Plugin.Request;

  constructor(address _link, address _oracle, bytes32 _specId) {
    _setPluginToken(_link);
    _setPluginOracle(_oracle);
    s_specId = _specId;
  }

  function gassyRequestEthereumPrice(uint256 _payment) public {
    Plugin.Request memory req = _buildPluginRequest(s_specId, address(this), this.gassyFulfill.selector);
    req._add("get", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](1);
    path[0] = "USD";
    req._addStringArray("path", path);
    _sendPluginRequest(req, _payment);
  }

  function gassyFulfill(bytes32 _requestId, bytes32) public recordPluginFulfillment(_requestId) {
    while (true) {}
  }

  function gassyMultiWordRequest(uint256 _payment) public {
    Plugin.Request memory req = _buildPluginRequest(s_specId, address(this), this.gassyMultiWordFulfill.selector);
    req._add("get", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](1);
    path[0] = "USD";
    req._addStringArray("path", path);
    _sendPluginRequest(req, _payment);
  }

  function gassyMultiWordFulfill(bytes32 _requestId, bytes memory) public recordPluginFulfillment(_requestId) {
    while (true) {}
  }
}
