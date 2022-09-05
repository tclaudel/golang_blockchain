package entity

import (
	"fmt"
	"sync"

	"github.com/tclaudel/golang_blockchain/internal/values"
	"go.uber.org/zap/zapcore"
)

type TransactionPool struct {
	sync.Mutex
	transactions []values.Transaction
}

func NewTransactionPool() *TransactionPool {
	return new(TransactionPool)
}

func (tp *TransactionPool) Len() int {
	return len(tp.transactions)
}

func (tp *TransactionPool) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for i, transaction := range tp.transactions {
		if err := enc.AddObject(fmt.Sprintf("transaction[%3d]", i), transaction); err != nil {
			return err
		}
	}

	return nil
}

func (tp *TransactionPool) Export() []values.Transaction {
	tp.Lock()
	dst := make([]values.Transaction, len(tp.transactions))
	copy(dst, tp.transactions)

	tp.transactions = nil
	tp.Unlock()
	return dst
}

func (tp *TransactionPool) Transactions() []values.Transaction {
	return tp.transactions
}

func (tp *TransactionPool) Append(transaction values.Transaction) {
	tp.Lock()
	tp.transactions = append(tp.transactions, transaction)
	tp.Unlock()
}
