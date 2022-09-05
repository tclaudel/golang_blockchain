package config

import (
	"fmt"

	"github.com/jinzhu/configor"
)

type BlockchainNode struct {
	Address string `default:"http://localhost:9092"`
}

type CliRepositories struct {
	Wallet Wallet
}

type Config struct {
	Log            Log
	BlockchainNode BlockchainNode
	Repositories   CliRepositories
}

func NewCliConfig() *Config {
	var cfg = new(Config)

	if err := configor.New(&configor.Config{
		ENVPrefix:            "GOLANG_BLOCKCHAIN_CLI",
		Debug:                false,
		Verbose:              false,
		AutoReload:           false,
		ErrorOnUnmatchedKeys: true,
	}).Load(cfg); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return cfg
}
