package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/uchebuego/towncrier/internal/blockchain"
	"github.com/uchebuego/towncrier/pkg/config"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the YAML config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config from %s: %v", *configPath, err)
	}

	bm, err := blockchain.NewBlockchainManager(cfg)
	if err != nil {
		log.Fatalf("Failed to create blockchain manager: %v", err)

	}

	bm.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down...")

}
