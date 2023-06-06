package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"math/big"
	"math/rand"
	"time"
)

type Oracle struct {
	client          *ethclient.Client
	contract        *bind.BoundContract
	contractAddress common.Address
	privateKey      *ecdsa.PrivateKey
	randSource      rand.Source
	chainID         *big.Int
	nextBlockToRead uint64
}

func NewOracle(client *ethclient.Client,
	contract *bind.BoundContract,
	contractAddress common.Address,
	privateKeyHex string,
	chainID *big.Int,
	salt int64) (*Oracle, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	source := rand.NewSource(time.Now().UnixNano() + salt)

	return &Oracle{
		client:          client,
		contract:        contract,
		contractAddress: contractAddress,
		privateKey:      privateKey,
		randSource:      source,
		chainID:         chainID,
	}, nil
}

func (o *Oracle) Run() error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			addresses, err := o.readRandomNumbersRequests()
			if err != nil {
				return err
			}
			fmt.Println("Addresses that requested random numbers:", addresses)

			for _, address := range addresses {
				err := o.sendRandomNumbers(address)
				if err != nil {
					log.Errorf("sending random numbers to address %v: %v", addresses, err)
				}
			}
		}
	}
}

func (o *Oracle) readRandomNumbersRequests() ([]common.Address, error) {
	var addresses []common.Address

	// Define a query to filter logs
	query := ethereum.FilterQuery{
		Addresses: []common.Address{o.contractAddress},
		FromBlock: big.NewInt(int64(o.nextBlockToRead)),
	}

	// Get logs from the blockchain
	logs, err := o.client.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, err
	}

	if len(logs) > 0 {
		if logs[len(logs)-1].BlockNumber < o.nextBlockToRead {
			return addresses, nil
		} else {
			o.nextBlockToRead = logs[len(logs)-1].BlockNumber + 1
		}
	} else {
		header, err := o.client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return nil, err
		}
		if o.nextBlockToRead < header.Number.Uint64() {
			o.nextBlockToRead++
		}
	}

	// Iterate over the logs
	for _, log := range logs {
		// Get the event name from the topics
		eventName := log.Topics[0].Hex()

		// Check if it's a RandomNumberRequested event
		if eventName == crypto.Keccak256Hash([]byte("RandomNumberRequested(address)")).Hex() {
			// Decode the address from the data field
			requesterAddress := common.BytesToAddress(log.Topics[1].Bytes())
			addresses = append(addresses, requesterAddress)
		}
	}

	return addresses, nil
}

func (o *Oracle) sendRandomNumbers(destination common.Address) error {
	// Create a new random source and generate random numbers
	r := rand.New(o.randSource)
	randomNumbers := make([]*big.Int, 1000)
	for i := 0; i < 1000; i++ {
		randomNumbers[i] = big.NewInt(r.Int63())
	}

	// Create a new transactor
	auth, err := bind.NewKeyedTransactorWithChainID(o.privateKey, o.chainID)
	if err != nil {
		return err
	}

	auth.GasLimit = uint64(30000000)        // in units
	auth.GasPrice = big.NewInt(20000000000) // in wei
	auth.Value = big.NewInt(0)              // in wei

	// Call the contract method
	tx, err := o.contract.Transact(auth, "forwardRandomNumbers", destination, randomNumbers)
	if err != nil {
		return err
	}

	fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex()) // for debug

	// Check the receipt
	receipt, err := bind.WaitMined(context.Background(), o.client, tx)
	if err != nil {
		return err
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return errors.New("transaction failed")
	}

	return nil
}
