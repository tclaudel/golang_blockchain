# Golang Blockchain

A simple blockchain implementation in Go
This project is composed of two parts:
- A blockchain node, which is a simple HTTP server that exposes a REST API to interact with the blockchain;
- A simple CLI (command line interface) that allows to interact with the blockchain node.

## Build

A Makefile is provided to build and run the project.
- Running `make all` will build the blockchain node and the CLI as well as their dependencies.
```bash
$> make
mkdir -p ./data/blockchain ./data/wallet
go mod tidy
oapi-codegen -package rest -generate client,types ./docs/swagger.yaml > ./pkg/interfaces/http/rest/api.gen.go
go build -o bin/server cmd/server/main.go
go build -o bin/cli cmd/cli/main.go
```

## Blockchain node

The blockchain node is a Bitcoin-like implementation of a blockchain.
It is composed of a REST API that allows interactions with the blockchain.
This API is composed of the following endpoints:
- `GET /transactions`: returns the list of transactions in the mempool
- `POST /transactions`: add a new transaction to the mempool
- `GET /blocks`: returns the list of blocks in the blockchain
- `POST /blocks`: append a block composed of the transactions in the mempool to the blockchain

The blockchain node has a mempool of transactions that are not yet included in a block.
The mempool is a list of transactions that are waiting to be included in a block.

The blockchain node has a simple proof of work algorithm to validate the blocks.
The blockchain node is also responsible for mining new blocks.
It is configured with a difficulty that defines the number of leading zeros that the hash of a block must have.
The blockchain node will try to find a nonce that will make the hash of a block have the required number of leading zeros.
The blockchain node will also reward the miner with a coinbase transaction.

The blockchain node needs some environment variables to be configured:
- `GOLANG_BLOCKCHAIN_LOG_LEVEL` (default: `info`): the log level of the blockchain node
- `GOLANG_BLOCKCHAIN_LOG_FORMAT` (default: `text`): the log format of the blockchain node
- `GOLANG_BLOCKCHAIN_HTTP_ADDRESS` (default: `:8080`): the HTTP address of the blockchain node
- `GOLANG_BLOCKCHAIN_HTTP_READTIMEOUT` (default: `5s`): the read timeout of the blockchain node
- `GOLANG_BLOCKCHAIN_HTTP_WRITETIMEOUT` (default: `10s`): the write timeout of the blockchain node
- `GOLANG_BLOCKCHAIN_REPOSITORIES_PROOFOFWORK_TYPE` (default: `memory`): the type of the proof of work repository
- `GOLANG_BLOCKCHAIN_REPOSITORIES_PROOFOFWORK_DIFFICULTY` (default: `3`): the difficulty of the proof of work
  (2: easiest, 3: easy, 4: medium, 5: hard)
- `GOLANG_BLOCKCHAIN_REPOSITORIES_BLOCKCHAIN_TYPE` (default: `badger`): the type of the blockchain repository
- `GOLANG_BLOCKCHAIN_REPOSITORIES_BLOCKCHAIN_BADGER_PATH` (default: `./data/blockchain`): the path of the badger database files
- `GOLANG_BLOCKCHAIN_REPOSITORIES_BLOCKCHAIN_BADGER_LOGGER` (default: `false`): toggles badger logging
- `GOLANG_BLOCKCHAIN_REPOSITORIES_WALLET_MININGREWARD` (default: `10`): the mining reward of the wallet
- `GOLANG_BLOCKCHAIN_REPOSITORIES_WALLET_TYPE` (default: `filesystem`): the type of the wallet repository
- `GOLANG_BLOCKCHAIN_REPOSITORIES_WALLET_FS_PATH` (default: `./data/wallet`): the path of the wallet files
- `GOLANG_BLOCKCHAIN_REPOSITORIES_WALLET_NAME` (default: `blockchain_wallet.json`): the default wallet name

## Run

```bash
./bin/server
```

## CLI

The cli is a simple command line interface that allows to interact with the blockchain node. 
It is composed of a set of commands that allows to interact with the blockchain node. 
The cli is configured with the following env variables:
- `GOLANG_BLOCKCHAIN_CLI_LOG_FORMAT` (default: `text`): the log format of the cli
- `GOLANG_BLOCKCHAIN_CLI_LOG_LEVEL` (default: `info`): the log level of the cli
- `GOLANG_BLOCKCHAIN_CLI_BLOCKCHAINNODE_ADDRESS` (default: `http://localhost:8080`): the address of the blockchain node
- `OLANG_BLOCKCHAIN_CLI_REPOSITORIES_WALLET_TYPE` (default: `filesystem`): the type of the wallet repository
- `GOLANG_BLOCKCHAIN_CLI_REPOSITORIES_WALLET_FS` (default: `./data/wallet`): the path of the wallet files

The cli set of commands is composed of the following commands:
- `blockchain`: returns the blockchain blocks
- `blockchain commit`: commit a new block to the blockchain
- `transactions`: returns the mempool transactions
- `transaction create`: create a new transaction
