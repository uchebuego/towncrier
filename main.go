package main

import (
	"flag"
	"log"

	"github.com/uchebuego/towncrier/blockchain"
	"github.com/uchebuego/towncrier/config"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the YAML config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config from %s: %v", *configPath, err)
	}

	abiJSON, err := blockchain.GetABI(cfg)
	if err != nil {
		log.Fatalf("Error retrieving ABI: %v", err)
	}

	listener, err := blockchain.NewEventListener(cfg.Blockchain.RPCUrl, cfg.Blockchain.ContractAddress, abiJSON)
	if err != nil {
		log.Fatalf("Failed to create event listener: %v", err)
	}

	err = listener.Listen(cfg.Blockchain.StartBlock, cfg.Blockchain.EventNames, cfg.WebhookURL)
	if err != nil {
		log.Fatalf("Error while listening to events: %v", err)
	}
}
