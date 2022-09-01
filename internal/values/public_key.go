package values

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"math/big"
)

type PublicKey struct {
	publicKey ecdsa.PublicKey
}

func NewPublicKey(publicKey ecdsa.PublicKey) PublicKey {
	return PublicKey{
		publicKey: publicKey,
	}
}

func PublicKeyFromStrings(x, y string) (PublicKey, error) {
	var (
		xInt = new(big.Int)
		yInt = new(big.Int)
	)

	if _, ok := xInt.SetString(x, 10); !ok {
		return PublicKey{}, fmt.Errorf("invalid x")
	}
	if _, ok := yInt.SetString(y, 10); !ok {
		return PublicKey{}, fmt.Errorf("invalid y")
	}

	var publicKey ecdsa.PublicKey
	publicKey.Curve = elliptic.P256()
	publicKey.X = xInt
	publicKey.Y = yInt

	return PublicKey{
		publicKey: publicKey,
	}, nil
}

func (pk PublicKey) Strings() (x string, y string) {
	return pk.publicKey.X.String(), pk.publicKey.Y.String()
}

func (pk PublicKey) Verify(data []byte, signature Signature) bool {
	return ecdsa.Verify(&pk.publicKey, data, signature.r, signature.s)
}
