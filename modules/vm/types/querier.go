package types

import (
	"fmt"

	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	QueryParameters    = "params"
	QueryContractCode  = "code"
	QueryContractState = "state"
	QueryStorage       = "storage"
	QueryTxLogs        = "logs"
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
	Value []byte `json:"value"`
}

func (q QueryResStorage) String() string {
	return string(q.Value)
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
