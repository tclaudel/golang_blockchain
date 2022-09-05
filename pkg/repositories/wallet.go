package repositories

import "github.com/tclaudel/golang_blockchain/internal/values"

type Wallet interface {
	Save(wallet values.Wallet) error
	Load(identifier string) (values.Wallet, error)
}
