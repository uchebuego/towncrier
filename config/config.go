package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type EventConfig struct {
	Name       string `yaml:"name"`
	WebhookURL string `yaml:"webhook_url,omitempty"`
}

type ContractConfig struct {
	Address         string        `yaml:"address"`
	ABIPath         string        `yaml:"abi,omitempty"`
	ABI             string        `yaml:"abi_json,omitempty"`
	ABIURL          string        `yaml:"abi_url,omitempty"`
	ABISource       string        `yaml:"abi_source,omitempty"`
	APIKey          string        `yaml:"api_key,omitempty"`
	WebhookURL      string        `yaml:"webhook_url,omitempty"`
	Events          []EventConfig `yaml:"events"`
	ContractAddress string        `yaml:"contract_address,omitempty"`
}

type BlockchainConfig struct {
	RPCURL     string           `yaml:"rpc_url"`
	WebhookURL string           `yaml:"webhook_url,omitempty"`
	Contracts  []ContractConfig `yaml:"contracts"`
	StartBlock uint64           `yaml:"start_block"`
}

type Config struct {
	Blockchains []BlockchainConfig `yaml:"blockchains"`
}

func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
