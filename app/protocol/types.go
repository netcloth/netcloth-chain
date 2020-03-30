package protocol

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

type Protocol interface {
	GetVersion() uint64
	GetRouter() sdk.Router
	GetInitChainer() sdk.InitChainer

	Load()
	Init(ctx sdk.Context)
}
