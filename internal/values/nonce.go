package values

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

// Nonce is a value that represents a nonce.
type Nonce struct {
	// nonce is the nonce.
	nonce int
}

func InitNonce() Nonce {
	return Nonce{nonce: 0}
}

func (n Nonce) Increase() Nonce {
	return Nonce{nonce: n.nonce + 1}
}

func NonceFromInt(nonce int) Nonce {
	return Nonce{nonce: nonce}
}

var GenesisNonce = NonceFromInt(0)

func (n Nonce) String() string {
	return fmt.Sprintf("%d", n.nonce)
}

func (n Nonce) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("nonce", n.String())
	return nil
}

func (n Nonce) Int() int {
	return n.nonce
}

func (n Nonce) Equal(other Nonce) bool {
	return n.nonce == other.nonce
}
