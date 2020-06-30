package types

import (
	"github.com/netcloth/netcloth-chain/app/v0/auth/exported"
	"github.com/netcloth/netcloth-chain/codec"
)

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	cdc.RegisterInterface((*exported.VestingAccount)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "nch/Account", nil)
	cdc.RegisterConcrete(&BaseVestingAccount{}, "nch/BaseVestingAccount", nil)
	cdc.RegisterConcrete(&ContinuousVestingAccount{}, "nch/ContinuousVestingAccount", nil)
	cdc.RegisterConcrete(&DelayedVestingAccount{}, "nch/DelayedVestingAccount", nil)
	cdc.RegisterConcrete(StdTx{}, "nch/StdTx", nil)
}

// ModuleCdc - generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
