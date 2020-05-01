package types

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/netcloth/netcloth-chain/app/v0/params"
	nchtypes "github.com/netcloth/netcloth-chain/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// Parameter store keys
var (
	KeyMintDenom                  = []byte("MintDenom")
	KeyInflationCutBackRate       = []byte("InflationCutBackRate")
	KeyNextInflationCutBackHeight = []byte("NextInflationCutBackHeight")
	KeyBlockProvision             = []byte("BlockProvision")
	KeyBlocksPerYear              = []byte("BlocksPerYear")
	KeyTotalSupplyCeiling         = []byte("TotalSupplyCeiling")
)

// mint parameters
type Params struct {
	MintDenom                  string  `json:"mint_denom" yaml:"mint_denom"`                         // type of coin to mint
	InflationCutBackRate       sdk.Dec `json:"inflation_cutback_rate" yaml:"inflation_cutback_rate"` // current annual inflate cutback  rate
	NextInflationCutBackHeight int64   `json:"next_inflation_cutback_height" yaml:"next_inflation_cutback_height"`
	BlockProvision             sdk.Dec `json:"block_provision" yaml:"block_provision"`
	BlocksPerYear              int64   `json:"blocks_per_year" yaml:"blocks_per_year"` // expected blocks per year
	TotalSupplyCeiling         sdk.Int `json:"total_supply_ceiling" yaml:"total_supply_ceiling"`
}

// ParamTable for minting module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(mintDenom string, inflationCutBackRate sdk.Dec, nextInflationCutBackHeight int64,
	blockProvision sdk.Dec, blocksPerYear int64, totalSupplyCeiling sdk.Int) Params {

	return Params{
		MintDenom:                  mintDenom,
		InflationCutBackRate:       inflationCutBackRate,
		NextInflationCutBackHeight: nextInflationCutBackHeight,
		BlockProvision:             blockProvision,
		BlocksPerYear:              blocksPerYear,
		TotalSupplyCeiling:         totalSupplyCeiling,
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		MintDenom:                  nchtypes.DefaultBondDenom,
		InflationCutBackRate:       sdk.NewDecWithPrec(90, 2),
		NextInflationCutBackHeight: int64(0),
		BlockProvision:             sdk.NewDec(11090830734911),
		BlocksPerYear:              int64(60 * 60 * 8766 / 5), // assuming 5 second block times
		TotalSupplyCeiling:         sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(21), nil)),
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateInflationCutBackRate(p.InflationCutBackRate); err != nil {
		return err
	}
	if err := validateNextInflationCutBackHeight(p.NextInflationCutBackHeight); err != nil {
		return err
	}
	if err := validateBlockProvision(p.BlockProvision); err != nil {
		return err
	}
	if err := validateBlocksPerYear(p.BlocksPerYear); err != nil {
		return err
	}
	if err := validateTotalSupplyCeiling(p.TotalSupplyCeiling); err != nil {
		return err
	}

	return nil
}

func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		params.NewParamSetPair(KeyInflationCutBackRate, &p.InflationCutBackRate, validateInflationCutBackRate),
		params.NewParamSetPair(KeyNextInflationCutBackHeight, &p.NextInflationCutBackHeight, validateNextInflationCutBackHeight),
		params.NewParamSetPair(KeyBlockProvision, &p.BlockProvision, validateBlockProvision),
		params.NewParamSetPair(KeyBlocksPerYear, &p.BlocksPerYear, validateBlocksPerYear),
		params.NewParamSetPair(KeyTotalSupplyCeiling, &p.TotalSupplyCeiling, validateTotalSupplyCeiling),
	}
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("validateMintDenom invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateInflationCutBackRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("validateInflationCutBackRate invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("inflation cutback rate cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("inflation cutback rate too large: %s", v)
	}

	return nil
}

func validateNextInflationCutBackHeight(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("validateNextInflationCutBackHeight invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("next inflation cutback height must be positive: %d", v)
	}

	return nil
}

func validateBlockProvision(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("validateBlockProvision invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("block provision cannot be negative: %s", v)
	}

	return nil
}

func validateBlocksPerYear(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("validateBlocksPerYear invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("blocks per year must be positive: %d", v)
	}

	return nil
}

func validateTotalSupplyCeiling(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("validateTotalSupplyCeiling invalid parameter type: %T", i)
	}

	if v.LTE(sdk.NewInt(int64(0))) {
		return fmt.Errorf("total supply ceiling must be positive: %d", v)
	}

	return nil
}
