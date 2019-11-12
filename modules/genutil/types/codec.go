package types

import (
	"github.com/netcloth/netcloth-chain/codec"
	authtypes "github.com/netcloth/netcloth-chain/modules/auth/types"
	stakingtypes "github.com/netcloth/netcloth-chain/modules/staking/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// ModuleCdc defines a generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

// TODO: abstract genesis transactions registration back to staking
// required for genesis transactions
func init() {
	ModuleCdc = codec.New()
	stakingtypes.RegisterCodec(ModuleCdc)
	authtypes.RegisterCodec(ModuleCdc)
	sdk.RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
