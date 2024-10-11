package blockchain

import (
	"fmt"

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
