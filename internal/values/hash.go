package values

import (
	"crypto/sha256"
	"fmt"

	"go.uber.org/zap/zapcore"
)

type Hash struct {
	hash []byte
}

var GenesisHash = HashFromBytes([]byte("genesis hash"))

func HashFromBytes(data []byte) Hash {
	hash := sha256.Sum256(data)
	return Hash{hash: hash[:]}
}

func (h Hash) String() string {
	return fmt.Sprintf("%x", h.hash)
}

func (h Hash) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("hash", h.String())
	return nil
}

func (h Hash) Bytes() []byte {
	return h.hash
}

func (h Hash) Equal(other Hash) bool {
	return h.String() == other.String()
}
