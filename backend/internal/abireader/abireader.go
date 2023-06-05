package abireader

import (
	"encoding/json"
	"os"
)

type contract struct {
	Abi json.RawMessage `json:"abi"`
}

func Read(path string) (string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var contract contract
	err = json.Unmarshal(file, &contract)
	if err != nil {
		return "", err
	}

	return string(contract.Abi), nil
}
