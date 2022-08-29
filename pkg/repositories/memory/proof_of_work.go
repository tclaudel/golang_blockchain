package memory

import (
	"github.com/tclaudel/golang_blockchain/internal/entity"
	"github.com/tclaudel/golang_blockchain/internal/values"
	"github.com/tclaudel/golang_blockchain/pkg/repositories"
	"go.uber.org/zap"
)

type ProofOfWork struct {
	logger     *zap.Logger
	difficulty values.MiningDifficulty
}

func (p *ProofOfWork) Mine(previousHash values.Hash, transactions []values.Transaction) (values.Nonce, error) {
	nonce := values.InitNonce()
	for {
		valid, err := p.validate(nonce, previousHash, transactions, p.difficulty)
		if err != nil {
			return values.Nonce{}, err
		}

		if !valid {
			nonce = nonce.Increase()
			continue
		}

		break
	}
	p.logger.Debug("Proof of work", zap.Int("nonce", nonce.Int()))

	return nonce, nil
}

func (p *ProofOfWork) validate(nonce values.Nonce, previousHash values.Hash, transactions []values.Transaction, difficulty values.MiningDifficulty) (bool, error) {
	guessBlock := entity.NewBlockFromValues(nonce, previousHash, transactions)
	hash, err := guessBlock.Hash()
	if err != nil {
		return false, err
	}
	guessHashStr := hash.String()
	valid := guessHashStr[:difficulty.Difficulty()] == difficulty.Zeros()
	if valid {
		p.logger.Debug("Valid Nonce Found", zap.Int("nonce", nonce.Int()), zap.String("hash", guessHashStr))
	}

	return valid, nil
}

func NewProofOfWork(logger *zap.Logger, miningDifficulty int) (repositories.ProofOfWork, error) {

	difficulty, err := values.MiningDifficultyFromInt(miningDifficulty)
	if err != nil {
		return nil, err
	}

	logger.Debug("New proof of work", zap.Int("difficulty", difficulty.Difficulty()), zap.String("zeros", difficulty.Zeros()))

	return &ProofOfWork{
		logger:     logger,
		difficulty: difficulty,
	}, nil
}
