package values

import (
	"go.uber.org/zap/zapcore"
)

const version = "0x00"

type Address struct {
	address string
}

func AddressFromString(address string) Address {
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
