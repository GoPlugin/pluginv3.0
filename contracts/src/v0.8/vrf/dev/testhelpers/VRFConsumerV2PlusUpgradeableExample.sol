// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {PliTokenInterface} from "../../../shared/interfaces/PliTokenInterface.sol";
import {IVRFCoordinatorV2Plus} from "../interfaces/IVRFCoordinatorV2Plus.sol";
import {VRFConsumerBaseV2Upgradeable} from "../VRFConsumerBaseV2Upgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {VRFV2PlusClient} from "../libraries/VRFV2PlusClient.sol";

contract VRFConsumerV2PlusUpgradeableExample is Initializable, VRFConsumerBaseV2Upgradeable {
  uint256[] public s_randomWords;
  uint256 public s_requestId;
  IVRFCoordinatorV2Plus public COORDINATOR;
  PliTokenInterface public PLITOKEN;
  uint256 public s_subId;
  uint256 public s_gasAvailable;

  function initialize(address _vrfCoordinator, address _pli) public initializer {
    __VRFConsumerBaseV2_init(_vrfCoordinator);
    COORDINATOR = IVRFCoordinatorV2Plus(_vrfCoordinator);
    PLITOKEN = PliTokenInterface(_pli);
  }

  // solhint-disable-next-line plugin-solidity/prefix-internal-functions-with-underscore
  function fulfillRandomWords(uint256 requestId, uint256[] memory randomWords) internal override {
    // solhint-disable-next-line gas-custom-errors
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
    // solhint-disable-next-line gas-custom-errors
    require(s_subId != 0, "sub not set");
    // Approve the pli transfer.
    PLITOKEN.transferAndCall(address(COORDINATOR), amount, abi.encode(s_subId));
  }

  function updateSubscription(address[] memory consumers) external {
    // solhint-disable-next-line gas-custom-errors
    require(s_subId != 0, "subID not set");
    for (uint256 i = 0; i < consumers.length; i++) {
      COORDINATOR.addConsumer(s_subId, consumers[i]);
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
    s_requestId = COORDINATOR.requestRandomWords(req);
    return s_requestId;
  }
}
