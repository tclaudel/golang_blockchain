package entity

import (
	"encoding/json"
	"time"

	"github.com/tclaudel/golang_blockchain/internal/values"
	"go.uber.org/zap/zapcore"
)

type Block struct {
	timestamp    values.Timestamp
	nonce        values.Nonce
	previousHash values.Hash
	transactions []values.Transaction
}

var Genesis = &Block{
	nonce:        values.GenesisNonce,
	previousHash: values.GenesisHash,
	timestamp:    values.TimestampNow(),
	transactions: nil,
}

func NewBlockFromValues(
	nonce values.Nonce,
	previousHash values.Hash,
	transactions []values.Transaction,
) *Block {
	block := Block{
		nonce:        nonce,
		previousHash: previousHash,
		timestamp:    values.TimestampNow(),
		transactions: transactions,
	}

	return &block
}

func (b *Block) Transactions() []values.Transaction {
	return b.transactions
}

func (b *Block) Hash() (values.Hash, error) {
	data, err := json.Marshal(b)
	if err != nil {
		return values.Hash{}, err
	}
	hash := values.HashFromBytes(data)

	return hash, nil
}

func (b *Block) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(
		struct {
			Nonce        int                  `json:"nonce"`
			PreviousHash []byte               `json:"previous_hash"`
			Timestamp    time.Time            `json:"timestamp"`
			Transactions []values.Transaction `json:"transactions"`
		}{
			Nonce:        b.nonce.Int(),
			PreviousHash: b.previousHash.Bytes(),
			Timestamp:    b.timestamp.Time(),
			Transactions: b.transactions,
		},
	)

	return data, err
}

func (b *Block) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if err := enc.AddObject("nonce", b.nonce); err != nil {
		return err
	}

	if err := enc.AddObject("previous_hash", b.previousHash); err != nil {
		return err
	}

	if err := enc.AddObject("timestamp", b.timestamp); err != nil {
		return err
	}

	if err := enc.AddArray("transaction", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
		for _, tx := range b.transactions {
			if err := enc.AppendObject(tx); err != nil {
				return err
			}
		}

		return nil
	})); err != nil {
		return err
	}

	return nil
}
