package main

import (
	"errors"

	"github.com/tclaudel/golang_blockchain/config"
	"github.com/tclaudel/golang_blockchain/internal/entity"
	"github.com/tclaudel/golang_blockchain/internal/log"
	"github.com/tclaudel/golang_blockchain/internal/values"
	"github.com/tclaudel/golang_blockchain/pkg/repositories"
	"github.com/tclaudel/golang_blockchain/pkg/repositories/memory"
	"go.uber.org/zap"
)

func main() {
	cfg := config.NewConfig()

	logger, err := log.New(cfg.Log.Format, cfg.Log.Level)
	if err != nil {
		panic(err)
	}
	logger.Debug("configuration", zap.Any("config", cfg))

	repositories, err := selectRepositories(logger, cfg.Repositories)
	if err != nil {
		panic(err)
	}

	ownerWallet := values.NewWallet()
	blockchainNode := entity.NewBlockchainNode(
		logger,
		ownerWallet,
		values.AmountFromFloat64(cfg.Blockchain.Reward),
		repositories)

	userAWallet := values.NewWallet()
	userBWallet := values.NewWallet()

	err = blockchainNode.AppendTransaction(userAWallet, ownerWallet.Address(), values.AmountFromFloat64(10))
	if err != nil {
		panic(err)
	}

	err = blockchainNode.AppendTransaction(userBWallet, ownerWallet.Address(), values.AmountFromFloat64(10))
	if err != nil {
		panic(err)
	}

	err = blockchainNode.Commit()
	if err != nil {
		panic(err)
	}

	ownerAmount := blockchainNode.CalculateTotalAmount(ownerWallet.Address())
	userAAmount := blockchainNode.CalculateTotalAmount(userAWallet.Address())
	userBAmount := blockchainNode.CalculateTotalAmount(userBWallet.Address())
	logger.Info("total amount of owner", zap.Float64("amount", ownerAmount))
	logger.Info("total amount of user A", zap.Float64("amount", userAAmount))
	logger.Info("total amount of user B", zap.Float64("amount", userBAmount))

}

type Repositories struct {
	proofOfWork repositories.ProofOfWork
}

func selectRepositories(logger *zap.Logger, cfgRepo config.Repositories) (repositories.Repositories, error) {
	var (
		err   error
		repos = new(Repositories)
	)

	switch cfgRepo.ProofOfWork.Type {
	case "memory":
		repos.proofOfWork, err = memory.NewProofOfWork(logger, cfgRepo.ProofOfWork.Difficulty)
		if err != nil {
			return nil, err
		}
	default:
		logger.Error("unknown proof of work type", zap.String("type", cfgRepo.ProofOfWork.Type))
		return nil, errors.New("unknown proof of work type")
	}

	return repos, nil
}

func (r *Repositories) ProofOfWork() repositories.ProofOfWork {
	return r.proofOfWork
}
