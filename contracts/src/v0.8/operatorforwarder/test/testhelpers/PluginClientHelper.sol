// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {PluginClient} from "../../../PluginClient.sol";

contract PluginClientHelper is PluginClient {
  bytes4 public constant FULFILL_SELECTOR = this.fulfill.selector;

  constructor(address link) {
    _setPluginToken(link);
  }

  function sendRequest(address op, uint256 payment) external returns (bytes32) {
    return _sendPluginRequestTo(op, _buildOperatorRequest(bytes32(hex"10"), FULFILL_SELECTOR), payment);
  }

  function cancelRequest(bytes32 requestId, uint256 payment, uint256 expiration) external {
    _cancelPluginRequest(requestId, payment, this.fulfill.selector, expiration);
  }

  function fulfill(bytes32) external {}
}
