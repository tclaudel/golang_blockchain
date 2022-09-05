package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/http/errors"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/http/rest"
)

func (s *Server) CommitTransactions(w http.ResponseWriter, r *http.Request) {
	logger := GetRequestLogger(r)
	_, err := GetRequestID(r)
	if err != nil {
		errors.ErrNilRequestID.Write(logger, w, http.StatusBadRequest)
		return
	}

	spew.Dump("CommitTransactions")
	blockEntity, err := s.blockchainNode.Commit()
	if err != nil {
		errors.ErrMarshalingJSON.Write(logger, w, http.StatusInternalServerError)
		return
	}

	txs, err := blockEntity.Transactions()
	if err != nil {
		errors.ErrMarshalingJSON.Write(logger, w, http.StatusInternalServerError)
		return
	}

	var transactions = make([]rest.Transaction, len(txs))
	for i, tx := range txs {
		transactions[i] = rest.Transaction{
			Amount:    tx.Amount().Float64(),
			Recipient: tx.RecipientAddress().String(),
			Sender:    tx.SenderAddress().String(),
		}
	}

	block := rest.Block{
		Hash:         blockEntity.Hash().String(),
		Nonce:        blockEntity.Nonce().Int(),
		PreviousHash: blockEntity.PreviousHash().String(),
		Timestamp:    blockEntity.Timestamp().Time().Format(time.RFC3339),
		Transactions: transactions,
	}

	data, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		errors.ErrMarshalingJSON.Write(logger, w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}
