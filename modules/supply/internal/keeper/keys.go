package keeper

import (
	"github.com/netcloth/netcloth-chain/modules/supply/internal/types"
)

// DefaultCodespace from the supply module
var DefaultCodespace string = types.ModuleName

// Keys for supply store
// Items are stored with the following key: values
//
// - 0x00: Supply
var (
	SupplyKey = []byte{0x00}
)
