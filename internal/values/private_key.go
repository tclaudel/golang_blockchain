package values

import (
	"crypto/ecdsa"
	"crypto/rand"
)

type PrivateKey struct {
	privateKey ecdsa.PrivateKey
}

func NewPrivateKey(privateKey ecdsa.PrivateKey) PrivateKey {
	return PrivateKey{
		privateKey: privateKey,
	}
}

func (p PrivateKey) Sign(data []byte) (Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, &p.privateKey, data)
	if err != nil {
		return Signature{}, err
	}

	return Signature{
		r: r,
		s: s,
	}, err
}
