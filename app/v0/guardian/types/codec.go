package types

import (
	"github.com/netcloth/netcloth-chain/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgAddProfiler{}, "nch/guardian/MsgAddProfiler", nil)
	cdc.RegisterConcrete(MsgDeleteProfiler{}, "nch/guardian/MsgDeleteProfiler", nil)
	cdc.RegisterConcrete(Guardian{}, "nch/guardian/Guardian", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
