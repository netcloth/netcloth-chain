package types

import (
    "github.com/NetCloth/netcloth-chain/codec"
)

func RegisterCodec(cdc *codec.Codec) {
    cdc.RegisterConcrete(MsgIPALClaim{}, "nch/IPALClaim", nil)
}

var ModuleCdc *codec.Codec

func init() {
    ModuleCdc = codec.New()
    RegisterCodec(ModuleCdc)
    codec.RegisterCrypto(ModuleCdc)
    ModuleCdc.Seal()
}
