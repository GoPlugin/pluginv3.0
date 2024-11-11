pragma solidity 0.4.24;

import "./Consumer.sol";

contract BasicConsumer is Consumer {

  constructor(address _link, address _oracle, bytes32 _specId) public {
    setPluginToken(_link);
    setPluginOracle(_oracle);
    specId = _specId;
  }

}
