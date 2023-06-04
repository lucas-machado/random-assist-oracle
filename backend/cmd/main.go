package main

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	abireader "random-assist-oracle/internal/contract"
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

	if contract != nil {
		log.Println("success!")
	}
}
