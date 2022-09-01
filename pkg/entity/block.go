package entity

import (
	"github.com/tclaudel/golang_blockchain/internal/values"
)

type Block interface {
	Timestamp() values.Timestamp
	Nonce() values.Nonce
	PreviousHash() values.Hash
	Transactions() ([]values.Transaction, error)
	Hash() values.Hash
}
