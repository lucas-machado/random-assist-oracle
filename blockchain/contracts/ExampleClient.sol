// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import "hardhat/console.sol";
import "./Bridge.sol"; // path to the Bridge contract file

contract ExampleClient {
    Bridge public bridge;
    uint256[] public randomNumbers;
    uint256 public counter;

    constructor(address _bridgeAddress) {
        bridge = Bridge(_bridgeAddress);
    }

    function receiveRandomNumbers(uint256[] memory _randomNumbers) public {
        console.log("Received random numbers! Number of elements: ", _randomNumbers.length);
        randomNumbers = _randomNumbers;
        counter = 0; // reset counter each time new numbers are received
    }

    function getRandomNumber() public returns (uint256) {
        require(counter < randomNumbers.length, "All numbers have been returned");
        return randomNumbers[counter++];
    }

    function bet() public {
        bridge.requestRandomNumbers();
    }
}
