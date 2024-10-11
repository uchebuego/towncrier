# EVM Event Listener Tool

This tool listens for specified events on multiple EVM-based blockchains, decodes them, and logs them out. It also sends the decoded data to webhooks defined in your YAML configuration file. The tool supports various methods of loading contract ABIs and allows users to define webhooks at different levels (blockchain, contract, event).

## Features

- Listen to events from multiple contracts across multiple blockchains
- Send decoded event data to a webhook
- Load contract ABIs from a file, URL, JSON string, or via Etherscan
- Flexible configuration using a YAML file

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
  - [Blockchain Level](#blockchain-level)
  - [Contract Level](#contract-level)
  - [Event Level](#event-level)
  - [ABI Loading Methods](#abi-loading-methods)
- [Running the Tool](#running-the-tool)
- [Payload Structure](#payload-structure)

## Installation

### Prerequisites

Ensure you have Go installed on your machine.

### Clone the repository

```bash
git clone https://github.com/uchebuego/towncrier
cd towncrier
```

### Build the binary

```bash
go build -o evm-event-listener main.go
```

## Usage

You need to pass the path to a YAML configuration file when running the tool:

```bash
./evm-event-listener -config config.yaml
```

## Configuration

The configuration is a YAML file that defines the blockchains, contracts, events, and webhooks the tool will use. Below is an example configuration.

### Example Configuration

```yaml
blockchains:
  - rpc_url: "https://mainnet.infura.io/v3/YOUR_INFURA_KEY"
    webhook_url: "https://webhook.site/blockchain" # Webhook at blockchain level
    contracts:
      - address: "0xContractAddress1"
        abi: "./abi/Contract1.json" # Load ABI from a file
        webhook_url: "https://webhook.site/contract1" # Webhook at contract level
        events:
          - name: "Transfer"
            webhook_url: "https://webhook.site/transfer" # Webhook at event level
      - address: "0xContractAddress2"
        abi_url: "https://example.com/abi/Contract2.json" # Load ABI from a URL
        events:
          - name: "Approval"
      - address: "0xContractAddress3"
        abi_json: |
          [
            {
              "constant": true,
              "inputs": [{"name": "_owner", "type": "address"}],
              "name": "balanceOf",
              "outputs": [{"name": "balance", "type": "uint256"}],
              "type": "function"
            }
          ]  # Load ABI from an embedded JSON string
        events:
          - name: "Transfer"
      - address: "0xContractAddress4"
        abi_source: "etherscan"
        contract_address: "0xContractAddress4"
        api_key: "YOUR_ETHERSCAN_API_KEY" # Load ABI from Etherscan
        events:
          - name: "Swap"
```

### Configuration Details

#### Blockchain Level

- **`rpc_url`**: RPC URL to connect to the blockchain (e.g., Infura or Alchemy).
- **`webhook_url`**: A default webhook URL for all contracts and events under this blockchain (optional).

#### Contract Level

- **`address`**: The Ethereum address of the contract.
- **`abi`**: Path to the ABI file (optional).
- **`abi_url`**: URL to fetch the ABI (optional).
- **`abi_json`**: ABI as a JSON string embedded in the YAML (optional).
- **`abi_source`**: Specify `etherscan` to fetch the ABI from Etherscan (optional).
- **`contract_address`**: Contract address used when fetching the ABI from services like Etherscan.
- **`api_key`**: API key for fetching ABIs from services like Etherscan.
- **`webhook_url`**: Webhook URL specific to this contract (optional).

#### Event Level

- **`name`**: The name of the event to listen for (e.g., "Transfer", "Approval").
- **`webhook_url`**: Webhook URL specific to this event (optional).

### ABI Loading Methods

You can load ABIs in multiple ways:

1. **From a file** (`abi: "./abi/Contract1.json"`)
2. **From a URL** (`abi_url: "https://example.com/abi/Contract2.json"`)
3. **As an embedded JSON string** (`abi_json: | [...]`)
4. **From Etherscan** (`abi_source: "etherscan"`, `contract_address: "0x...", api_key: "YOUR_API_KEY"`)

## Running the Tool

After configuring your YAML file, run the tool with the following command:

```bash
./evm-event-listener -config /path/to/config.yaml
```

The tool will:

- Connect to the specified blockchains.
- Start listening for the events defined in the configuration.
- Decode and log the event data.
- Send the decoded event data to the specified webhook URLs.

### Command-Line Flags

- `-config`: Path to the YAML configuration file.

## Payload Structure

The data sent to the webhook will be structured into three sections: `metadata`, `topics`, and `data`.

### Example Payload

```json
{
  "metadata": {
    "transactionHash": "0x123...",
    "blockNumber": 1234567,
    "logIndex": 1,
    "transactionIndex": 2,
    "blockHash": "0xabc...",
    "removed": false,
    "contractAddress": "0xContractAddress1"
  },
  "topics": {
    "from": "0xFromAddress",
    "to": "0xToAddress"
  },
  "data": {
    "value": 1000000000000000000
  }
}
```

### Metadata

- **`transactionHash`**: The transaction hash.
- **`blockNumber`**: The block number where the event occurred.
- **`logIndex`**: The index of the log in the block.
- **`transactionIndex`**: The index of the transaction in the block.
- **`blockHash`**: The hash of the block where the event occurred.
- **`removed`**: Whether the log was removed due to a chain reorganization.
- **`contractAddress`**: The address of the contract emitting the event.

### Topics

The indexed parameters from the event (e.g., `from`, `to` for a `Transfer` event).

### Data

The non-indexed event parameters (e.g., `value` for a `Transfer` event).

## Conclusion

This tool provides a flexible way to listen to Ethereum events and forward them to webhooks with minimal setup. The multiple ABI loading methods and configurable webhooks make it easy to adapt the tool for a variety of use cases.

Feel free to customize the tool further to suit your needs!
