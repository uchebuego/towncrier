package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/uchebuego/towncrier/webhook"
)

type EventListener struct {
	client      *ethclient.Client
	contractABI abi.ABI
	address     common.Address
}

func NewEventListener(rpcURL, contractAddress, abiJSON string) (*EventListener, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %v", err)
	}

	contractABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	return &EventListener{
		client:      client,
		contractABI: contractABI,
		address:     common.HexToAddress(contractAddress),
	}, nil
}

func (el *EventListener) Listen(startBlock uint64, eventNames []string, webhookURL string) error {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{el.address},
		FromBlock: big.NewInt(int64(startBlock)),
	}

	logs := make(chan types.Log)
	sub, err := el.client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return fmt.Errorf("failed to subscribe to logs: %v", err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Printf("Subscription error: %v", err)
			return err
		case vLog := <-logs:
			el.processLog(vLog, eventNames, webhookURL)
		}
	}
}

func (el *EventListener) processLog(vLog types.Log, eventNames []string, webhookURL string) {
	for _, eventName := range eventNames {
		event, ok := el.contractABI.Events[eventName]
		if !ok {
			continue
		}

		if vLog.Topics[0] != event.ID {
			continue
		}

		dataMap := make(map[string]interface{})
		err := el.contractABI.UnpackIntoMap(dataMap, eventName, vLog.Data)
		if err != nil {
			log.Printf("Failed to unpack log data: %v", err)
			continue
		}

		payload := map[string]interface{}{
			"metadata": map[string]interface{}{
				"transactionHash":  vLog.TxHash.Hex(),
				"blockNumber":      vLog.BlockNumber,
				"blockHash":        vLog.BlockHash.Hex(),
				"transactionIndex": vLog.TxIndex,
				"logIndex":         vLog.Index,
				"removed":          vLog.Removed,
				"contractAddress":  vLog.Address.Hex(),
			},
			"topics": map[string]interface{}{
				"from": common.HexToAddress(vLog.Topics[1].Hex()),
				"to":   common.HexToAddress(vLog.Topics[2].Hex()),
			},
			"data": dataMap,
		}

		eventData, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			log.Printf("Error marshalling event data: %v", err)
			continue
		}

		log.Printf("Event: %s - Data: %s", eventName, string(eventData))

		err = webhook.Send(webhookURL, payload)
		if err != nil {
			log.Printf("Failed to send event to webhook: %v", err)
		}
	}
}
