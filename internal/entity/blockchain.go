package entity

import (
	"fmt"

	"github.com/tclaudel/golang_blockchain/internal/values"
	"github.com/tclaudel/golang_blockchain/pkg/repositories"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const reward = 1.0

type BlockchainNode struct {
	transactionsPool *TransactionPool
	chain            []*Block
	logger           *zap.Logger
	repositories     repositories.Repositories
	owner            values.Wallet
	miningReward     values.Amount
}

var blockchainWallet = values.NewWallet()

func NewBlockchainNode(logger *zap.Logger, wallet values.Wallet, repositories repositories.Repositories) *BlockchainNode {
	bc := &BlockchainNode{
		transactionsPool: NewTransactionPool(),
		chain:            []*Block{Genesis},
		logger:           logger,
		repositories:     repositories,
		owner:            wallet,
		miningReward:     values.AmountFromFloat64(reward),
	}

	logger.Info("created new blockchain node", zap.String("owner", wallet.Address().String()))
	return bc
}

func (bc *BlockchainNode) blocks() []*Block {
	return bc.chain
}

func (bc *BlockchainNode) CalculateTotalAmount(address values.Address) float64 {
	var totalAmount float64 // = 20.0
	for _, block := range bc.chain {
		for _, tx := range block.Transactions() {
			if tx.IsSender(address) {
				totalAmount -= tx.Amount().Float64()
			}
			if tx.IsRecipient(address) {
				totalAmount += tx.Amount().Float64()
			}
		}
	}
	return totalAmount
}

func (bc *BlockchainNode) AppendTransaction(senderWallet values.Wallet, recipientAddress values.Address, value values.Amount) error {
	tx, err := values.NewTransaction(senderWallet, recipientAddress, value)
	if err != nil {
		return err
	}

	verified, err := tx.Verify(senderWallet.PublicKey)
	if !verified {
		bc.logger.Error("transaction verification failed", zap.Error(err))
		return err
	}

	if bc.CalculateTotalAmount(senderWallet.Address()) < value.Float64() {
		bc.logger.Error("insufficient funds", zap.Error(err))
		return nil
	}

	bc.transactionsPool.Append(tx)
	return nil
}

func (bc *BlockchainNode) Commit() error {
	if bc.transactionsPool.Len() == 0 {
		return fmt.Errorf("no transactions to commit")
	}

	tx, err := values.NewTransaction(blockchainWallet, bc.owner.Address(), bc.miningReward)
	if err != nil {
		return err
	}

	bc.transactionsPool.Append(tx)

	previousHash, err := bc.lastBlock().Hash()
	if err != nil {
		bc.logger.Error("failed to calculate previous hash", zap.Error(err))
		return err
	}

	nonce, err := bc.repositories.ProofOfWork().Mine(previousHash, bc.transactionsPool.Transactions())
	if err != nil {
		bc.logger.Error("failed to mine block", zap.Error(err))
		return err
	}

	newBlock := NewBlockFromValues(nonce, previousHash, bc.transactionsPool.Transactions())
	bc.logger.Debug("created new block", zap.Object("block", newBlock))
	bc.chain = append(bc.chain, newBlock)
	bc.transactionsPool.Flush()
	return nil
}

func (bc *BlockchainNode) lastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *BlockchainNode) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if err := enc.AddObject("transactions_pool", bc.transactionsPool); err != nil {
		return err
	}

	if err := enc.AddArray("chain", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
		for _, block := range bc.chain {
			if err := enc.AppendObject(block); err != nil {
				return err
			}
		}
		return nil
	})); err != nil {
		return err
	}

	return nil
}
