package types

import (
	"github.com/netcloth/netcloth-chain/app/v0/supply/exported"
	"github.com/netcloth/netcloth-chain/codec"
)

// RegisterCodec registers the account types and interface
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.ModuleAccountI)(nil), nil)
	cdc.RegisterInterface((*exported.SupplyI)(nil), nil)
	cdc.RegisterConcrete(&ModuleAccount{}, "nch/ModuleAccount", nil)
	cdc.RegisterConcrete(&Supply{}, "nch/Supply", nil)
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	ModuleCdc = cdc.Seal()
}
