package types

import (
	"github.com/netcloth/netcloth-chain/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "nch/MsgSend", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "nch/MsgMultiSend", nil)
}

// ModuleCdc codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
