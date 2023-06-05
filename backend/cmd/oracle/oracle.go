package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"math/rand"
	"time"
)

type Oracle struct {
	client     *ethclient.Client
	contract   *bind.BoundContract
	privateKey *ecdsa.PrivateKey
	randSource rand.Source
	chainID    *big.Int
}

func NewOracle(client *ethclient.Client, contract *bind.BoundContract, privateKeyHex string, chainID *big.Int) (*Oracle, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	source := rand.NewSource(time.Now().UnixNano())

	return &Oracle{
		client:     client,
		contract:   contract,
		privateKey: privateKey,
		randSource: source,
		chainID:    chainID,
	}, nil
}

func (o *Oracle) SendRandomNumbers() error {
	// Create a new random source and generate random numbers
	r := rand.New(o.randSource)
	randomNumbers := make([]*big.Int, 1000)
	for i := 0; i < 1000; i++ {
		randomNumbers[i] = big.NewInt(int64(r.Int63()))
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
	tx, err := o.contract.Transact(auth, "receiveRandomNumbers", randomNumbers)
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
