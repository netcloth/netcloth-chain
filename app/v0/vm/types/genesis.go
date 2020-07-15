package types

import (
	"github.com/netcloth/netcloth-chain/hexutil"
	sdk "github.com/netcloth/netcloth-chain/types"
)

type (
	// GenesisState vm genesis state, include params, vm storage, vm codes, vm logs
	GenesisState struct {
		Params  Params              `json:"params"`
		Storage []Storage           `json:"storage"`
		Codes   map[string]sdk.Code `json:"codes"`
		VMLogs  VMLogs              `json:"vm_logs"`
	}

	// Storage vm storage of k, v pairs
	Storage struct {
		Key   hexutil.Bytes `json:"key"`
		Value hexutil.Bytes `json:"value"`
	}

	// VMLogs vm logs, include log index, logs which conclude k(txHash), v(logs) pairs
	VMLogs struct {
		Logs     map[string]string `json:"logs"`
		LogIndex int64             `json:"log_index"`
	}
)

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:  DefaultParams(),
		Storage: make([]Storage, 0),
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
