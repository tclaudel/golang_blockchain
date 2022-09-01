package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tclaudel/golang_blockchain/pkg/interfaces/http/errors"
)

func (s *Server) GetBlockchain(w http.ResponseWriter, r *http.Request) {
	logger := GetRequestLogger(r)
	_, err := GetRequestID(r)
	if err != nil {
		errors.ErrNilRequestID.Write(logger, w, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var blocks []Block

	bs, err := s.blockchainNode.Blocks()
	if err != nil {
		errors.ErrMarshalingJSON.Write(logger, w, http.StatusInternalServerError)
		return
	}
	for _, block := range bs {
		var transactions []Transaction
		txs, err := block.Transactions()
		if err != nil {
			errors.ErrMarshalingJSON.Write(logger, w, http.StatusInternalServerError)
			return
		}
		for _, transaction := range txs {
			transactions = append(transactions, Transaction{
				Sender:    transaction.SenderAddress().String(),
				Recipient: transaction.RecipientAddress().String(),
				Amount:    transaction.Amount().Float64(),
			})
		}

		blocks = append(blocks, Block{
			Hash:         block.Hash().String(),
			Timestamp:    block.Timestamp().Time().Format(time.RFC3339),
			Nonce:        block.Nonce().Int(),
			PreviousHash: block.PreviousHash().String(),
			Transactions: transactions,
		})
	}
	data, err := json.MarshalIndent(blocks, "", "  ")
	if err != nil {
		errors.ErrMarshalingJSON.Write(logger, w, http.StatusInternalServerError)
		return
	}
	w.Write(data)
	return
}
