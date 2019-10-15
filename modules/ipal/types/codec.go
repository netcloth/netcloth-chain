package types

import (
	"github.com/NetCloth/netcloth-chain/codec"
)

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgIPALClaim{}, "nch/IPALCLaim", nil)
	cdc.RegisterConcrete(MsgServiceNodeClaim{}, "nch/ServerNodeClaim", nil)
}

// module wide codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
