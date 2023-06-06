// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

contract Bridge {
    uint256[] public randomNumbers;
    uint256 public counter;

    event RandomNumberRequested(address indexed requester);

    function requestRandomNumbers() public {
        emit RandomNumberRequested(msg.sender);
    }
}
