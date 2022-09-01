package types

import (
	"time"
)

type Block struct {
	Hash         string         `json:"hash"`
	Timestamp    time.Time      `json:"timestamp"`
	Nonce        int            `json:"nonce"`
	PreviousHash string         `json:"previous_hash"`
	Transactions []*Transaction `json:"transactions"`
}
