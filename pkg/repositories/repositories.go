package repositories

type Repositories interface {
	ProofOfWork() ProofOfWork
	Wallet() Wallet
	Blockchain() Blockchain
	Close() error
}
