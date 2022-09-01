package values

import (
	"crypto/sha256"
	"fmt"

	"go.uber.org/zap/zapcore"
)

type Hash struct {
	hash []byte
}

func GenesisHash() Hash {
	hash := sha256.Sum256([]byte("Genesis"))

	return HashFromBytes(hash[:])
}

func HashFromBytes(data []byte) Hash {
	return Hash{hash: data}
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
