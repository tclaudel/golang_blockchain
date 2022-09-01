package entity

import (
	"github.com/tclaudel/golang_blockchain/internal/values"
	"github.com/tclaudel/golang_blockchain/pkg/entity"
	"go.uber.org/zap/zapcore"
)

type Block struct {
	timestamp    values.Timestamp
	nonce        values.Nonce
	previousHash values.Hash
	transactions []values.Transaction
	hash         values.Hash
}

func (b *Block) Hash() values.Hash {
	return b.hash
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

func (b *Block) Transactions() ([]values.Transaction, error) {
	return b.transactions, nil
}

var Genesis = &Block{
	nonce:        values.GenesisNonce,
	previousHash: values.GenesisHash(),
	timestamp:    values.TimestampNow(),
	transactions: nil,
	hash:         values.GenesisHash(),
}

func (b *Block) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if err := enc.AddObject("hash", b.hash); err != nil {
		return err
	}

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

func (b *Block) FromEntity(block entity.Block) error {
	var transactions []values.Transaction
	txs, err := block.Transactions()
	if err != nil {
		return err
	}
	for _, txEntity := range txs {
		pkKey, err := txEntity.SenderPublicKey()
		if err != nil {
			return err
		}

		sig, err := txEntity.Signature()
		if err != nil {
			return err
		}

		transactions = append(transactions, values.TransactionFromValues(
			pkKey,
			txEntity.SenderAddress(),
			txEntity.RecipientAddress(),
			txEntity.Amount(),
			sig,
		))
	}

	*b = Block{
		timestamp:    block.Timestamp(),
		nonce:        block.Nonce(),
		previousHash: block.PreviousHash(),
		transactions: transactions,
		hash:         block.Hash(),
	}

	return nil
}
