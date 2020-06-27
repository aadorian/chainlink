pragma solidity ^0.6.0;

import "../interfaces/AggregatorInterface.sol";
import "../interfaces/AggregatorV3Interface.sol";

/**
 * @title MockV3Aggregator
 * @notice Based on the FluxAggregator contract
 * @notice Use this contract when you need to test
 * other contract's ability to read data from an
 * aggregator contract, but how the aggregator got
 * its answer is unimportant
 */
contract MockV3Aggregator is AggregatorInterface, AggregatorV3Interface {
  uint256 constant public override version = 0;

  uint8 public override decimals;
  int256 public override latestAnswer;
  uint256 public override latestTimestamp;
  uint256 public override latestRound;

  mapping(uint256 => int256) public override getAnswer;
  mapping(uint256 => uint256) public override getTimestamp;
  mapping(uint256 => uint256) private getStartedAt;
  mapping(uint256 => uint256) private getAnsweredInRound;

  constructor(
    uint8 _decimals,
    int256 _initialAnswer
  ) public {
    decimals = _decimals;
    updateAnswer(_initialAnswer);
  }

  function updateAnswer(
    int256 _answer
  ) public {
    latestAnswer = _answer;
    latestTimestamp = block.timestamp;
    latestRound++;
    getAnswer[latestRound] = _answer;
    getTimestamp[latestRound] = block.timestamp;
    getStartedAt[latestRound] = block.timestamp;
    getAnsweredInRound[latestRound] = latestRound;
  }

  function updateRoundData(
    uint256 _roundId,
    int256 _answer,
    uint256 _timestamp,
    uint256 _startedAt,
    uint256 _answeredInRound
  ) public {
    latestRound = _roundId;
    latestAnswer = _answer;
    latestTimestamp = _timestamp;
    getAnswer[latestRound] = _answer;
    getTimestamp[latestRound] = _timestamp;
    getStartedAt[latestRound] = _startedAt;
    getAnsweredInRound[latestRound] = _answeredInRound;
  }

  function getRoundData(uint256 _roundId)
    external
    view
    override
    returns (
      uint256 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint256 answeredInRound
    )
  {
    return (
      _roundId,
      getAnswer[_roundId],
      getStartedAt[_roundId],
      getTimestamp[_roundId],
      getAnsweredInRound[_roundId]
    );
  }

  function latestRoundData()
    external
    view
    override
    returns (
      uint256 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint256 answeredInRound
    )
  {
    return (
      latestRound,
      getAnswer[latestRound],
      getStartedAt[latestRound],
      getTimestamp[latestRound],
      getAnsweredInRound[latestRound]
    );
  }

  function description()
    external
    view
    override
    returns (string memory)
  {
    return "v0.6/tests/MockV3Aggregator.sol";
  }

  function onTokenTransfer(address, uint256, bytes calldata _data)
    external
  {
    require(_data.length == 0, "transfer doesn't accept calldata");
  }
}
