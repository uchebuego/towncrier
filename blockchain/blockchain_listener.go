package blockchain

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/uchebuego/towncrier/config"
)

type BlockchainListener struct {
	client     *ethclient.Client
	contracts  []*ContractListener
	webhookURL string
	startBlock uint64
}

func (bl *BlockchainListener) Listen() {
	for _, contract := range bl.contracts {
		go contract.Listen()
	}
}

func NewBlockchainListener(cfg config.BlockchainConfig) (*BlockchainListener, error) {
	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %v", err)
	}

	startBlock := cfg.StartBlock
	if startBlock == 0 {
		startBlock, err = getCurrentBlockNumber(client)
		if err != nil {
			return nil, fmt.Errorf("failed to get current block number: %v", err)
		}
	}

	var contracts []*ContractListener
	for _, contractConfig := range cfg.Contracts {
		contractListener, err := NewContractListener(client, contractConfig, cfg.WebhookURL, startBlock)
		if err != nil {
			return nil, fmt.Errorf("failed to create contract listener: %v", err)
		}
		contracts = append(contracts, contractListener)
	}

	return &BlockchainListener{
		client:     client,
		contracts:  contracts,
		webhookURL: cfg.WebhookURL,
		startBlock: startBlock,
	}, nil
}
