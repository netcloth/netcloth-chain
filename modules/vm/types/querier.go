package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	QueryParameters = "params"
	QueryCode       = "code"
	QueryState      = "state"
	QueryStorage    = "storage"
	QueryTxLogs     = "logs"
	EstimateGas     = "estimate_gas"
)

type QueryCodeParams struct {
	AccAddr sdk.AccAddress
}

// QueryResCode is response type for code query
type QueryResCode struct {
	Value []byte `json:"value"`
}

type QueryLogs struct {
	Logs []*Log `json:"logs"`
}

func (q QueryResCode) String() string {
	return string(q.Value)
}

// QueryResStorage is response type for storage query
type QueryResStorage struct {
	Value sdk.Hash `json:"value"`
}

func (q QueryResStorage) String() string {
	return q.Value.String()
}

func (q QueryLogs) String() string {
	return string(fmt.Sprintf("%+v", q.Logs))
}

func NewQueryCodeParams(AccAddr sdk.AccAddress) QueryCodeParams {
	return QueryCodeParams{
		AccAddr: AccAddr,
	}
}

type QueryContractStateParams struct {
	From sdk.AccAddress
	To   sdk.AccAddress
	Data []byte
}

// creates a new instance of QueryProposalParams
func NewQueryContractStateParams(from, to sdk.AccAddress, data []byte) QueryContractStateParams {
	return QueryContractStateParams{
		From: from,
		To:   to,
		Data: data,
	}
}

type FeeResult struct {
	Gas uint64
}

func (r FeeResult) String() string {
	return fmt.Sprintf("Gas = %d\n", r.Gas)
}

type QueryFeeParams struct {
	From sdk.AccAddress
	To   sdk.AccAddress
	Data Data
}

type Data []byte

func (aa Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(aa))
}

func (aa *Data) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	d, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	*aa = append(*aa, d...)
	return nil
}

func NewQueryFeeParams(from, to sdk.AccAddress, data []byte) QueryFeeParams {
	return QueryFeeParams{
		From: from,
		To:   to,
		Data: data,
	}
}
