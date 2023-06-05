package main

import (
	"errors"
	"fmt"
	"github.com/ardanlabs/conf"
)

type Config struct {
	PrivateKey      string `conf:"env:PRIVATE_KEY"`
	BlockchainURL   string `conf:"default:http://localhost:8545,env:BLOCKCHAIN_URL"`
	ContractABIPath string `conf:"default:contract/Bridge.json,env:CONTRACT_ABI_PATH"`
	ContractAddress string `conf:"default:0x5fbdb2315678afecb367f032d93f642f64180aa3,env:CONTRACT_ADDRESS"`
	ChainId         string `conf:"default:1337,env:CHAIN_ID"`
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
