package values

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
)

type PrivateKey struct {
	privateKey ecdsa.PrivateKey
}

func NewPrivateKey(privateKey ecdsa.PrivateKey) PrivateKey {
	return PrivateKey{
		privateKey: privateKey,
	}
}

func (p PrivateKey) String() string {
	return fmt.Sprintf("%s", p.privateKey.D.String())
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
