package clicfg

import (
	"errors"

	"github.com/tclaudel/golang_blockchain/config"
	"github.com/tclaudel/golang_blockchain/internal/log"
	"github.com/tclaudel/golang_blockchain/pkg/repositories"
	"github.com/tclaudel/golang_blockchain/pkg/repositories/fs"
	"github.com/tclaudel/golang_blockchain/pkg/repositories/memory"
	"go.uber.org/zap"
)

var (
	Cfg          *config.Config
	Logger       *zap.Logger
	Repositories repositories.Repositories
)

func InitConfig() {
	var err error

	Cfg = config.NewConfig()
	Logger, err = log.New(Cfg.Log.Format, Cfg.Log.Level)
	if err != nil {
		panic(err)
	}

	Repositories, err = selectRepositories(Logger, Cfg.Repositories)
	if err != nil {
		panic(err)
	}
}

type RepositoriesHolder struct {
	proofOfWork repositories.ProofOfWork
	wallet      repositories.Wallet
}

func (r *RepositoriesHolder) Blockchain() repositories.Blockchain {
	//TODO implement me
	panic("implement me")
}

func (r *RepositoriesHolder) Close() error {
	//TODO implement me
	panic("implement me")
}

func selectRepositories(logger *zap.Logger, cfgRepo config.Repositories) (repositories.Repositories, error) {
	var (
		err   error
		repos = new(RepositoriesHolder)
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

	switch cfgRepo.Wallet.Type {
	case "filesystem":
		repos.wallet = fs.NewWalletRepository(logger, cfgRepo.Wallet)
	default:
		logger.Error("unknown wallet type", zap.String("type", cfgRepo.Wallet.Type))
		return nil, errors.New("unknown wallet type")

	}

	return repos, nil
}

func (r *RepositoriesHolder) ProofOfWork() repositories.ProofOfWork {
	return r.proofOfWork
}

func (r *RepositoriesHolder) Wallet() repositories.Wallet {
	return r.wallet
}
