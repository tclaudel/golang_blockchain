package config

type Log struct {
	Format string `default:"text"`
	Level  string `default:"debug"`
}

type Wallet struct {
	MiningReward float64 `default:"100"`
	Type         string  `default:"filesystem"`
	FS           struct {
		Path string `default:"./data/wallet"`
	}
	Name string `default:"blockchain_wallet.json"`
}
