package main

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"random-assist-oracle/internal/abireader"
	"strconv"
	"strings"
)
import log "github.com/sirupsen/logrus"

func main() {
	log.Println("starting random assist oracle")
	cfg, err := LoadConfig()

	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	client, err := ethclient.Dial(cfg.BlockchainURL)

	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum network: %v", err)
	}

	contractABI, err := abireader.Read(cfg.ContractABIPath)

	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	address := common.HexToAddress(cfg.ContractAddress)
	contract := bind.NewBoundContract(address, parsedABI, client, client, client)

	chainId := new(big.Int)
	_, ok := chainId.SetString(cfg.ChainId, 10)

	if !ok {
		log.Fatalf("could not convert chain id environment variable %s to big.Int", cfg.ChainId)
	}

	salt, err := strconv.ParseInt(cfg.Salt, 10, 64)

	if err != nil {
		log.Fatalf("parsing salt: %v", err)
	}

	oracle, err := NewOracle(client, contract, address, cfg.PrivateKey, chainId, salt)

	if err != nil {
		log.Fatalf("creating oracle: %v", err)
	}

	log.Fatal(oracle.Run())
}
