package blockchain

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/uchebuego/towncrier/config"
)

func GetABI(cfg *config.Config) (string, error) {
	if cfg.Blockchain.ABI != "" {
		log.Println("Using ABI from YAML config.")
		return cfg.Blockchain.ABI, nil
	}

	if cfg.Blockchain.ABIFile != "" {
		abiData, err := os.ReadFile(cfg.Blockchain.ABIFile)
		if err != nil {
			return "", fmt.Errorf("failed to load ABI from file: %v", err)
		}
		log.Println("Using ABI from file.")
		return string(abiData), nil
	}

	if cfg.Blockchain.EtherscanAPIKey != "" && cfg.Blockchain.ContractAddress != "" {
		abiData, err := fetchABIFromEtherscan(cfg.Blockchain.ContractAddress, cfg.Blockchain.EtherscanAPIKey)
		if err != nil {
			return "", fmt.Errorf("failed to fetch ABI from Etherscan: %v", err)
		}
		log.Println("Using ABI fetched from Etherscan.")
		return abiData, nil
	}

	return "", fmt.Errorf("no ABI provided inline, in file, or available via Etherscan")
}

func fetchABIFromEtherscan(contractAddress, apiKey string) (string, error) {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=contract&action=getabi&address=%s&apikey=%s", contractAddress, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
