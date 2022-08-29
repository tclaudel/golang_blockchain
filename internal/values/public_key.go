package values

import "crypto/ecdsa"

type PublicKey struct {
	publicKey ecdsa.PublicKey
}

func NewPublicKey(publicKey ecdsa.PublicKey) PublicKey {
	return PublicKey{
		publicKey: publicKey,
	}
}

func (p PublicKey) Verify(data []byte, signature Signature) bool {
	return ecdsa.Verify(&p.publicKey, data, signature.r, signature.s)
}
