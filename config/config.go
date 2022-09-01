package config

import (
	"fmt"

	"github.com/jinzhu/configor"
)

type Log struct {
	Format string `default:"default"`
	Level  string `default:"info"`
}

type ProofOfWork struct {
	Type       string `default:"memory"`
	Difficulty int    `default:"3"`
}

type Wallet struct {
	InitialAmount float64 `default:"100"`
	Type          string  `default:"filesystem"`
	FS            struct {
		Path string `default:"./data/wallet"`
	}
	Name string `default:"owner"`
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

type Repositories struct {
	ProofOfWork ProofOfWork
	Blockchain  Blockchain
	Wallet      Wallet
}

type Config struct {
	Log          Log
	Http         Http
	Repositories Repositories
}

func NewConfig() *Config {
	var cfg = new(Config)

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
