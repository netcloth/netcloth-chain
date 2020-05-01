package types

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/netcloth/netcloth-chain/app/v0/params"
	"github.com/netcloth/netcloth-chain/codec"
	nchtypes "github.com/netcloth/netcloth-chain/types"
	sdk "github.com/netcloth/netcloth-chain/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

// Staking params default values
const (
	// DefaultUnbondingTime reflects three weeks in seconds as the default
	// unbonding time.
	// TODO: Justify our choice of default here.
	DefaultUnbondingTime time.Duration = time.Hour * 24 * 14

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
	KeyUnbondingTime               = []byte("UnbondingTime")
	KeyMaxValidators               = []byte("MaxValidators")
	KeyMaxValidatorsExtending      = []byte("MaxValidatorsExtending")
	KeyMaxValidatorsExtendingSpeed = []byte("MaxValidatorsExtendingSpeed")
	KeyNextExtendingTime           = []byte("NextExtendingTime")
	KeyMaxEntries                  = []byte("KeyMaxEntries")
	KeyBondDenom                   = []byte("BondDenom")
	KeyMaxLever                    = []byte("MaxLever")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for staking
type Params struct {
	UnbondingTime               time.Duration `json:"unbonding_time" yaml:"unbonding_time"` // time duration of unbonding
	MaxValidators               uint16        `json:"max_validators" yaml:"max_validators"` // maximum number of validators (max uint16 = 65535)
	MaxValidatorsExtending      uint16        `json:"max_validators_extending" yaml:"max_validators_extending"`
	MaxValidatorsExtendingSpeed uint16        `json:"max_validators_extending_speed" yaml:"max_validators_extending_speed"`
	NextExtendingTime           int64         `json:"next_extending_time" yaml:"next_extending_time"`
	MaxEntries                  uint16        `json:"max_entries" yaml:"max_entries"` // max entries for either unbonding delegation or redelegation (per pair/trio)
	// note: we need to be a bit careful about potential overflow here, since this is user-determined
	BondDenom string  `json:"bond_denom" yaml:"bond_denom"` // bondable coin denomination
	MaxLever  sdk.Dec `json:"max_lever" yaml:"max_lever"`   // max lever: total user delegate / self delegate < max_lever
}

// NewParams creates a new Params instance
func NewParams(unbondingTime time.Duration, maxValidators, maxValidatorsExtending, maxValidatorsExtendingSpeed uint16, nextExtendingTime int64, maxEntries uint16,
	bondDenom string, maxLeverRate sdk.Dec) Params {

	return Params{
		UnbondingTime:               unbondingTime,
		MaxValidators:               maxValidators,
		MaxValidatorsExtending:      maxValidatorsExtending,
		MaxValidatorsExtendingSpeed: maxValidatorsExtendingSpeed,
		NextExtendingTime:           nextExtendingTime,
		MaxEntries:                  maxEntries,
		BondDenom:                   bondDenom,
		MaxLever:                    maxLeverRate,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyUnbondingTime, &p.UnbondingTime, validateUnbondingTime),
		params.NewParamSetPair(KeyMaxValidators, &p.MaxValidators, validateMaxValidators),
		params.NewParamSetPair(KeyMaxValidatorsExtending, &p.MaxValidatorsExtending, validateMaxValidatorsExtending),
		params.NewParamSetPair(KeyMaxValidatorsExtendingSpeed, &p.MaxValidatorsExtendingSpeed, validateMaxValidatorsExtendingSpeed),
		params.NewParamSetPair(KeyNextExtendingTime, &p.NextExtendingTime, validateNextExtendingTime),
		params.NewParamSetPair(KeyMaxEntries, &p.MaxEntries, validateMaxEntries),
		params.NewParamSetPair(KeyBondDenom, &p.BondDenom, validateBondDenom),
		params.NewParamSetPair(KeyMaxLever, &p.MaxLever, validateMaxLever),
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
		tmtime.Now().Unix()+MaxValidatorsExtendingInterval,
		DefaultMaxEntries,
		nchtypes.DefaultBondDenom,
		DefaultMaxLever)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
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

func validateUnbondingTime(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("validateUnbondingTime invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("unbonding time must be positive: %d", v)
	}

	return nil
}

func validateNextExtendingTime(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("validateNextExtendingTime invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("unbonding time must be positive: %d", v)
	}

	return nil
}

func validateMaxValidators(i interface{}) error {
	v, ok := i.(uint16)
	if !ok {
		return fmt.Errorf("validateMaxValidators invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max validators must be positive: %d", v)
	}

	return nil
}

func validateMaxValidatorsExtending(i interface{}) error {
	v, ok := i.(uint16)
	if !ok {
		return fmt.Errorf("validateMaxValidatorsExtending invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max validators must be positive: %d", v)
	}

	return nil
}

func validateMaxValidatorsExtendingSpeed(i interface{}) error {
	v, ok := i.(uint16)
	if !ok {
		return fmt.Errorf("validateMaxValidatorsExtendingSpeed invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max validators must be positive: %d", v)
	}

	return nil
}

func validateMaxEntries(i interface{}) error {
	v, ok := i.(uint16)
	if !ok {
		return fmt.Errorf("validateMaxEntries invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max entries must be positive: %d", v)
	}

	return nil
}

func validateMaxLever(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("validateMaxLever invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("staking max lever cannot be negative: %s", v)
	}

	return nil
}

func validateBondDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("validateBondDenom invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("bond denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

// validate a set of params
func (p Params) Validate() error {
	if err := validateUnbondingTime(p.UnbondingTime); err != nil {
		return err
	}
	if err := validateMaxValidators(p.MaxValidators); err != nil {
		return err
	}
	if err := validateMaxValidatorsExtending(p.MaxValidatorsExtending); err != nil {
		return err
	}
	if err := validateMaxValidatorsExtendingSpeed(p.MaxValidatorsExtendingSpeed); err != nil {
		return err
	}
	if err := validateNextExtendingTime(p.NextExtendingTime); err != nil {
		return err
	}
	if err := validateMaxEntries(p.MaxEntries); err != nil {
		return err
	}
	if err := validateBondDenom(p.BondDenom); err != nil {
		return err
	}
	if err := validateMaxLever(p.MaxLever); err != nil {
		return err
	}

	return nil
}
