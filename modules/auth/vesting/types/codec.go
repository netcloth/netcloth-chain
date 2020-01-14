package types

import (
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/auth/vesting/exported"
)

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.VestingAccount)(nil), nil)
	cdc.RegisterConcrete(&BaseVestingAccount{}, "nch/BaseVestingAccount", nil)
	cdc.RegisterConcrete(&ContinuousVestingAccount{}, "nch/ContinuousVestingAccount", nil)
	cdc.RegisterConcrete(&DelayedVestingAccount{}, "nch/DelayedVestingAccount", nil)
	cdc.RegisterConcrete(&PeriodicVestingAccount{}, "nch/PeriodicVestingAccount", nil)
}

// VestingCdc module wide codec
var VestingCdc *codec.Codec

func init() {
	VestingCdc = codec.New()
	RegisterCodec(VestingCdc)
	VestingCdc.Seal()
}
