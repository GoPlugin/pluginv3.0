pragma solidity ^0.8.0;

import {MaliciousPlugin} from "./MaliciousPlugin.sol";
import {MaliciousPlugined, Plugin} from "./MaliciousPlugined.sol";
import {PluginRequestInterface} from "../../../interfaces/PluginRequestInterface.sol";

contract MaliciousRequester is MaliciousPlugined {
  uint256 private constant ORACLE_PAYMENT = 1 ether;
  uint256 private s_expiration;

  constructor(address _link, address _oracle) {
    setLinkToken(_link);
    setOracle(_oracle);
  }

  function maliciousWithdraw() public {
    MaliciousPlugin.WithdrawRequest memory req = newWithdrawRequest(
      "specId",
      address(this),
      this.doesNothing.selector
    );
    pluginWithdrawRequest(req, ORACLE_PAYMENT);
  }

  function request(bytes32 _id, address _target, bytes memory _callbackFunc) public returns (bytes32 requestId) {
    Plugin.Request memory req = newRequest(_id, _target, bytes4(keccak256(_callbackFunc)));
    s_expiration = block.timestamp + 5 minutes; // solhint-disable-line not-rely-on-time
    return pluginRequest(req, ORACLE_PAYMENT);
  }

  function maliciousPrice(bytes32 _id) public returns (bytes32 requestId) {
    Plugin.Request memory req = newRequest(_id, address(this), this.doesNothing.selector);
    return pluginPriceRequest(req, ORACLE_PAYMENT);
  }

  function maliciousTargetConsumer(address _target) public returns (bytes32 requestId) {
    Plugin.Request memory req = newRequest("specId", _target, bytes4(keccak256("fulfill(bytes32,bytes32)")));
    return pluginTargetRequest(_target, req, ORACLE_PAYMENT);
  }

  function maliciousRequestCancel(bytes32 _id, bytes memory _callbackFunc) public {
    PluginRequestInterface oracle = PluginRequestInterface(oracleAddress());
    oracle.cancelOracleRequest(
      request(_id, address(this), _callbackFunc),
      ORACLE_PAYMENT,
      this.maliciousRequestCancel.selector,
      s_expiration
    );
  }

  function doesNothing(bytes32, bytes32) public pure {} // solhint-disable-line no-empty-blocks
}
