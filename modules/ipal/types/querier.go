package types

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	QueryServiceNodeList = "list"
	QueryServiceNode     = "node"
	QueryServiceNodes    = "nodes"
	QueryParameters      = "params"
)

type QueryServiceNodeParams struct {
	AccAddr sdk.AccAddress
}

type QueryServiceNodesParams struct {
	AccAddrs []sdk.AccAddress `json:"acc_addrs"`
}

func NewQueryServiceNodeParams(AccAddr sdk.AccAddress) QueryServiceNodeParams {
	return QueryServiceNodeParams{
		AccAddr: AccAddr,
	}
}

func NewQueryServiceNodesParams(AccAddrs []sdk.AccAddress) QueryServiceNodesParams {
	return QueryServiceNodesParams{
		AccAddrs: AccAddrs,
	}
}
