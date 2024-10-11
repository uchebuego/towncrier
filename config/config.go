package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Blockchain struct {
		RPCUrl          string   `yaml:"rpc_url"`
		ContractAddress string   `yaml:"contract_address"`
		StartBlock      uint64   `yaml:"start_block"`
		EventNames      []string `yaml:"event_names"`
		ABI             string   `yaml:"abi"`
		ABIFile         string   `yaml:"abi_file"`
		EtherscanAPIKey string   `yaml:"etherscan_api_key"`
	} `yaml:"blockchain"`
	WebhookURL string `yaml:"webhook_url"`
}

func LoadConfig(configFile string) (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	log.Printf("Config loaded from %s", configFile)
	return &config, nil
}
