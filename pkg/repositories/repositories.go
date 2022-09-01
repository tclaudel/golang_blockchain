package repositories

type ServerRepositories interface {
	ProofOfWork() ProofOfWork
	Wallet() Wallet
	Blockchain() Blockchain
	Close() error
}

type CliRepositories interface {
	Wallet() Wallet
}
