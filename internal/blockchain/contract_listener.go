package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/uchebuego/towncrier/internal/webhook"
	"github.com/uchebuego/towncrier/pkg/config"
)

type ContractListener struct {
	client            *ethclient.Client
	contractABI       abi.ABI
	address           common.Address
	events            []config.EventConfig
	defaultWebhookURL string
	startBlock        uint64
}

func NewContractListener(client *ethclient.Client, cfg config.ContractConfig, blockchainWebhookURL string, blockChainStartBlock uint64) (*ContractListener, error) {
	contractABI, err := LoadABI(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load contract ABI: %v", err)
	}

	return &ContractListener{
		client:            client,
		contractABI:       contractABI,
		address:           common.HexToAddress(cfg.Address),
		events:            cfg.Events,
		defaultWebhookURL: resolveWebhookURL(cfg.WebhookURL, blockchainWebhookURL),
		startBlock:        blockChainStartBlock,
	}, nil
}

func (cl *ContractListener) Listen() {
	for _, eventConfig := range cl.events {
		go cl.listenToEvent(eventConfig)
	}
}

func (cl *ContractListener) listenToEvent(eventConfig config.EventConfig) {
	event, ok := cl.contractABI.Events[eventConfig.Name]
	if !ok {
		log.Printf("Event %s not found in ABI", eventConfig.Name)
		return
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{cl.address},
		Topics:    [][]common.Hash{{event.ID}},
		FromBlock: big.NewInt(int64(cl.startBlock)),
	}

	logs := make(chan types.Log)
	sub, err := cl.client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Printf("Failed to subscribe to logs for event %s: %v", eventConfig.Name, err)
		return
	}

	webhookURL := resolveWebhookURL(eventConfig.WebhookURL, cl.defaultWebhookURL)

	for {
		select {
		case err := <-sub.Err():
			log.Printf("Subscription error for event %s: %v", eventConfig.Name, err)
			return
		case vLog := <-logs:
			cl.processLog(vLog, eventConfig.Name, webhookURL)
		}
	}
}

func (cl *ContractListener) processLog(vLog types.Log, eventName, webhookURL string) {
	dataMap := make(map[string]interface{})

	event, ok := cl.contractABI.Events[eventName]
	if !ok {
		log.Printf("Event %s not found", eventName)
		return
	}

	err := cl.contractABI.UnpackIntoMap(dataMap, eventName, vLog.Data)
	if err != nil {
		log.Printf("Failed to unpack log for event %s: %v", eventName, err)
		return
	}

	topics, err := getTopics(event, vLog.Topics)
	if err != nil {
		log.Printf("Failed to parse topics for event %s: %v", eventName, err)
		return
	}

	payload := map[string]interface{}{
		"metadata": map[string]interface{}{
			"transactionHash":  vLog.TxHash.Hex(),
			"blockNumber":      vLog.BlockNumber,
			"logIndex":         vLog.Index,
			"transactionIndex": vLog.TxIndex,
			"blockHash":        vLog.BlockHash.Hex(),
			"removed":          vLog.Removed,
			"contractAddress":  vLog.Address.Hex(),
		},
		"topics": topics,
		"data":   dataMap,
	}

	eventData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Printf("Error marshalling event data: %v", err)
		return
	}

	log.Printf("Event: %s - Data: %s", eventName, string(eventData))

	err = webhook.Send(webhookURL, payload)
	if err != nil {
		log.Printf("Failed to send event to webhook: %v", err)
	}
}
