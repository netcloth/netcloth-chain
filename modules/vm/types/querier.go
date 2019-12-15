package types

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	QueryParameters   = "params"
	QueryContractCode = "code"
	QueryStorage      = "storage"
)

type QueryCodeParams struct {
	AccAddr sdk.AccAddress
}

// QueryResCode is response type for code query
type QueryResCode struct {
	Value []byte `json:"value"`
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

func NewQueryCodeParams(AccAddr sdk.AccAddress) QueryCodeParams {
	return QueryCodeParams{
		AccAddr: AccAddr,
	}
}
