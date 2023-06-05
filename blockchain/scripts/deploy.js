const hre = require("hardhat");

async function main() {
  const Example = await hre.ethers.getContractFactory("Bridge");
  const example = await Example.deploy();

  await example.deployed();

  console.log(`Example deployed to ${example.address}`);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
