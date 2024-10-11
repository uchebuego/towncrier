package blockchain

import (
	"fmt"

	"github.com/uchebuego/towncrier/pkg/config"
)

type BlockchainManager struct {
	blockchains []*BlockchainListener
}

func (bm *BlockchainManager) Start() {
	for _, blockchain := range bm.blockchains {
		go blockchain.Listen()
	}
}

func NewBlockchainManager(cfg *config.Config) (*BlockchainManager, error) {
	var blockchains []*BlockchainListener
	for _, blockchainConfig := range cfg.Blockchains {
		listener, err := NewBlockchainListener(blockchainConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create listener for blockchain: %v", err)
		}
		blockchains = append(blockchains, listener)
	}

	return &BlockchainManager{
		blockchains: blockchains,
	}, nil
}
