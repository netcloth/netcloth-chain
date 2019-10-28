package keeper

import (
    "time"
    "github.com/NetCloth/netcloth-chain/modules/aipal/types"
    sdk "github.com/NetCloth/netcloth-chain/types"
)

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

func (k Keeper) InsertUnBondingQueue(ctx sdk.Context, unBonding types.UnBonding, endTime time.Time) {
    s := k.GetUnBondingQueueTimeSlice(ctx, endTime)
    s = append(s, unBonding)
    k.SetUnBondingQueueTimeSlice(ctx, endTime, s)
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
