package values

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"

	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
)

type PublicKey struct {
	publicKey ecdsa.PublicKey
	address   Address
}

func NewPublicKey(publicKey ecdsa.PublicKey) PublicKey {
	return PublicKey{
		publicKey: publicKey,
		address:   Address{},
	}
}

func (pk PublicKey) Address() Address {
	if pk.address.address != "" {
		return pk.address
	}
	// 1. SHA256 of public key
	h1 := sha256.New()
	h1.Write(pk.publicKey.X.Bytes())
	h1.Write(pk.publicKey.Y.Bytes())
	digest1 := h1.Sum(nil)

	// 2. RipeMD160 of hashed public key
	h2 := ripemd160.New()
	h2.Write(digest1)
	digest2 := h2.Sum(nil)

	// 3. Add version bytes in front of public key
	digest3 := bytes.Join([][]byte{[]byte(version), digest2}, []byte{})

	// 4. Rehash public key with SHA256
	h4 := sha256.New()
	h4.Write(digest3)
	digest4 := h4.Sum(nil)

	// 5. Rehash public key with SHA256
	h5 := sha256.New()
	h5.Write(digest3)
	digest5 := h5.Sum(nil)

	// 6. Take the first 4 bytes of the second SHA256 hash. This is the address checksum
	checksum := digest5[:4]

	// 7. Add the checksum at the end of extended RIPEMD160 hash
	digest6 := bytes.Join([][]byte{digest4, checksum}, []byte{})

	// 8. Base58 encode
	address := base58.Encode(digest6)

	return Address{address: address}

}

const sep = ","

func PublicKeyFromString(data string) (PublicKey, error) {
	fields := strings.Split(data, sep)
	if len(fields) != 2 {
		return PublicKey{}, fmt.Errorf("invalid data")
	}

	return PublicKeyFromStrings(fields[0], fields[1])
}

func (pk PublicKey) String() string {
	return fmt.Sprintf("%s%s%s", pk.publicKey.X.String(), sep, pk.publicKey.Y.String())
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
