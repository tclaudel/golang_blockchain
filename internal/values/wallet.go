package values

import "encoding/json"

const WalletPath = "./wallet"

type Wallet struct {
	identifier Identifier
	KeyPair
	address Address
}

func NewWallet(identifier Identifier) Wallet {
	keyPair := GenerateKeyPair()

	wallet := Wallet{
		identifier: IdentifierFromString(func() string {
			if identifier.String() == "" {
				return keyPair.Address().String()
			}
			return identifier.String()
		}()),
		KeyPair: keyPair,
		address: keyPair.Address(),
	}

	return wallet
}

func (w Wallet) Identifier() Identifier {
	return w.identifier
}

func (w Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Address string  `json:"address"`
		KeyPair KeyPair `json:"key_pair"`
	}{
		Address: w.address.String(),
		KeyPair: w.KeyPair,
	})
}

func (w *Wallet) UnmarshalJSON(data []byte) error {
	var v struct {
		Address string  `json:"address"`
		KeyPair KeyPair `json:"key_pair"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	w.address = AddressFromString(v.Address)
	w.KeyPair = v.KeyPair

	return nil
}

func (w Wallet) Address() Address {
	return w.PublicKey.Address()
}
