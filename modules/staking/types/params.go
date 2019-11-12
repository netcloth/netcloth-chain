package types

import (
	"bytes"
	"fmt"
	"time"

	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/params"
	nchtypes "github.com/netcloth/netcloth-chain/types"
	sdk "github.com/netcloth/netcloth-chain/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

// Staking params default values
const (
	// DefaultUnbondingTime reflects three weeks in seconds as the default
	// unbonding time.
	// TODO: Justify our choice of default here.
	DefaultUnbondingTime time.Duration = time.Hour * 24 * 7

	// Default maximum number of bonded validators
	DefaultMaxValidators uint16 = 100

	DefaultMaxValidatorsExtending uint16 = 300

	DefaultMaxValidatorsExtendingSpeed uint16 = 10

	MaxValidatorsExtendingInterval = 60 * 60 * 8766

	// Default maximum entries in a UBD/RED pair
	DefaultMaxEntries uint16 = 7
)

var (
	// Default maximum lever
	DefaultMaxLever sdk.Dec = sdk.NewDec(20)
)

// nolint - Keys for parameter access
var (
	KeyUnbondingTime                = []byte("UnbondingTime")
	KeyMaxValidators                = []byte("MaxValidators")
	KeyMaxValidatorsExtending       = []byte("MaxValidatorsExtending")
	KeyMaxValidatorsExtendingSpeed  = []byte("MaxValidatorsExtendingSpeed")
	KeyNextExtendingTime            = []byte("NextExtendingTime")
	KeyMaxEntries                   = []byte("KeyMaxEntries")
	KeyBondDenom                    = []byte("BondDenom")
	KeyMaxLever                     = []byte("MaxLever")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for staking
type Params struct {
	UnbondingTime               time.Duration `json:"unbonding_time" yaml:"unbonding_time"` // time duration of unbonding
	MaxValidators               uint16        `json:"max_validators" yaml:"max_validators"` // maximum number of validators (max uint16 = 65535)
	MaxValidatorsExtending      uint16        `json:"max_validators_extending" yaml:"max_validators_extending"`
	MaxValidatorsExtendingSpeed uint16        `json:"max_validators_extending_speed" yaml:"max_validators_extending_speed"`
	NextExtendingTime           int64         `json:"next_extending_time" yaml:"next_extending_time"`
	MaxEntries                  uint16        `json:"max_entries" yaml:"max_entries"`       // max entries for either unbonding delegation or redelegation (per pair/trio)
	// note: we need to be a bit careful about potential overflow here, since this is user-determined
	BondDenom                   string        `json:"bond_denom" yaml:"bond_denom"` // bondable coin denomination
	MaxLever                    sdk.Dec       `json:"max_lever" yaml:"max_lever"`   // max lever: total user delegate / self delegate < max_lever
}

// NewParams creates a new Params instance
func NewParams(unbondingTime time.Duration, maxValidators, maxValidatorsExtending, maxValidatorsExtendingSpeed uint16, nextExtendingTime int64, maxEntries uint16,
	bondDenom string, maxLeverRate sdk.Dec) Params {

	return Params{
		UnbondingTime                   : unbondingTime,
		MaxValidators                   : maxValidators,
		MaxValidatorsExtending          : maxValidatorsExtending,
		MaxValidatorsExtendingSpeed     : maxValidatorsExtendingSpeed,
		NextExtendingTime               : nextExtendingTime,
		MaxEntries                      : maxEntries,
		BondDenom                       : bondDenom,
		MaxLever                        : maxLeverRate,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyUnbondingTime, &p.UnbondingTime},
		{KeyMaxValidators, &p.MaxValidators},
		{KeyMaxValidatorsExtending, &p.MaxValidatorsExtending},
		{KeyMaxValidatorsExtendingSpeed, &p.MaxValidatorsExtendingSpeed},
		{KeyNextExtendingTime, &p.NextExtendingTime},
		{KeyMaxEntries, &p.MaxEntries},
		{KeyBondDenom, &p.BondDenom},
		{KeyMaxLever, &p.MaxLever},
	}
}

// Equal returns a boolean determining if two Param types are identical.
// TODO: This is slower than comparing struct fields directly
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultUnbondingTime,
		DefaultMaxValidators,
		DefaultMaxValidatorsExtending,
		DefaultMaxValidatorsExtendingSpeed,
		tmtime.Now().Unix() + MaxValidatorsExtendingInterval,
		DefaultMaxEntries,
		nchtypes.DefaultBondDenom,
		DefaultMaxLever)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Unbonding Time                 : %s
  Max Validators                 : %d
  Max Validators Extending       : %d
  Max Validators Extending Speed : %d
  Next Extending Time            : %d
  Max Entries                    : %d
  Bonded Coin Denom              : %s
  Max Lever                      : %s`,
  p.UnbondingTime,
  p.MaxValidators,
  p.MaxValidatorsExtending,
  p.MaxValidatorsExtendingSpeed,
  p.NextExtendingTime,
  p.MaxEntries,
  p.BondDenom,
  p.MaxLever)
}

// unmarshal the current staking params value from store key or panic
func MustUnmarshalParams(cdc *codec.Codec, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}
	return params
}

// unmarshal the current staking params value from store key
func UnmarshalParams(cdc *codec.Codec, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &params)
	if err != nil {
		return
	}
	return
}

// validate a set of params
func (p Params) Validate() error {
	if p.BondDenom == "" {
		return fmt.Errorf("staking parameter BondDenom can't be an empty string")
	}
	if p.MaxValidators == 0 {
		return fmt.Errorf("staking parameter MaxValidators must be a positive integer")
	}
	return nil
}
