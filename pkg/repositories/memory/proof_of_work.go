package memory

import (
	"crypto/sha256"
	"encoding/json"
	"time"

	"github.com/tclaudel/golang_blockchain/internal/values"
	"github.com/tclaudel/golang_blockchain/pkg/entity"
	"github.com/tclaudel/golang_blockchain/pkg/repositories"
	"go.uber.org/zap"
)

type Transaction struct {
	senderPublicKey  values.PublicKey
	senderAddress    values.Address
	recipientAddress values.Address
	amount           values.Amount
	signature        values.Signature
	timestamp        values.Timestamp
}

func (t Transaction) SenderPublicKey() (values.PublicKey, error) {
	return t.senderPublicKey, nil
}

func (t Transaction) SenderAddress() values.Address {
	return t.senderAddress
}

func (t Transaction) RecipientAddress() values.Address {
	return t.recipientAddress
}

func (t Transaction) Amount() values.Amount {
	return t.amount
}

func (t Transaction) Signature() (values.Signature, error) {
	return t.signature, nil
}

func (t Transaction) Timestamp() (values.Timestamp, error) {
	return values.TimestampNow(), nil
}

type Block struct {
	timestamp    values.Timestamp
	nonce        values.Nonce
	previousHash values.Hash
	transactions []Transaction
	hash         values.Hash
}

func (b *Block) Timestamp() values.Timestamp {
	return b.timestamp
}

func (b *Block) Nonce() values.Nonce {
	return b.nonce
}

func (b *Block) PreviousHash() values.Hash {
	return b.previousHash
}

func (b *Block) Hash() values.Hash {
	return b.hash
}

func (b *Block) Transactions() ([]values.Transaction, error) {
	var transactions = make([]values.Transaction, len(b.transactions))
	for i, transaction := range b.transactions {
		t, err := transaction.Timestamp()
		if err != nil {
			return nil, err
		}

		transactions[i] = values.TransactionFromValues(
			t,
			transaction.senderPublicKey,
			transaction.senderAddress,
			transaction.recipientAddress,
			transaction.amount,
			transaction.signature)
	}

	return transactions, nil
}

func (b *Block) Entity() entity.Block {
	return b
}

func (b *Block) IncreaseNonce() {
	b.nonce = b.nonce.Increase()
}

func (b Block) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(
		struct {
			Nonce        int           `json:"nonce"`
			PreviousHash []byte        `json:"previous_hash"`
			Timestamp    time.Time     `json:"timestamp"`
			Transactions []Transaction `json:"transactions"`
		}{
			Nonce:        b.nonce.Int(),
			PreviousHash: b.previousHash.Bytes(),
			Timestamp:    b.timestamp.Time(),
			Transactions: b.transactions,
		},
	)

	return data, err
}

type ProofOfWork struct {
	logger     *zap.Logger
	difficulty values.MiningDifficulty
}

func (pow *ProofOfWork) Mine(previousHash values.Hash, transactions []entity.Transaction) (entity.Block, error) {
	var blockTransactions = make([]Transaction, len(transactions))
	for i, transaction := range transactions {
		pk, err := transaction.SenderPublicKey()
		if err != nil {
			pow.logger.Error("Failed to get sender public key", zap.Error(err))
			return nil, err
		}

		sig, err := transaction.Signature()
		if err != nil {
			pow.logger.Error("Failed to get signature", zap.Error(err))
			return nil, err
		}

		blockTransactions[i] = Transaction{
			senderPublicKey:  pk,
			senderAddress:    transaction.SenderAddress(),
			recipientAddress: transaction.RecipientAddress(),
			amount:           transaction.Amount(),
			signature:        sig,
		}
	}

	block := &Block{
		timestamp:    values.TimestampFromTime(time.Now()),
		nonce:        values.InitNonce(),
		previousHash: previousHash,
		transactions: blockTransactions,
	}

	for {
		valid, err := pow.validate(block, pow.difficulty)
		if err != nil {
			pow.logger.Error("Failed to validate block", zap.Error(err))
			return nil, err
		}

		if !valid {
			block.IncreaseNonce()
			continue
		}

		break
	}

	pow.logger.Debug("Hash", zap.Int("nonce", block.nonce.Int()))

	return block, nil
}

func (pow *ProofOfWork) validate(block *Block, difficulty values.MiningDifficulty) (bool, error) {
	data, err := json.Marshal(block)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(data)
	block.hash = values.HashFromBytes(hash[:])

	return block.hash.String()[:difficulty.Difficulty()] == difficulty.Zeros(), nil
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
