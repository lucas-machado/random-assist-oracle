// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

// Uncomment this line to use console.log
import "hardhat/console.sol";

contract Bridge {
    uint256[] public randomNumbers;
    uint256 public counter;

    function receiveRandomNumbers(uint256[] memory _randomNumbers) public {
        console.log("Received random numbers! Number of elements: ", _randomNumbers.length);
        randomNumbers = _randomNumbers;
        counter = 0; // reset counter each time new numbers are received
    }


    function getRandomNumber() public returns (uint256) {
        require(counter < randomNumbers.length, "All numbers have been returned");
        return randomNumbers[counter++];
    }
}
