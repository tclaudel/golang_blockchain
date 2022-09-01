package config

import (
	"fmt"

	"github.com/jinzhu/configor"
)

type ProofOfWork struct {
	Type       string `default:"memory"`
	Difficulty int    `default:"3"`
}

type Http struct {
	Address      string `default:":9092"`
	ReadTimeout  int    `default:"5"`
	WriteTimeout int    `default:"5"`
}

type Blockchain struct {
	Type   string `default:"badger"`
	Badger struct {
		Path   string `default:"./data/blockchain"`
		Logger bool   `default:"false"`
	}
}

type NodeRepositories struct {
	ProofOfWork ProofOfWork
	Blockchain  Blockchain
	Wallet      Wallet
}

type ServerConfig struct {
	Log          Log
	Http         Http
	Repositories NodeRepositories
}

func NewServerConfig() *ServerConfig {
	var cfg = new(ServerConfig)

	if err := configor.New(&configor.Config{
		ENVPrefix:            "GOLANG_BLOCKCHAIN",
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
