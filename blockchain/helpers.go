package blockchain

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getCurrentBlockNumber(client *ethclient.Client) (uint64, error) {
	var blockNumber hexutil.Big
	err := client.Client().Call(&blockNumber, "eth_blockNumber")
	if err != nil {
		return 0, fmt.Errorf("failed to get current block number: %v", err)
	}
	return blockNumber.ToInt().Uint64(), nil
}

func resolveWebhookURL(specific, general string) string {
	if specific != "" {
		return specific
	}
	return general
}

func getTopics(event abi.Event, topics []common.Hash) (interface{}, error) {
	indexed := make([]abi.Argument, 0)

	for _, input := range event.Inputs {
		if input.Indexed {
			indexed = append(indexed, input)
		}
	}

	parsed := make(map[string]interface{})

	err := abi.ParseTopicsIntoMap(parsed, indexed, topics[1:])
	if err != nil {
		return err, nil
	}

	return parsed, nil
}
