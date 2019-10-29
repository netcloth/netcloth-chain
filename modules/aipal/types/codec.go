package types

import (
    "github.com/NetCloth/netcloth-chain/codec"
)

func RegisterCodec(cdc *codec.Codec) {
    cdc.RegisterConcrete(MsgServiceNodeClaim{}, "nch/ServiceNodeClaim", nil)
}

var ModuleCdc *codec.Codec

func init() {
    ModuleCdc = codec.New()
    RegisterCodec(ModuleCdc)
    codec.RegisterCrypto(ModuleCdc)
    ModuleCdc.Seal()
}
