// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {PluginClient, PluginRequestInterface, LinkTokenInterface} from "../../../PluginClient.sol";
import {Plugin} from "../../../Plugin.sol";

contract Consumer is PluginClient {
  using Plugin for Plugin.Request;

  bytes32 internal s_specId;
  bytes32 internal s_currentPrice;

  event RequestFulfilled(
    bytes32 indexed requestId, // User-defined ID
    bytes32 indexed price
  );

  function requestEthereumPrice(string memory _currency, uint256 _payment) public {
    requestEthereumPriceByCallback(_currency, _payment, address(this));
  }

  function requestEthereumPriceByCallback(string memory _currency, uint256 _payment, address _callback) public {
    Plugin.Request memory req = _buildPluginRequest(s_specId, _callback, this.fulfill.selector);
    req._add("get", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](1);
    path[0] = _currency;
    req._addStringArray("path", path);
    _sendPluginRequest(req, _payment);
  }

  function cancelRequest(
    address _oracle,
    bytes32 _requestId,
    uint256 _payment,
    bytes4 _callbackFunctionId,
    uint256 _expiration
  ) public {
    PluginRequestInterface requested = PluginRequestInterface(_oracle);
    requested.cancelOracleRequest(_requestId, _payment, _callbackFunctionId, _expiration);
  }

  function withdrawLink() public {
    LinkTokenInterface _link = LinkTokenInterface(_pluginTokenAddress());
    // solhint-disable-next-line gas-custom-errors
    require(_link.transfer(msg.sender, _link.balanceOf(address(this))), "Unable to transfer");
  }

  function addExternalRequest(address _oracle, bytes32 _requestId) external {
    _addPluginExternalRequest(_oracle, _requestId);
  }

  function fulfill(bytes32 _requestId, bytes32 _price) public recordPluginFulfillment(_requestId) {
    emit RequestFulfilled(_requestId, _price);
    s_currentPrice = _price;
  }

  function getCurrentPrice() public view returns (bytes32) {
    return s_currentPrice;
  }
}
