package protocol

import (
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
)

type Protocol interface {
	GetVersion() uint64
	GetRouter() sdk.Router
	GetInitChainer() sdk.InitChainer
	GetBeginBlocker() sdk.BeginBlocker
	GetEndBlocker() sdk.EndBlocker

	Load()
	Init(ctx sdk.Context)
	GetCodec() *codec.Codec
}
