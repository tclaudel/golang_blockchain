package entity

import "github.com/tclaudel/golang_blockchain/internal/values"

type Transaction interface {
	SenderPublicKey() (values.PublicKey, error)
	SenderAddress() values.Address
	RecipientAddress() values.Address
	Amount() values.Amount
	Signature() (values.Signature, error)
}
