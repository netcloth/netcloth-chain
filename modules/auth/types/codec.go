package types

import (
	"github.com/NetCloth/netcloth-chain/codec"
	"github.com/NetCloth/netcloth-chain/modules/auth/exported"
)

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "nch/Account", nil)
	cdc.RegisterConcrete(StdTx{}, "nch/StdTx", nil)
}

// module wide codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
