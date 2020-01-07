package types

import (
	"fmt"

	"github.com/netcloth/netcloth-chain/modules/params"
)

const (
	DefaultMaxCodeSize uint64 = 1024 * 1024
	CallCreateDepth    uint64 = 1024
)

var (
	KeyMaxCodeSize       = []byte("MaxCodeSize")
	KeyVMOpGasParams     = []byte("VMOpGasParams")
	KeyVMCommonGasParams = []byte("VMCommonGasParams")

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

	DefaultVMCommonGasParams = VMCommonGasParams{CreateDataGas: 200} //protocol_params.go::CreateDataGas
)

type VMCommonGasParams struct {
	CreateDataGas uint64 `json:"create_data_gas" yaml:"create_data_gas"`
}

type Params struct {
	MaxCodeSize       uint64            `json:"max_code_size" yaml:"max_code_size"`
	VMOpGasParams     [256]uint64       `json:"vm_op_gas_params" yaml:"vm_op_gas_params"`
	VMCommonGasParams VMCommonGasParams `json:"vm_common_gas_params" yaml:"vm_common_gas_params"`
}

var _ params.ParamSet = (*Params)(nil)

func NewParams(maxCodeSize uint64, vmOpGasParams [256]uint64, vmCommonGasParams VMCommonGasParams) Params {
	return Params{
		MaxCodeSize:       maxCodeSize,
		VMOpGasParams:     vmOpGasParams,
		VMCommonGasParams: vmCommonGasParams,
	}
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyMaxCodeSize, &p.MaxCodeSize},
		{KeyVMOpGasParams, p.VMOpGasParams},
		{KeyVMCommonGasParams, p.VMCommonGasParams},
	}
}

func DefaultParams() Params {
	return NewParams(
		DefaultMaxCodeSize,
		DefaultVMOpGasParams,
		DefaultVMCommonGasParams,
	)
}

func (p Params) String() string {
	return fmt.Sprintf(`Params:
  MaxCodeSize   : %v`,
		p.MaxCodeSize)
}
