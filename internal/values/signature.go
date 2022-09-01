package values

import (
	"fmt"
	"math/big"
)

type Signature struct {
	r *big.Int
	s *big.Int
}

func SignatureFromStrings(r, s string) (Signature, error) {
	var (
		rInt = new(big.Int)
		sInt = new(big.Int)
	)

	if _, ok := rInt.SetString(r, 10); !ok {
		return Signature{}, fmt.Errorf("invalid r")
	}
	if _, ok := sInt.SetString(s, 10); !ok {
		return Signature{}, fmt.Errorf("invalid s")
	}

	return Signature{
		r: rInt,
		s: sInt,
	}, nil
}

func (sig Signature) String() string {
	return fmt.Sprintf("%s%s", sig.r.String(), sig.s.String())
}

func (sig Signature) Strings() (r string, s string) {
	return sig.r.String(), sig.s.String()
}
