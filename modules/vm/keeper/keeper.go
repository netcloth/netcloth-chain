package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/staking/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// keeper of the staking store
type Keeper struct {
	storeKey  sdk.StoreKey
	storeTKey sdk.StoreKey
	cdc       *codec.Codec

	// codespace
	codespace sdk.CodespaceType
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(cdc *codec.Codec, key, tkey sdk.StoreKey,
	codespace sdk.CodespaceType) Keeper {

	return Keeper{
		storeKey:  key,
		storeTKey: tkey,
		cdc:       cdc,
		codespace: codespace,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("modules/%s", types.ModuleName))
}
