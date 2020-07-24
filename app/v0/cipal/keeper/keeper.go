package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/netcloth/netcloth-chain/app/v0/cipal/types"
	"github.com/netcloth/netcloth-chain/app/v0/params"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
)

type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramstore params.Subspace
}

func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, paramstore params.Subspace) Keeper {
	return Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
		paramstore: paramstore.WithKeyTable(ParamKeyTable()),
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("modules/%s", types.ModuleName))
}

func (k Keeper) GetCIPALObject(ctx sdk.Context, userAddress string) (obj types.CIPALObject, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetCIPALObjectKey(userAddress))
	ctx.Logger().Info(string(types.GetCIPALObjectKey(userAddress)))
	if value == nil {
		return obj, false
	}

	obj = types.MustUnmarshalCIPALObject(k.cdc, value)
	return obj, true
}

func (k Keeper) GetAllCIPALObjects(ctx sdk.Context) (objs []types.CIPALObject) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.CIPALObjectKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		obj := types.MustUnmarshalCIPALObject(k.cdc, iterator.Value())
		objs = append(objs, obj)
	}
	return objs
}

// get the set of all cipal object with no limits
func (k Keeper) GetCIPALObjectCount(ctx sdk.Context) (count int) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.CIPALObjectKey)
	defer iterator.Close()

	count = 0
	for ; iterator.Valid(); iterator.Next() {
		count++
	}
	return count
}

func (k Keeper) SetCIPALObject(ctx sdk.Context, obj types.CIPALObject) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalCIPALObject(k.cdc, obj)
	store.Set(types.GetCIPALObjectKey(obj.UserAddress), bz)
	//ctx.Logger().Info(string(types.GetCIPALObjectKey(obj.UserAddress)))
}
