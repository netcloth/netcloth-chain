package types

import (
	"github.com/netcloth/netcloth-chain/codec"
)

// RegisterCodec - register the sdk message type
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgContract{}, "nch/MsgContract", nil)
}

// ModuleCdc - generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
