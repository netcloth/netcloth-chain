package types

import sdk "github.com/netcloth/netcloth-chain/types"

const (
	QueryParameters   = "params"
	QueryContractCode = "code"
)

type QueryCodeParams struct {
	AccAddr sdk.AccAddress
}

func NewQueryCodeParams(AccAddr sdk.AccAddress) QueryCodeParams {
	return QueryCodeParams{
		AccAddr: AccAddr,
	}
}
