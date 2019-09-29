package types

import (
	"github.com/NetCloth/netcloth-chain/codec"
)

// module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers all necessary param module types with a given codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(ParameterChangeProposal{}, "nch/ParameterChangeProposal", nil)
}
