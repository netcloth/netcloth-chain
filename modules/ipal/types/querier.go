package types

import (
	sdk "github.com/NetCloth/netcloth-chain/types"
)

const (
	QueryServiceNodeList = "list"
	QueryServiceNode     = "node"
	QueryParameters      = "params"
)

type QueryServiceNodeParams struct {
	AccAddr sdk.AccAddress
}

func NewQueryServiceNodeParams(AccAddr sdk.AccAddress) QueryServiceNodeParams {
	return QueryServiceNodeParams{
		AccAddr: AccAddr,
	}
}
