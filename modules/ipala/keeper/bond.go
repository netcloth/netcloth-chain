package keeper

import (
    "time"
    "github.com/NetCloth/netcloth-chain/modules/ipala/types"
    sdk "github.com/NetCloth/netcloth-chain/types"
)

func (k Keeper) GetUnBonds(ctx sdk.Context, aa sdk.AccAddress) (unBonds types.UnBonds, found bool) {
    store := ctx.KVStore(k.storeKey)
    key := types.GetUnBondsKey(aa)
    value := store.Get(key)
    if value == nil {
        return unBonds, false
    }

    unBonds = types.MustUnmarshalUnstaking(k.cdc, value)//TODO check register type to codec????????????????
    return unBonds, true
}

func(k Keeper) SetUnBonds(ctx sdk.Context, unBonds types.UnBonds) {
    store := ctx.KVStore(k.storeKey)
    value := types.MustMarshalUnstaking(k.cdc, unBonds)
    key := types.GetUnBondsKey(unBonds.AccountAddress)
    store.Set(key, value)
}

func (k Keeper) RemoveUnBonds(ctx sdk.Context, unBonds types.UnBonds) {
    store := ctx.KVStore(k.storeKey)
    key := types.GetUnBondsKey(unBonds.AccountAddress)
    store.Delete(key)
}

func (k Keeper) SetUnBondsEntry(ctx sdk.Context, aa sdk.AccAddress, endTime time.Time, amt sdk.Coin) types.UnBonds {
    unBonds, found := k.GetUnBonds(ctx, aa)
    if found {
        unBonds.AddEntry(endTime, amt)
    } else {
        unBonds = types.NewUnBonds(aa, endTime, amt)
    }
    k.SetUnBonds(ctx, unBonds)
    return unBonds
}

func (k Keeper) GetUnBondingQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (unBondings types.UnBondings) {
    store := ctx.KVStore(k.storeKey)
    value := store.Get(types.GetUnBondingKey(timestamp))
    if value == nil {
        return []types.UnBonding{}
    }
    k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &unBondings)
    return unBondings
}

func (k Keeper) SetUnBondingQueueTimeSlice(ctx sdk.Context, timestamp time.Time, todos types.UnBondings) {
    store := ctx.KVStore(k.storeKey)
    value := k.cdc.MustMarshalBinaryLengthPrefixed(todos)
    store.Set(types.GetUnBondingKey(timestamp), value)
}

func (k Keeper) InsertUnBondingQueue(ctx sdk.Context, unBonds types.UnBonds, endTime time.Time) {
    currentUnBondings :=  k.GetUnBondingQueueTimeSlice(ctx, endTime)
    unBondings := types.UnBondings{}

    for _, v := range unBonds.Entries {
        unBonding := types.NewUnBonding(unBonds.AccountAddress, v.Amount, v.EndTime)
        unBondings = append(unBondings, unBonding)
    }

    if len(currentUnBondings) == 0 {
        k.SetUnBondingQueueTimeSlice(ctx, endTime, unBondings)
    } else {
        currentUnBondings = append(currentUnBondings, unBondings...)
        k.SetUnBondingQueueTimeSlice(ctx, endTime, currentUnBondings)
    }
}

func (k Keeper) UnBondingQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
    store := ctx.KVStore(k.storeKey)
    return store.Iterator(types.UnBondingKey, sdk.InclusiveEndBytes(types.GetUnBondingKey(endTime)))
}

func (k Keeper) DequeueAllMatureUnBondingQueue(ctx sdk.Context, curTime time.Time) (matureUnBondings []types.UnBonding) {
    store := ctx.KVStore(k.storeKey)

    itr := k.UnBondingQueueIterator(ctx, ctx.BlockHeader().Time)
    for ; itr.Valid(); itr.Next() {
        tMatureUnBondings := types.UnBondings{}
        v := itr.Value()
        k.cdc.MustUnmarshalBinaryLengthPrefixed(v, &tMatureUnBondings)
        matureUnBondings = append(matureUnBondings, tMatureUnBondings...)
        //TODO check tMatureUnBondings.clear()????????
        store.Delete(itr.Key())
    }

    return matureUnBondings
}

func (k Keeper) DoUnBond(ctx sdk.Context, unBonding types.UnBonding) sdk.Error {
    //TODO check
    return k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, unBonding.AccountAddress, sdk.NewCoins(unBonding.Amount))
}
