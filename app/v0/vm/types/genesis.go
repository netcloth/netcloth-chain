package types

import (
	"github.com/netcloth/netcloth-chain/hexutil"
	sdk "github.com/netcloth/netcloth-chain/types"
)

type (
	GenesisState struct {
		Params  Params              `json:"params"`
		Storage []GenesisStorage    `json:"storage"`
		Codes   map[string]sdk.Code `json:"codes"`
		VMLogs  VMLogs              `json:"vm_logs""`
	}

	GenesisStorage struct {
		Key   hexutil.Bytes `json:"key"`
		Value hexutil.Bytes `json:"value"`
	}

	VMLogs struct {
		Logs     map[string]string `json:"logs"`
		LogIndex int64             `json:"log_index"`
	}
)

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:  DefaultParams(),
		Storage: make([]GenesisStorage, 0),
	}
}

func NewGenesisState(params Params) GenesisState {
	return GenesisState{Params: params}
}

func ValidateGenesis(data GenesisState) error {
	maxCodeSize := data.Params.MaxCodeSize
	if err := validateMaxCodeSize(maxCodeSize); err != nil {
		return err
	}

	vmOpGasParams := data.Params.VMOpGasParams
	if err := validateVMOpGasParams(vmOpGasParams); err != nil {
		return err
	}

	vmCommonGasParams := data.Params.VMCommonGasParams
	return validateVMCommonGasParams(vmCommonGasParams)
}
