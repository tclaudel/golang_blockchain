package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/tclaudel/golang_blockchain/config"
	"github.com/tclaudel/golang_blockchain/internal/entity"
	"github.com/tclaudel/golang_blockchain/internal/log"
	"github.com/tclaudel/golang_blockchain/internal/values"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/http"
	"github.com/tclaudel/golang_blockchain/pkg/repositories"
	"github.com/tclaudel/golang_blockchain/pkg/repositories/badger"
	"github.com/tclaudel/golang_blockchain/pkg/repositories/fs"
	"github.com/tclaudel/golang_blockchain/pkg/repositories/memory"
	"go.uber.org/zap"
)

const serverTimeout = 5 * time.Second

func main() {
	cfg := config.NewServerConfig()

	logger, err := log.New(cfg.Log.Format, cfg.Log.Level)
	if err != nil {
		panic(err)
	}
	logger.Debug("configuration", zap.Any("config", cfg))

	repositories, err := selectRepositories(logger, cfg.Repositories)
	if err != nil {
		panic(err)
	}
	defer repositories.Close()

	nodeWallet := values.NewWallet(values.IdentifierFromString("blockchain"))
	if err := repositories.Wallet().Save(nodeWallet); err != nil {
		panic(err)
	}

	ownerWallet := values.NewWallet(
		values.IdentifierFromString(cfg.Repositories.Wallet.Name))
	repositories.Wallet().Save(ownerWallet)
	blockchainNode, err := entity.NewBlockchainNode(
		nodeWallet,
		logger,
		ownerWallet,
		values.AmountFromFloat64(cfg.Repositories.Wallet.InitialAmount),
		repositories,
	)

	httpServer := http.NewServer(cfg, logger, blockchainNode, repositories)

	var stop = make(chan error, 1)
	go httpServer.Start(cfg, logger, stop)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
		stop <- fmt.Errorf("received Interrupt signal")
	}()
	err = <-stop

	logger.With(zap.Error(err)).Info("Shutting down services")
	stopCtx, cancel := context.WithTimeout(context.Background(), serverTimeout)
	defer cancel()
	httpServer.Shutdown(stopCtx, logger)
}

type Repositories struct {
	proofOfWork repositories.ProofOfWork
	wallet      repositories.Wallet
	blockchain  repositories.Blockchain
}

func selectRepositories(logger *zap.Logger, cfgRepo config.NodeRepositories) (repositories.ServerRepositories, error) {
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

	switch cfgRepo.Wallet.Type {
	case "filesystem":
		repos.wallet = fs.NewWalletRepository(logger, cfgRepo.Wallet)
	default:
		logger.Error("unknown wallet type", zap.String("type", cfgRepo.Wallet.Type))
		return nil, errors.New("unknown wallet type")
	}

	switch cfgRepo.Blockchain.Type {
	case "badger":
		repos.blockchain, err = badger.NewBlockchainRepository(logger, cfgRepo.Blockchain)
		if err != nil {
			return nil, err
		}
	case "memory":
		repos.blockchain, err = memory.NewBlockchainRepository(logger, cfgRepo.Blockchain)
		if err != nil {
			return nil, err
		}

	default:
		logger.Error("unknown blockchain type", zap.String("type", cfgRepo.Blockchain.Type))
		return nil, errors.New("unknown blockchain type")
	}

	return repos, nil
}

func (r *Repositories) ProofOfWork() repositories.ProofOfWork {
	return r.proofOfWork
}

func (r *Repositories) Wallet() repositories.Wallet {
	return r.wallet
}

func (r *Repositories) Blockchain() repositories.Blockchain {
	return r.blockchain
}

func (r *Repositories) Close() error {
	r.blockchain.Close()
	return nil
}
