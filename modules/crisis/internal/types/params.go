package types

import (
	"github.com/NetCloth/netcloth-chain/modules/params"
	sdk "github.com/NetCloth/netcloth-chain/types"
)

// Default parameter namespace
const (
	DefaultParamspace = ModuleName
)

var (
	// key for constant fee parameter
	ParamStoreKeyConstantFee = []byte("ConstantFee")
)

// type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		ParamStoreKeyConstantFee, sdk.Coin{},
	)
}