package repositories

import (
	"github.com/tclaudel/golang_blockchain/internal/values"
	"github.com/tclaudel/golang_blockchain/pkg/entity"
)

type ProofOfWork interface {
	Mine(previousHash values.Hash, transactions []entity.Transaction) (entity.Block, error)
}
