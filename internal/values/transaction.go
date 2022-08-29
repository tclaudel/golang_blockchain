package values

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"go.uber.org/zap/zapcore"
)

type Transaction struct {
	senderPublicKey  PublicKey
	senderAddress    Address
	recipientAddress Address
	value            Amount
	signature        Signature
}

func NewTransaction(senderWallet Wallet, recipientAddress Address, value Amount) (Transaction, error) {
	var err error

	transaction := Transaction{
		senderAddress:    senderWallet.Address(),
		recipientAddress: recipientAddress,
		value:            value,
	}

	txHash, err := transaction.Hash()
	if err != nil {
		return Transaction{}, err
	}

	transaction.signature, err = senderWallet.Sign(txHash)
	if err != nil {
		return Transaction{}, err
	}

	return transaction, err
}

func (t Transaction) Verify(publicKey PublicKey) (bool, error) {
	hash, err := t.Hash()
	if err != nil {
		return false, err
	}

	verified := publicKey.Verify(hash, t.signature)
	if !verified {
		return false, fmt.Errorf("transaction verification failed")
	}

	return true, nil
}

func (t Transaction) Hash() ([]byte, error) {
	data, err := json.Marshal(struct {
		SenderAddress    string  `json:"sender_address"`
		RecipientAddress string  `json:"recipient_address"`
		Value            float64 `json:"value"`
	}{
		SenderAddress:    t.senderAddress.String(),
		RecipientAddress: t.recipientAddress.String(),
		Value:            t.value.Float64(),
	})

	hash := sha256.Sum256(data)

	return hash[:], err
}

func (t Transaction) IsSender(address Address) bool {
	return t.senderAddress.Equal(address)
}

func (t Transaction) IsRecipient(address Address) bool {
	return t.recipientAddress.Equal(address)
}

func (t Transaction) Amount() Amount {
	return t.value
}

func (t Transaction) String() string {
	return fmt.Sprintf("%s -> %s: %s", t.senderAddress, t.recipientAddress, t.value)
}

func (t Transaction) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if err := enc.AddObject("senderAddress", t.senderAddress); err != nil {
		return err
	}

	if err := enc.AddObject("recipientAddress", t.recipientAddress); err != nil {
		return err
	}

	if err := enc.AddObject("value", t.value); err != nil {
		return err
	}
	return nil
}
