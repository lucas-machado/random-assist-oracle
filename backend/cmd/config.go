package main

import (
	"errors"
	"fmt"
	"github.com/ardanlabs/conf"
)

type Config struct {
	Key             string `conf:"env:PRIVATE_KEY"`
	BlockchainURL   string `conf:"default:http://localhost:8545,env:BLOCKCHAIN_URL"`
	ContractABIPath string `conf:"default:contract/contract.json,env:CONTRACT_ABI_PATH"`
	ContractAddress string `conf:"env:CONTRACT_ADDRESS"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	help, err := conf.ParseOSArgs("APP", &cfg)

	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil, fmt.Errorf("parsing config: %w", err)
		}
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return &cfg, nil
}
