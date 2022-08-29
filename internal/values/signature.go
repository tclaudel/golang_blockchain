package values

import (
	"fmt"
	"math/big"
)

type Signature struct {
	r *big.Int
	s *big.Int
}

func (s Signature) String() string {
	return fmt.Sprintf("%s%s", s.r.String(), s.s.String())
}
