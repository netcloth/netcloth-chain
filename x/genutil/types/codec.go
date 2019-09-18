package types

import (
	"github.com/NetCloth/netcloth-chain/codec"
	authtypes "github.com/NetCloth/netcloth-chain/x/auth/types"
	stakingtypes "github.com/NetCloth/netcloth-chain/x/staking/types"
	sdk "github.com/NetCloth/netcloth-chain/types"
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
