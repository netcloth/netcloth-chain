package types

import (
	"fmt"

	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	QueryParameters = "params"
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
type EstimateGasResult struct {
	Gas uint64
}

func (r EstimateGasResult) String() string {
	return fmt.Sprintf("Gas = %d\n", r.Gas)
}
