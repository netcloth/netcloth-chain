package types

import (
	"fmt"
	"github.com/NetCloth/netcloth-chain/codec"
	"github.com/NetCloth/netcloth-chain/modules/params"
	sdk "github.com/NetCloth/netcloth-chain/types"
	"time"
)

const (
	//DefaultUnbondingTime = time.Hour * 24 * 7
	DefaultUnbondingTime = time.Minute * 5
)

var (
	DefaultMinStakeShares = sdk.NewCoin("unch", sdk.NewInt(1000000 * 1))
)

var (
	KeyUnbondingTime    = []byte("UnbondingTime")
	KeyMinStakeShares   = []byte("MinStakeShares")
)

type Params struct {
	UnbondingTime		time.Duration   `json:"unbonding_time" yaml:"unbonding_time"`
	MinStakeShares      sdk.Coin        `json:"min_stake" yaml:"min_bond"`
}

var _ params.ParamSet = (*Params)(nil)

func NewParams(unbondingTime time.Duration, minStakeShares sdk.Coin) Params {
	return Params {
		UnbondingTime: unbondingTime,
		MinStakeShares: minStakeShares,
	}
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs {
		{KeyUnbondingTime, &p.UnbondingTime},
		{KeyMinStakeShares, &p.MinStakeShares},
	}
}

func DefaultParams() Params {
	return NewParams(
		DefaultUnbondingTime,
		DefaultMinStakeShares,
		)
}

func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Unbonding Time    : %s
  Min StakeShares   : %v`,
  p.UnbondingTime,
  p.MinStakeShares)
}

func MustUnmarshalParams(cdc *codec.Codec, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}
	return params
}

func UnmarshalParams(cdc *codec.Codec, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &params)
	if err != nil {
		return
	}
	return
}

func (p Params) Validate() error {
	//TODO
	return nil
}

