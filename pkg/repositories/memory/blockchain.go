package memory

import (
	"github.com/tclaudel/golang_blockchain/config"
	dEntity "github.com/tclaudel/golang_blockchain/internal/entity"
	"github.com/tclaudel/golang_blockchain/pkg/entity"
	"github.com/tclaudel/golang_blockchain/pkg/repositories"
	"go.uber.org/zap"
)

type Blockchain struct {
	chain []entity.Block
}

func (b Blockchain) Chain() []entity.Block {
	return b.chain
}

type BlockchainRepository struct {
	logger     *zap.Logger
	blockchain *Blockchain
}

func (b *BlockchainRepository) Append(block entity.Block) error {
	b.blockchain.chain = append(b.blockchain.chain, block)
	return nil
}

func (b *BlockchainRepository) Close() error {
	return nil
}

func (b *BlockchainRepository) Get() (entity.Blockchain, error) {
	return b.blockchain, nil
}

func NewBlockchainRepository(logger *zap.Logger, blockchainCfg config.Blockchain) (repositories.Blockchain, error) {
	bcRepo := &BlockchainRepository{
		logger:     logger,
		blockchain: new(Blockchain),
	}

	err := bcRepo.Append(dEntity.Genesis)
	if err != nil {
		return nil, err
	}

	return bcRepo, nil
}
