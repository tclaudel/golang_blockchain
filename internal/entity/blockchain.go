package entity

import (
	"fmt"
	"sync"
	"time"

	"github.com/tclaudel/golang_blockchain/internal/values"
	"github.com/tclaudel/golang_blockchain/pkg/entity"
	"github.com/tclaudel/golang_blockchain/pkg/repositories"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const reward = 1.0

type BlockchainNode struct {
	nodeWallet          values.Wallet
	transactionsPool    *TransactionPool
	logger              *zap.Logger
	repositories        repositories.ServerRepositories
	owner               values.Wallet
	miningReward        values.Amount
	walletInitialAmount values.Amount
	commitMutex         sync.Mutex
	miningInterval      time.Duration
}

func NewBlockchainNode(nodeWallet values.Wallet, logger *zap.Logger, wallet values.Wallet, walletInitialAmount values.Amount, repositories repositories.ServerRepositories, miningInterval time.Duration) (*BlockchainNode, error) {
	bc := &BlockchainNode{
		nodeWallet:          nodeWallet,
		transactionsPool:    NewTransactionPool(),
		logger:              logger,
		repositories:        repositories,
		owner:               wallet,
		miningReward:        values.AmountFromFloat64(reward),
		walletInitialAmount: walletInitialAmount,
		miningInterval:      miningInterval,
	}

	logger.Info("initialized blockchain node", zap.String("owner", wallet.Address().String()))
	return bc, nil
}

func (bc *BlockchainNode) CronMining() {
	_, err := bc.Commit()
	if err != nil {
		bc.logger.Error("failed to commit", zap.Error(err))
	}
	_ = time.AfterFunc(bc.miningInterval, bc.CronMining)
}

func (bc *BlockchainNode) TransactionPool() []values.Transaction {
	return bc.transactionsPool.Transactions()
}

func (bc *BlockchainNode) Blocks() ([]*Block, error) {
	chain, err := bc.repositories.Blockchain().Get()
	if err != nil {

		return nil, err
	}

	var blockchain = make([]*Block, len(chain.Chain()))
	for i, block := range chain.Chain() {
		var b Block
		err := b.FromEntity(block)
		if err != nil {
			bc.logger.Error("failed to unmarshal block", zap.Error(err))
			return nil, err
		}
		blockchain[i] = &b
	}

	return blockchain, nil
}

func (bc *BlockchainNode) CalculateTotalAmount(address values.Address) (float64, error) {
	var totalAmount = 0.0
	blocks, err := bc.Blocks()
	if err != nil {
		return 0, err
	}

	for _, block := range blocks {
		txs, err := block.Transactions()
		if err != nil {
			bc.logger.Error("failed to get transactions from block", zap.Error(err))
			return 0, err
		}

		for _, tx := range txs {
			if tx.IsSender(address) {
				totalAmount -= tx.Amount().Float64()
			}
			if tx.IsRecipient(address) {
				totalAmount += tx.Amount().Float64()
			}
		}
	}
	totalAmount += bc.walletInitialAmount.Float64()

	return totalAmount, nil
}

func (bc *BlockchainNode) AppendTransaction(
	t values.Timestamp,
	senderPublicKey values.PublicKey,
	senderAddress values.Address,
	recipientAddress values.Address,
	value values.Amount,
	sig values.Signature) (values.Transaction, error) {
	tx, err := values.VerifyTransaction(t, senderPublicKey, senderAddress, recipientAddress, value, sig)
	if err != nil {
		bc.logger.Error("failed to create transaction", zap.Error(err))
		return values.Transaction{}, err
	}

	totalAmount, err := bc.CalculateTotalAmount(senderAddress)
	if err != nil {
		bc.logger.Error("failed to calculate total amount", zap.Error(err))
		return values.Transaction{}, err
	}
	if totalAmount < value.Float64() {
		bc.logger.Error("insufficient funds", zap.Error(err))
		return values.Transaction{}, nil
	}

	bc.transactionsPool.Append(tx)
	return tx, nil
}

func (bc *BlockchainNode) Commit() (*Block, error) {
	bc.commitMutex.Lock()
	defer bc.commitMutex.Unlock()

	txs := bc.transactionsPool.Export()

	if len(txs) == 0 {
		return nil, fmt.Errorf("no transactions to commit")
	}

	tx, err := values.NewTransaction(bc.nodeWallet, bc.owner.Address(), bc.miningReward)
	if err != nil {
		return nil, err
	}

	txs = append(txs, tx)

	lstBlock, err := bc.lastBlock()
	if err != nil {
		return nil, err
	}

	block, err := bc.repositories.ProofOfWork().Mine(lstBlock.Hash(), func() []entity.Transaction {
		var transactions = make([]entity.Transaction, len(txs))
		for i, tx := range txs {
			transactions[i] = tx
		}
		return transactions
	}())
	if err != nil {
		bc.logger.Error("failed to mine block", zap.Error(err))
		return nil, err
	}

	bTxs, err := block.Transactions()
	if err != nil {
		bc.logger.Error("failed to get transactions from block", zap.Error(err))
		return nil, err
	}
	var transactions = make([]values.Transaction, len(bTxs))
	for i, tx := range bTxs {
		pk, err := tx.SenderPublicKey()
		if err != nil {
			bc.logger.Error("failed to get public key", zap.Error(err))
			return nil, err
		}

		sig, err := tx.Signature()
		if err != nil {
			bc.logger.Error("failed to get signature", zap.Error(err))
			return nil, err
		}

		t, err := tx.Timestamp()
		if err != nil {
			bc.logger.Error("failed to get timestamp", zap.Error(err))
			return nil, err
		}

		transactions[i] = values.TransactionFromValues(
			t,
			pk,
			tx.SenderAddress(),
			tx.RecipientAddress(),
			tx.Amount(),
			sig,
		)
	}

	newBlock := &Block{
		timestamp:    block.Timestamp(),
		nonce:        block.Nonce(),
		previousHash: block.PreviousHash(),
		transactions: transactions,
		hash:         block.Hash(),
	}

	err = bc.repositories.Blockchain().Append(newBlock)
	if err != nil {
		bc.logger.Error("failed to append block to blockchain", zap.Error(err))
		return nil, err
	}

	bc.logger.Debug("created new block", zap.Object("block", newBlock))

	return newBlock, nil
}

func (bc *BlockchainNode) lastBlock() (*Block, error) {
	blocks, err := bc.Blocks()
	if err != nil {
		return nil, err
	}

	return blocks[0], nil
}

func (bc *BlockchainNode) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if err := enc.AddObject("transactions_pool", bc.transactionsPool); err != nil {
		return err
	}

	blocks, err := bc.Blocks()
	if err != nil {
		return err
	}

	if err := enc.AddArray("chainCache", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
		for _, block := range blocks {
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
