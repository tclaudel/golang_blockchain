package values

import (
	"bytes"
	"crypto/sha256"

	"github.com/mr-tron/base58"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/ripemd160"
)

const version = "0x00"

type Address struct {
	address string
}

func GenerateAddress(keyPair KeyPair) Address {
	// 1. SHA256 of public key
	h1 := sha256.New()
	h1.Write(keyPair.publicKey.X.Bytes())
	h1.Write(keyPair.publicKey.Y.Bytes())
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

func (a Address) String() string {
	return a.address
}

func (a Address) Equal(other Address) bool {
	return a.String() == other.String()
}

func (a Address) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("address", a.String())
	return nil
}
