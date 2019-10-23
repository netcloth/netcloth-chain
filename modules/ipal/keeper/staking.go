package keeper

import (
	"github.com/NetCloth/netcloth-chain/modules/ipal/types"
	sdk "github.com/NetCloth/netcloth-chain/types"
	"time"
)

func (k Keeper) GetUnStaking(ctx sdk.Context, accountAddress sdk.AccAddress) (unstaking types.Unstaking, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUnStakingKey(accountAddress)
	value := store.Get(key)
	if value == nil {
		return unstaking, false
	}

	unstaking = types.MustUnmarshalUnstaking(k.cdc, value)
	return unstaking, true
}

func(k Keeper) SetUnStaking(ctx sdk.Context, unstaking types.Unstaking) {
	store := ctx.KVStore(k.storeKey)
	value := types.MustMarshalUnstaking(k.cdc, unstaking)
	key := types.GetUnStakingKey(unstaking.AccountAddress)
	store.Set(key, value)
}

func (k Keeper) RemoveUnStaking(ctx sdk.Context, unstaking types.Unstaking) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUnStakingKey(unstaking.AccountAddress)
	store.Delete(key)
}

func (k Keeper) SetUnstakingEntry(ctx sdk.Context, accountAddress sdk.AccAddress, creationHeight int64, minTime time.Time, amount sdk.Coin) types.Unstaking {
	unstaking, found := k.GetUnStaking(ctx, accountAddress)
	if found {
		unstaking.AddEntry(creationHeight, minTime, amount)
	} else {
		unstaking = types.NewUnstaking(accountAddress, creationHeight, minTime, amount)
	}
	k.SetUnStaking(ctx, unstaking)
	return unstaking
}

func (k Keeper) GetUnStakingQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (todos []types.UnStakingTODO) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetUnstakingTimeKey(timestamp))
	if value == nil {
		return []types.UnStakingTODO{}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &todos)
	return todos
}

func (k Keeper) SetUnStakingQueueTimeSlice(ctx sdk.Context, timestamp time.Time, todos []types.UnStakingTODO) {
	store := ctx.KVStore(k.storeKey)
	value := k.cdc.MustMarshalBinaryLengthPrefixed(todos)
	store.Set(types.GetUnstakingTimeKey(timestamp), value)
}

func (k Keeper) InsertUnStakingQueue(ctx sdk.Context, unstaking types.Unstaking, completionTime time.Time) {
	timeSlice :=  k.GetUnStakingQueueTimeSlice(ctx, completionTime)
	todos := []types.UnStakingTODO{}

	for _, v := range unstaking.Entries {
		todo := types.UnStakingTODO {
			AccountAddress: unstaking.AccountAddress,
			Amount:         v.Amount,
			EndTime:        v.EndTime,
		}
		todos = append(todos, todo)
	}

	if len(timeSlice) == 0 {
		k.SetUnStakingQueueTimeSlice(ctx, completionTime, todos)
	} else {
		for _, v := range todos {
			timeSlice = append(timeSlice, v)
		}
		k.SetUnStakingQueueTimeSlice(ctx, completionTime, timeSlice)
	}
}

func (k Keeper) UnStakingQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.UnStakingTODOKey, sdk.InclusiveEndBytes(types.GetUnstakingTimeKey(endTime)))
}

func (k Keeper) DequeueAllMatureUnStakingQueue(ctx sdk.Context, curTime time.Time) (matureUnstakings []types.UnStakingTODO) {
	store := ctx.KVStore(k.storeKey)

	unStakingTimeSliceIterator := k.UnStakingQueueIterator(ctx, ctx.BlockHeader().Time)
	for ; unStakingTimeSliceIterator.Valid(); unStakingTimeSliceIterator.Next() {
		timeslice := []types.UnStakingTODO{}
		value := unStakingTimeSliceIterator.Value()
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &timeslice)
		matureUnstakings = append(matureUnstakings, timeslice...)
		store.Delete(unStakingTimeSliceIterator.Key())
	}

	return matureUnstakings
}

func (k Keeper) DoUnStaking(ctx sdk.Context, todo types.UnStakingTODO) sdk.Error {
	return k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, ModuleAccount, todo.AccountAddress, sdk.NewCoins(todo.Amount))
}