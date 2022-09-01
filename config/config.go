package config

type Log struct {
	Format string `default:"default"`
	Level  string `default:"info"`
}

type Wallet struct {
	InitialAmount float64 `default:"100"`
	Type          string  `default:"filesystem"`
	FS            struct {
		Path string `default:"./data/wallet"`
	}
	Name string `default:"owner"`
}
