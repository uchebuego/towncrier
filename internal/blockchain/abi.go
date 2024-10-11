package blockchain

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/uchebuego/towncrier/pkg/config"
)

func LoadABI(contractConfig config.ContractConfig) (abi.ABI, error) {
	if contractConfig.ABIPath != "" {
		abiJSON, err := os.ReadFile(contractConfig.ABIPath)
		if err != nil {
			return abi.ABI{}, fmt.Errorf("failed to read ABI file: %v", err)
		}
		return abi.JSON(strings.NewReader(string(abiJSON)))
	}

	if contractConfig.ABIURL != "" {
		resp, err := http.Get(contractConfig.ABIURL)
		if err != nil {
			return abi.ABI{}, fmt.Errorf("failed to fetch ABI from URL: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return abi.ABI{}, fmt.Errorf("failed to fetch ABI, status: %s", resp.Status)
		}

		abiJSON, err := io.ReadAll(resp.Body)
		if err != nil {
			return abi.ABI{}, fmt.Errorf("failed to read ABI from response: %v", err)
		}

		return abi.JSON(strings.NewReader(string(abiJSON)))
	}

	if contractConfig.ABI != "" {
		return abi.JSON(strings.NewReader(contractConfig.ABI))
	}

	if contractConfig.ABISource == "etherscan" && contractConfig.APIKey != "" {
		return loadABIFromEtherscan(contractConfig.ContractAddress, contractConfig.APIKey)
	}

	return abi.ABI{}, errors.New("no valid ABI source provided")
}

func loadABIFromEtherscan(contractAddress, apiKey string) (abi.ABI, error) {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=contract&action=getabi&address=%s&apikey=%s", contractAddress, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to fetch ABI from Etherscan: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return abi.ABI{}, fmt.Errorf("failed to fetch ABI from Etherscan, status: %s", resp.Status)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to decode Etherscan response: %v", err)
	}

	if result["status"] != "1" {
		return abi.ABI{}, fmt.Errorf("failed to retrieve ABI: %v", result["message"])
	}

	abiString := result["result"].(string)
	return abi.JSON(strings.NewReader(abiString))
}
