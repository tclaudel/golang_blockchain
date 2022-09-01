package clicfg

import (
	"errors"

	"github.com/tclaudel/golang_blockchain/config"
	"github.com/tclaudel/golang_blockchain/internal/log"
	"github.com/tclaudel/golang_blockchain/pkg/repositories"
	"github.com/tclaudel/golang_blockchain/pkg/repositories/fs"
	"go.uber.org/zap"
)

var (
	Cfg          *config.Config
	Logger       *zap.Logger
	Repositories repositories.CliRepositories
)

func InitConfig() {
	var err error

	Cfg = config.NewCliConfig()
	Logger, err = log.New(Cfg.Log.Format, Cfg.Log.Level)
	if err != nil {
		panic(err)
	}

	Repositories, err = selectCliRepositories(Logger, Cfg.Repositories)
	if err != nil {
		panic(err)
	}
}

type RepositoriesHolder struct {
	proofOfWork repositories.ProofOfWork
	wallet      repositories.Wallet
}

func selectCliRepositories(logger *zap.Logger, cfgRepo config.CliRepositories) (repositories.CliRepositories, error) {
	var (
		err   error
		repos = new(RepositoriesHolder)
	)

	_ = err

	switch cfgRepo.Wallet.Type {
	case "filesystem":
		repos.wallet = fs.NewWalletRepository(logger, cfgRepo.Wallet)
	default:
		logger.Error("unknown wallet type", zap.String("type", cfgRepo.Wallet.Type))
		return nil, errors.New("unknown wallet type")

	}

	return repos, nil
}
func (r *RepositoriesHolder) Wallet() repositories.Wallet {
	return r.wallet
}
