package types

import (
	"encoding/json"

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
		Key   hexutil.Bytes `json:"k"`
		Value hexutil.Bytes `json:"v"`
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
	if err := validateMaxCodeSize(data.Params.MaxCodeSize); err != nil {
		return err
	}

	if err := validateMaxCallCreateDepth(data.Params.MaxCallCreateDepth); err != nil {
		return err
	}

	vmOpGasParams := data.Params.VMOpGasParams
	if err := validateVMOpGasParams(vmOpGasParams); err != nil {
		return err
	}

	return validateVMCommonGasParams(data.Params.VMContractCreationGasParams)
}

// Equal judge GenesisState equal
func (a GenesisState) Equal(b GenesisState) bool {
	aJSON, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bJSON, err := json.Marshal(b)
	if err != nil {
		return false
	}
	aJSONSorted, err := sdk.SortJSON(aJSON)
	if err != nil {
		return false
	}
	bJSONSorted, err := sdk.SortJSON(bJSON)
	if err != nil {
		return false
	}
	return string(aJSONSorted) == string(bJSONSorted)
}

// EqualWithoutParams judge GenesisState equal without compare GenesisState.Params
func (a GenesisState) EqualWithoutParams(b GenesisState) bool {
	a.Params = Params{}
	b.Params = Params{}
	return a.Equal(b)
}
