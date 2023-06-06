// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

interface IRandomNumberRequester {
    function receiveRandomNumbers(uint256[] calldata _randomNumbers) external;
}

contract Bridge {
    event RandomNumberRequested(address indexed requester);

    function requestRandomNumbers() public {
        emit RandomNumberRequested(msg.sender);
    }

    function forwardRandomNumbers(address requester, uint256[] memory _randomNumbers) public {
        require(_randomNumbers.length > 0, "No random numbers provided");
        IRandomNumberRequester(requester).receiveRandomNumbers(_randomNumbers);
    }
}
