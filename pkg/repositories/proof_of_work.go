package repositories

import "github.com/tclaudel/golang_blockchain/internal/values"

type ProofOfWork interface {
	Mine(previousHash values.Hash, transactions []values.Transaction) (values.Nonce, error)
}
