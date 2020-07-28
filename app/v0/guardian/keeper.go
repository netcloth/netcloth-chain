package guardian

import (
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// Keeper defines the guardian store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

// NewKeeper creates a new guardian Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		storeKey: key,
		cdc:      cdc,
	}
	return keeper
}

func (k Keeper) AddProfiler(ctx sdk.Context, guardian Guardian) error {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(guardian)
	store.Set(GetProfilerKey(guardian.Address), bz)
	return nil
}

func (k Keeper) DeleteProfiler(ctx sdk.Context, address sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetProfilerKey(address))
	return nil
}

func (k Keeper) GetProfiler(ctx sdk.Context, addr sdk.AccAddress) (guardian Guardian, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(GetProfilerKey(addr))
	if bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &guardian)
		return guardian, true
	}
	return guardian, false
}

func (k Keeper) ProfilersIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetProfilersSubspaceKey())
}
