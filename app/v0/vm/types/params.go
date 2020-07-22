package types

import (
	"encoding/json"
	"fmt"

	"github.com/netcloth/netcloth-chain/app/v0/params"
)

const (
	defaultMaxCodeSize     = 1024 * 1024
	defaultCallCreateDepth = 1024

	defaultContractCreationGas = 53000
	defaultGasPerByte          = 200
)

// nolint
var (
	KeyMaxCodeSize                 = []byte("MaxCodeSize")
	KeyMaxCallCreateDepth          = []byte("MaxCallCreateDepth")
	KeyVMOpGasParams               = []byte("VMOpGasParams")
	KeyVMContractCreationGasParams = []byte("VMContractCreationGasParams")

	DefaultVMOpGasParams = [256]uint64{
		0, 3, 5, 3, 5, 5, 5, 5, 8, 8, 0, 5, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, //0-31
		30, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 700, 2, 2, 2, 3, 2, 3, 2, 3, 2, 700, 700, 2, 3, 700, //32-63
		20, 2, 2, 2, 0, 2, 2, 5, 0, 0, 0, 0, 0, 0, 0, 0, 2, 3, 3, 3, 800, 0, 8, 10, 2, 2, 2, 1, 0, 0, 0, 0, //64-95
		3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, //96-127
		3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, //128-159
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, //160-191
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, //192-223
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32000, 700, 700, 0, 700, 32000, 0, 0, 0, 0, 700, 0, 0, 0, 0, 0, //224-255
	}

	vmContractCreationGasParams = VMContractCreationGasParams{Gas: defaultContractCreationGas, GasPerByte: defaultGasPerByte}
)

// VMContractCreationGasParams contract creation gas params
type VMContractCreationGasParams struct {
	Gas        uint64 `json:"gas" yaml:"gas"`
	GasPerByte uint64 `json:"gas_per_byte" yaml:"gas_per_byte"`
}

type Params struct {
	MaxCodeSize                 uint64                      `json:"max_code_size" yaml:"max_code_size"`
	MaxCallCreateDepth          uint64                      `json:"max_call_create_depth" yaml:"max_call_create_depth"`
	VMOpGasParams               [256]uint64                 `json:"vm_op_gas_params" yaml:"vm_op_gas_params"`
	VMContractCreationGasParams VMContractCreationGasParams `json:"vm_contract_creation_gas_params" yaml:"vm_contract_creation_gas_params"`
}

var _ params.ParamSet = (*Params)(nil)

// NewParams return Params
func NewParams(maxCodeSize, callCreateDepth uint64, vmOpGasParams [256]uint64, vmContractCreationGasParams VMContractCreationGasParams) Params {
	return Params{
		MaxCodeSize:                 maxCodeSize,
		MaxCallCreateDepth:          callCreateDepth,
		VMOpGasParams:               vmOpGasParams,
		VMContractCreationGasParams: vmContractCreationGasParams,
	}
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyMaxCodeSize, &p.MaxCodeSize, validateMaxCodeSize),
		params.NewParamSetPair(KeyMaxCallCreateDepth, &p.MaxCallCreateDepth, validateMaxCallCreateDepth),
		params.NewParamSetPair(KeyVMOpGasParams, &p.VMOpGasParams, validateVMOpGasParams),
		params.NewParamSetPair(KeyVMContractCreationGasParams, &p.VMContractCreationGasParams, validateVMCommonGasParams),
	}
}

func DefaultParams() Params {
	return NewParams(
		defaultMaxCodeSize,
		defaultCallCreateDepth,
		DefaultVMOpGasParams,
		vmContractCreationGasParams,
	)
}

func (p Params) String() string {
	d, _ := json.Marshal(p)
	return string(d)
}

func validateMaxCodeSize(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("MaxCodeSize'type must be uint64: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max code size must be positive: %d", v)
	}

	return nil
}

func validateMaxCallCreateDepth(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("MaxCallCreateDepth'type must be uint64: %T", i)
	}

	return nil
}

func validateVMOpGasParams(i interface{}) error {
	return nil
}

func validateVMCommonGasParams(i interface{}) error {
	v, ok := i.(VMContractCreationGasParams)
	if !ok {
		return fmt.Errorf("invalid type: %T", i)
	}

	if v.Gas == 0 {
		return fmt.Errorf("gas must be positive: %d", v.Gas)
	}

	if v.GasPerByte == 0 {
		return fmt.Errorf("gas_per_byte must be positive: %d", v.GasPerByte)
	}

	return nil
}
