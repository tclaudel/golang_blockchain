package repositories

import "github.com/tclaudel/golang_blockchain/pkg/entity"

type Blockchain interface {
	Append(entity.Block) error
	Get() (entity.Blockchain, error)
	Close() error
}
