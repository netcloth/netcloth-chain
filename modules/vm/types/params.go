package types

import (
	"fmt"

	"github.com/netcloth/netcloth-chain/modules/params"
)

const (
	DefaultMaxCodeSize = 1024 * 1024

	CallCreateDepth uint64 = 1024 // Maximum depth of call/create stack.

	CreateAccountGas uint64 = 200 //

)

var (
	KeyMaxCodeSize = []byte("MaxCodeSize")
)

type Params struct {
	MaxCodeSize uint64 `json:"max_code_size" yaml:"max_code_size"`
}

var _ params.ParamSet = (*Params)(nil)

func NewParams(maxCodeSize uint64) Params {
	return Params{
		MaxCodeSize: maxCodeSize,
	}
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyMaxCodeSize, &p.MaxCodeSize},
	}
}

func DefaultParams() Params {
	return NewParams(
		DefaultMaxCodeSize,
	)
}

func (p Params) String() string {
	return fmt.Sprintf(`Params:
  MaxCodeSize   : %v`,
		p.MaxCodeSize)
}
