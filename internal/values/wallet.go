package values

type Wallet struct {
	KeyPair
	address Address
}

func NewWallet() Wallet {
	keyPair := GenerateKeyPair()
	address := GenerateAddress(keyPair)

	return Wallet{
		KeyPair: keyPair,
		address: address,
	}
}

func (w Wallet) Address() Address {
	return w.address
}
