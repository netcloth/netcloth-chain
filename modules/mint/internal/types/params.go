package types

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/netcloth/netcloth-chain/modules/params"
	nchtypes "github.com/netcloth/netcloth-chain/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// Parameter store keys
var (
	KeyMintDenom          = []byte("MintDenom")
	KeyInflationRate      = []byte("InflationRate")
	KeyNextInflateHeight  = []byte("NextInflateHeight")
	KeyBlockProvision     = []byte("BlockProvision")
	KeyBlocksPerYear      = []byte("BlocksPerYear")
	KeyTotalSupplyCeiling = []byte("TotalSupplyCeiling")
)

// mint parameters
type Params struct {
	MintDenom          string  `json:"mint_denom" yaml:"mint_denom"`         // type of coin to mint
	InflationRate      sdk.Dec `json:"inflation_rate" yaml:"inflation_rate"` // current annual inflation rate
	NextInflateHeight  int64   `json:"next_inflate_height" yaml:"next_inflate_height"`
	BlockProvision     sdk.Dec `json:"block_provision" yaml:"block_provision"`
	BlocksPerYear      int64   `json:"blocks_per_year" yaml:"blocks_per_year"` // expected blocks per year
	TotalSupplyCeiling sdk.Int `json:"total_supply_ceiling" yaml:"total_supply_ceiling"`
}

// ParamTable for minting module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(mintDenom string, inflationRate sdk.Dec, nextInflateHeight int64,
	blockProvision sdk.Dec, blocksPerYear int64, totalSupplyCeiling sdk.Int) Params {

	return Params{
		MintDenom:          mintDenom,
		InflationRate:      inflationRate,
		NextInflateHeight:  nextInflateHeight,
		BlockProvision:     blockProvision,
		BlocksPerYear:      blocksPerYear,
		TotalSupplyCeiling: totalSupplyCeiling,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		MintDenom:          nchtypes.DefaultBondDenom,
		InflationRate:      sdk.NewDecWithPrec(90, 2),
		NextInflateHeight:  int64(0),
		BlockProvision:     sdk.NewDecWithPrec(int64(11090830734910133), 2),
		BlocksPerYear:      int64(60 * 60 * 8766 / 5), // assuming 5 second block times
		TotalSupplyCeiling: sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(21), nil)),
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateInflationRate(p.InflationRate); err != nil {
		return err
	}
	if err := validateNextInflateHeight(p.NextInflateHeight); err != nil {
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
		params.NewParamSetPair(KeyInflationRate, &p.InflationRate, validateInflationRate),
		params.NewParamSetPair(KeyNextInflateHeight, &p.NextInflateHeight, validateNextInflateHeight),
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

func validateInflationRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("validateInflationRate invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("inflation rate cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("inflation rate too large: %s", v)
	}

	return nil
}

func validateNextInflateHeight(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("validateNextInflateHeight invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("next inflate height must be positive: %d", v)
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
