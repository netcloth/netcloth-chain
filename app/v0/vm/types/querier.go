package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	QueryParameters = "params"
	QueryState      = "state"
	QueryCode       = "code"
	QueryStorage    = "storage"
	QueryTxLogs     = "logs"
	EstimateGas     = "estimate_gas"
	QueryCall       = "call"
)

// for query logs
type QueryLogsResult struct {
	Logs []*Log `json:"logs"`
}

func (q QueryLogsResult) String() string {
	return fmt.Sprintf("%+v", q.Logs)
}

// for query storage
type QueryStorageResult struct {
	Value sdk.Hash `json:"value"`
}

func (q QueryStorageResult) String() string {
	return q.Value.String()
}

// for Gas Estimate
type SimulationResult struct {
	GasVM    uint64
	GasNonVM uint64
	Res      string
}

func (r SimulationResult) String() string {
	return fmt.Sprintf("GasVM = %d\nGasNonVM = %d\nRes = %s", r.GasVM, r.GasNonVM, r.Res)
}

type VMQueryResult struct {
	GasVM    uint64
	GasNonVM uint64
	Result   []interface{}
}

func (r VMQueryResult) String() string {
	j, err := json.Marshal(r)
	if err != nil {
		return fmt.Sprintf("GasVM = %d\nGasNonVM = %d\nValues = %s", r.GasVM, r.GasNonVM, err.Error())
	}
	return string(j)
}

type QueryStateParams struct {
	ShowCode     bool `json:"show_code" yaml:"show_code"`
	ContractOnly bool `json:"contract_only" yaml:"contract_only"`
}
