// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IVRFCoordinatorV2Plus, IVRFSubscriptionV2Plus} from "../dev/interfaces/IVRFCoordinatorV2Plus.sol";
import {VRFV2PlusClient} from "../dev/libraries/VRFV2PlusClient.sol";
import {VRFConsumerBaseV2Plus} from "../dev/VRFConsumerBaseV2Plus.sol";
import {PliTokenInterface} from "../../shared/interfaces/PliTokenInterface.sol";

contract VRFConsumerV2Plus is VRFConsumerBaseV2Plus {
  uint256[] public s_randomWords;
  uint256 public s_requestId;
  IVRFCoordinatorV2Plus internal COORDINATOR;
  PliTokenInterface internal PLITOKEN;
  uint256 public s_subId;
  uint256 public s_gasAvailable;

  constructor(address vrfCoordinator, address pli) VRFConsumerBaseV2Plus(vrfCoordinator) {
    COORDINATOR = IVRFCoordinatorV2Plus(vrfCoordinator);
    PLITOKEN = PliTokenInterface(pli);
  }

  function fulfillRandomWords(uint256 requestId, uint256[] calldata randomWords) internal override {
    require(requestId == s_requestId, "request ID is incorrect");

    s_gasAvailable = gasleft();
    s_randomWords = randomWords;
  }

  function createSubscriptionAndFund(uint96 amount) external {
    if (s_subId == 0) {
      s_subId = COORDINATOR.createSubscription();
      COORDINATOR.addConsumer(s_subId, address(this));
    }
    // Approve the pli transfer.
    PLITOKEN.transferAndCall(address(COORDINATOR), amount, abi.encode(s_subId));
  }

  function topUpSubscription(uint96 amount) external {
    require(s_subId != 0, "sub not set");
    // Approve the pli transfer.
    PLITOKEN.transferAndCall(address(COORDINATOR), amount, abi.encode(s_subId));
  }

  function updateSubscription(address[] memory consumers) external {
    require(s_subId != 0, "subID not set");
    for (uint256 i = 0; i < consumers.length; i++) {
      COORDINATOR.addConsumer(s_subId, consumers[i]);
    }
  }

  function requestRandomness(VRFV2PlusClient.RandomWordsRequest calldata req) external returns (uint256) {
    s_requestId = COORDINATOR.requestRandomWords(req);
    return s_requestId;
  }
}
