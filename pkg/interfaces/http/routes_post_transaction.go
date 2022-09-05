package http

import (
	"encoding/json"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/tclaudel/golang_blockchain/internal/values"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/http/errors"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/http/rest"
)

func (s *Server) PostTransaction(w http.ResponseWriter, r *http.Request) {
	logger := GetRequestLogger(r)
	_, err := GetRequestID(r)
	if err != nil {
		errors.ErrNilRequestID.Write(logger, w, http.StatusBadRequest)
		return
	}

	var transaction rest.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		errors.ErrMarshalingJSON.Write(logger, w, http.StatusBadRequest)
		return
	}

	spew.Dump(transaction)

	pk, err := values.PublicKeyFromString(transaction.SenderPublicKey)
	if err != nil {
		errors.ErrInvalidPublicKey.Write(logger, w, http.StatusBadRequest)
		return
	}

	sig, err := values.SignatureFromString(transaction.Signature)
	if err != nil {
		errors.ErrInvalidSignature.Write(logger, w, http.StatusBadRequest)
		return
	}

	tx, err := s.blockchainNode.AppendTransaction(
		values.TimestampFromTime(transaction.Timestamp),
		pk,
		values.AddressFromString(transaction.SenderAddress),
		values.AddressFromString(transaction.RecipientAddress),
		values.AmountFromFloat64(transaction.Amount),
		sig,
	)
	if err != nil {
		errors.ErrAppendingTransaction.Write(logger, w, http.StatusInternalServerError)
		return
	}

	data, err := json.MarshalIndent(rest.Transaction{
		Amount:    tx.Amount().Float64(),
		Recipient: tx.RecipientAddress().String(),
		Sender:    tx.SenderAddress().String(),
	}, "", "  ")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}
