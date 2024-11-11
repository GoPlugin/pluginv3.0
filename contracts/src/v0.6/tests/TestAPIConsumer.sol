/** This example code is designed to quickly deploy an example contract using Remix.
 *  If you have never used Remix, try our example walkthrough: https://docs.chain.link/docs/example-walkthrough
 *  You will need testnet ETH and PLI.
 *     - Kovan ETH faucet: https://faucet.kovan.network/
 *     - Kovan PLI faucet: https://kovan.chain.link/
 */
// SPDX-License-Identifier: MIT
pragma solidity ^0.6.0;

import "../interfaces/LinkTokenInterface.sol";
import "../PluginClient.sol";
import "../vendor/Ownable.sol";

/**
 * @title TestAPIConsumer is an example contract which requests data from
 * the Plugin network
 * @dev This contract is designed to work on multiple networks, including
 * local test networks
 */
contract TestAPIConsumer is PluginClient, Ownable {
  uint256 public currentRoundID = 0;
  uint256 public data;
  bytes4 public selector;

  event PerfMetricsEvent(uint256 roundID, bytes32 requestId, uint256 timestamp);


  /**
   * @notice Deploy the contract with a specified address for the PLI
   * and Oracle contract addresses
   * @dev Sets the storage for the specified addresses
   * @param _link The address of the PLI token contract
   */
  constructor(address _link) public {
    if (_link == address(0)) {
      setPublicPluginToken();
    } else {
      setPluginToken(_link);
    }
  }

  /**
   * @notice Returns the address of the PLI token
   * @dev This is the public implementation for pluginTokenAddress, which is
   * an internal method of the PluginClient contract
   */
  function getPluginToken() public view returns (address) {
    return pluginTokenAddress();
  }

  /**
   * @notice Creates a request to the specified Oracle contract address
   * @dev This function ignores the stored Oracle contract address and
   * will instead send the request to the address specified
   * @param _oracle The Oracle contract address to send the request to
   * @param _jobId The bytes32 JobID to be executed
   * @param _url The URL to fetch data from
   * @param _path The dot-delimited path to parse of the response
   * @param _times The number to multiply the result by
   */
  function createRequestTo(
    address _oracle,
    bytes32 _jobId,
    uint256 _payment,
    string memory _url,
    string memory _path,
    int256 _times
  )
    public
    onlyOwner
    returns (bytes32 requestId)
  {
    selector = this.fulfill.selector;
    Plugin.Request memory req = buildPluginRequest(_jobId, address(this), this.fulfill.selector);
    req.add("get", _url);
    req.add("path", _path);
    req.addInt("times", _times);
    requestId = sendPluginRequestTo(_oracle, req, _payment);
  }

  /**
   * @notice The fulfill method from requests created by this contract
   * @dev The recordPluginFulfillment protects this function from being called
   * by anyone other than the oracle address that the request was sent to
   * @param _requestId The ID that was generated for the request
   * @param _data The answer provided by the oracle
   */
  function fulfill(bytes32 _requestId, uint256 _data)
    public
  {
    data = _data;
    currentRoundID += 1;
    emit PerfMetricsEvent(currentRoundID, _requestId, block.timestamp);
  }

  /**
   * @notice Allows the owner to withdraw any PLI balance on the contract
   */
  function withdrawLink() public onlyOwner {
    LinkTokenInterface link = LinkTokenInterface(pluginTokenAddress());
    require(link.transfer(msg.sender, link.balanceOf(address(this))), "Unable to transfer");
  }

  /**
   * @notice Call this method if no response is received within 5 minutes
   * @param _requestId The ID that was generated for the request to cancel
   * @param _payment The payment specified for the request to cancel
   * @param _callbackFunctionId The bytes4 callback function ID specified for
   * the request to cancel
   * @param _expiration The expiration generated for the request to cancel
   */
  function cancelRequest(
    bytes32 _requestId,
    uint256 _payment,
    bytes4 _callbackFunctionId,
    uint256 _expiration
  )
    public
    onlyOwner
  {
    cancelPluginRequest(_requestId, _payment, _callbackFunctionId, _expiration);
  }
}
