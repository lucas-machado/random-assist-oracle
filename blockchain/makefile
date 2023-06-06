.PHONY: node
node:
	npx hardhat node

.PHONY: dev
dev:
	npx hardhat compile
	npx hardhat run scripts/deploy_bridge.js --network localhost
	npx hardhat run scripts/deploy_example_client.js --network localhost