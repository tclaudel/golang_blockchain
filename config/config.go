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

type Repositories struct {
	ProofOfWork ProofOfWork
}

type Config struct {
	Log          Log
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
