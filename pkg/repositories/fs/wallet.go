package fs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/tclaudel/golang_blockchain/config"
	"github.com/tclaudel/golang_blockchain/internal/values"
	"github.com/tclaudel/golang_blockchain/pkg/repositories"
	"go.uber.org/zap"
)

type WalletRepository struct {
	logger *zap.Logger
	path   string
}

type Wallet struct {
	privateKey string
	address    string
}

func (w WalletRepository) Save(wallet values.Wallet) error {
	data, err := json.MarshalIndent(wallet, "", "  ")
	if err != nil {
		w.logger.Error("failed to serialize wallet", zap.Error(err))
		return err
	}

	path := filepath.Join(w.path, fmt.Sprintf(wallet.Identifier().String()))
	if filepath.Ext(path) != ".json" {
		path += ".json"
	}
	w.logger.Debug("saving wallet", zap.String("path", path))
	err = ioutil.WriteFile(
		path,
		data,
		0644)
	if err != nil {
		w.logger.Error("failed to write wallet", zap.Error(err))
		return err
	}

	return nil
}

func (w WalletRepository) BatchSave(wallets []values.Wallet) error {
	for _, wallet := range wallets {
		err := w.Save(wallet)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w WalletRepository) Load(wallets []values.Wallet) error {
	//TODO implement me
	panic("implement me")
}

func NewWalletRepository(logger *zap.Logger, cfg config.Wallet) repositories.Wallet {
	walletRepository := WalletRepository{
		logger: logger,
		path:   cfg.FS.Path,
	}

	logger.Debug("created wallet repository", zap.String("path", walletRepository.path))
	return walletRepository
}
