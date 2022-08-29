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
	Type       string
	Difficulty int `default:"3"`
}

type Blockchain struct {
	UserAddress string  `default:"OxOTCLAUDEL"`
	Reward      float64 `default:"1.0"`
}

type Repositories struct {
	ProofOfWork ProofOfWork
}

type Config struct {
	Log          Log
	Blockchain   Blockchain
	Repositories Repositories
}

func NewConfig() *Config {
	var cfg = new(Config)

	if err := configor.New(&configor.Config{
		ENVPrefix:            "GOLANG_BLOCKCHAIN",
		Debug:                true,
		Verbose:              true,
		AutoReload:           false,
		ErrorOnUnmatchedKeys: true,
	}).Load(cfg); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return cfg
}
