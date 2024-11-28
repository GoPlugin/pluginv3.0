pragma solidity ^0.8.0;

import {Plugined, Plugin} from "./Plugined.sol";

// solhint-disable
contract MaliciousConsumer is Plugined {
  uint256 private constant ORACLE_PAYMENT = 1 ether;
  uint256 private s_expiration;

  constructor(address _link, address _oracle) payable {
    setLinkToken(_link);
    setOracle(_oracle);
  }

  fallback() external payable {} // solhint-disable-line no-empty-blocks

  function requestData(bytes32 _id, bytes memory _callbackFunc) public {
    Plugin.Request memory req = newRequest(_id, address(this), bytes4(keccak256(_callbackFunc)));
    s_expiration = block.timestamp + 5 minutes;
    pluginRequest(req, ORACLE_PAYMENT);
  }

  function assertFail(bytes32, bytes32) public pure {
    assert(1 == 2);
  }

  function cancelRequestOnFulfill(bytes32 _requestId, bytes32) public {
    _cancelPluginRequest(_requestId, ORACLE_PAYMENT, this.cancelRequestOnFulfill.selector, s_expiration);
  }

  function remove() public {
    selfdestruct(payable(address(0)));
  }

  function stealEthCall(bytes32 _requestId, bytes32) public recordPluginFulfillment(_requestId) {
    (bool success, ) = address(this).call{value: 100}("");
    require(success, "Call failed");
  }

  function stealEthSend(bytes32 _requestId, bytes32) public recordPluginFulfillment(_requestId) {
    require(payable(address(this)).send(100), "Send failed");
  }

  function stealEthTransfer(bytes32 _requestId, bytes32) public recordPluginFulfillment(_requestId) {
    payable(address(this)).transfer(100);
  }

  function doesNothing(bytes32, bytes32) public pure {} // solhint-disable-line no-empty-blocks
}
