package repositories

import "github.com/tclaudel/golang_blockchain/internal/values"

type Wallet interface {
	Save(wallet values.Wallet) error
	BatchSave(wallets []values.Wallet) error
	Load(wallets []values.Wallet) error
}
