package types

import (
	"fmt"
	"time"

	"github.com/netcloth/netcloth-chain/modules/params"
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	DefaultUnbondingTime = time.Hour * 24 * 7
)

var (
	DefaultMinBond = sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(sdk.NativeTokenFraction))
)

var (
	KeyUnbondingTime = []byte("UnbondingTime")
	KeyMinBond       = []byte("MinBond")
)

type Params struct {
	UnbondingTime time.Duration `json:"unbonding_time" yaml:"unbonding_time"`
	MinBond       sdk.Coin      `json:"min_bond" yaml:"min_bond"`
}

var _ params.ParamSet = (*Params)(nil)

func NewParams(unbondingTime time.Duration, minBond sdk.Coin) Params {
	return Params{
		UnbondingTime: unbondingTime,
		MinBond:       minBond,
	}
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyUnbondingTime, &p.UnbondingTime, validateUnbondingTime),
		params.NewParamSetPair(KeyMinBond, &p.MinBond, validateMinBond),
	}
}

func DefaultParams() Params {
	return NewParams(
		DefaultUnbondingTime,
		DefaultMinBond,
	)
}

func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Unbonding Time    : %s
  Min Bond   : %v`,
		p.UnbondingTime,
		p.MinBond)
}

func validateUnbondingTime(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("unbonding time must be positive: %d", v)
	}

	return nil
}

func validateMinBond(i interface{}) error {
	// TODO
	return nil
}
