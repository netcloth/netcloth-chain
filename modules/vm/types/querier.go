package types

import (
	"fmt"

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

// QueryResStorage is response type for storage query
type QueryResStorage struct {
	Value []byte `json:"value"`
}

func NewQueryCodeParams(AccAddr sdk.AccAddress) QueryCodeParams {
	return QueryCodeParams{
		AccAddr: AccAddr,
	}
}

func (p QueryResStorage) String() string {
	return fmt.Sprintf(`storage:
value   : %v`, p.Value)
}
