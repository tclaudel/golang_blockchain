package values

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

type Amount struct {
	value float64
}

func AmountFromFloat64(value float64) Amount {
	return Amount{value: value}
}

func (v Amount) Float64() float64 {
	return v.value
}

func (v Amount) String() string {
	return fmt.Sprintf("%f", v.value)
}

func (v Amount) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("value", v.String())
	return nil
}
