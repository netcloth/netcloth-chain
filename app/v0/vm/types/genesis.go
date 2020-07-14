package types

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

type (
	GenesisState struct {
		Params    Params              `json:"params"`
		Contracts []Contract          `json:"contracts"`
		Codes     map[string]sdk.Code `json:"codes"`
		VMLogs    VMLogs              `json:"vm_logs""`
	}

	Contract struct {
		Address  sdk.AccAddress   `json:"address"`
		Coins    sdk.Coins        `json:"coins"`
		CodeHash sdk.Hash         `json:"code_hash"`
		Storage  []GenesisStorage `json:"storage"`
	}

	GenesisStorage struct {
		Key   sdk.Hash `json:"key"`
		Value sdk.Hash `json:"value"`
	}

	VMLogs struct {
		Logs     map[string]string `json:"logs"`
		LogIndex int64             `json:"log_index"`
	}
)

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:    DefaultParams(),
		Contracts: make([]Contract, 0, 10240),
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
