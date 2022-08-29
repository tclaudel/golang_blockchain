package values

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
)

type KeyPair struct {
	PublicKey
	PrivateKey
}

func GenerateKeyPair() KeyPair {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	return KeyPair{
		PublicKey:  NewPublicKey(privateKey.PublicKey),
		PrivateKey: NewPrivateKey(*privateKey),
	}
}

func (k KeyPair) String() string {
	return fmt.Sprintf("%s %s", k.PublicKey, k.PrivateKey)
}
