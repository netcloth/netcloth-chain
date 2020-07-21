package params

import (
	"fmt"

	"github.com/netcloth/netcloth-chain/app/v0/params/subspace"
	"github.com/netcloth/netcloth-chain/app/v0/params/types"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"

	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the global paramstore
type Keeper struct {
	cdc    *codec.Codec
	key    sdk.StoreKey
	tkey   sdk.StoreKey
	spaces map[string]*Subspace
}

// NewKeeper constructs a params keeper
// nolint
func NewKeeper(cdc *codec.Codec, key *sdk.KVStoreKey, tkey *sdk.TransientStoreKey) (k Keeper) {
	k = Keeper{
		cdc:    cdc,
		key:    key,
		tkey:   tkey,
		spaces: make(map[string]*Subspace),
	}

	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("modules/%s", types.ModuleName))
}

// Allocate subspace used for keepers
func (k Keeper) Subspace(s string) Subspace {
	_, ok := k.spaces[s]
	if ok {
		panic("subspace already occupied")
	}

	if s == "" {
		panic("cannot use empty string for subspace")
	}

	space := subspace.NewSubspace(k.cdc, k.key, k.tkey, s)
	k.spaces[s] = &space

	return space
}

// Get existing substore from keeper
func (k Keeper) GetSubspace(s string) (Subspace, bool) {
	space, ok := k.spaces[s]
	if !ok {
		return Subspace{}, false
	}
	return *space, ok
}
