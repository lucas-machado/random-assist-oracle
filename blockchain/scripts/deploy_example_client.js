const hre = require("hardhat");

async function main() {
    const bridgeAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3";

    const Example = await hre.ethers.getContractFactory("ExampleClient");
    const example = await Example.deploy(bridgeAddress);

    await example.deployed();

    console.log(`Example deployed to ${example.address}`);

    const betTxn = await example.bet();
    await betTxn.wait();

    console.log(`Called bet function on ExampleClient`);
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
