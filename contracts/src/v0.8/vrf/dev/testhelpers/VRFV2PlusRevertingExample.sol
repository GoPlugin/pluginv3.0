// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {PliTokenInterface} from "../../../shared/interfaces/PliTokenInterface.sol";
import {VRFConsumerBaseV2Plus} from "../VRFConsumerBaseV2Plus.sol";
import {VRFV2PlusClient} from "../libraries/VRFV2PlusClient.sol";

// VRFV2RevertingExample will always revert. Used for testing only, useless in prod.
contract VRFV2PlusRevertingExample is VRFConsumerBaseV2Plus {
  uint256[] public s_randomWords;
  uint256 public s_requestId;
  // solhint-disable-next-line plugin-solidity/prefix-storage-variables-with-s-underscore
  PliTokenInterface internal PLITOKEN;
  uint256 public s_subId;
  uint256 public s_gasAvailable;

  constructor(address vrfCoordinator, address pli) VRFConsumerBaseV2Plus(vrfCoordinator) {
    PLITOKEN = PliTokenInterface(pli);
  }

  // solhint-disable-next-line plugin-solidity/prefix-internal-functions-with-underscore
  function fulfillRandomWords(uint256, uint256[] calldata) internal pure override {
    // solhint-disable-next-line gas-custom-errors, reason-string
    revert();
  }

  function createSubscriptionAndFund(uint96 amount) external {
    if (s_subId == 0) {
      s_subId = s_vrfCoordinator.createSubscription();
      s_vrfCoordinator.addConsumer(s_subId, address(this));
    }
    // Approve the pli transfer.
    PLITOKEN.transferAndCall(address(s_vrfCoordinator), amount, abi.encode(s_subId));
  }

  function topUpSubscription(uint96 amount) external {
    // solhint-disable-next-line gas-custom-errors
    require(s_subId != 0, "sub not set");
    // Approve the pli transfer.
    PLITOKEN.transferAndCall(address(s_vrfCoordinator), amount, abi.encode(s_subId));
  }

  function updateSubscription(address[] memory consumers) external {
    // solhint-disable-next-line gas-custom-errors
    require(s_subId != 0, "subID not set");
    for (uint256 i = 0; i < consumers.length; i++) {
      s_vrfCoordinator.addConsumer(s_subId, consumers[i]);
    }
  }

  function requestRandomness(
    bytes32 keyHash,
    uint256 subId,
    uint16 minReqConfs,
    uint32 callbackGasLimit,
    uint32 numWords
  ) external returns (uint256) {
    VRFV2PlusClient.RandomWordsRequest memory req = VRFV2PlusClient.RandomWordsRequest({
      keyHash: keyHash,
      subId: subId,
      requestConfirmations: minReqConfs,
      callbackGasLimit: callbackGasLimit,
      numWords: numWords,
      extraArgs: "" // empty extraArgs defaults to pli payment
    });
    s_requestId = s_vrfCoordinator.requestRandomWords(req);
    return s_requestId;
  }
}
