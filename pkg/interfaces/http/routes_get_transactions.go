package http

import (
	"encoding/json"
	"net/http"

	"github.com/tclaudel/golang_blockchain/pkg/interfaces/http/errors"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/http/rest"
)

func (s *Server) GetTransactionPool(w http.ResponseWriter, r *http.Request) {
	logger := GetRequestLogger(r)
	_, err := GetRequestID(r)
	if err != nil {
		errors.ErrNilRequestID.Write(logger, w, http.StatusBadRequest)
		return
	}

	transactionPool := s.blockchainNode.TransactionPool()

	var transactions = make([]rest.Transaction, len(transactionPool))
	for i, tx := range transactionPool {
		transactions[i] = rest.Transaction{
			Amount:    tx.Amount().Float64(),
			Recipient: tx.RecipientAddress().String(),
			Sender:    tx.SenderAddress().String(),
		}
	}

	data, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		errors.ErrMarshalingJSON.Write(logger, w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}