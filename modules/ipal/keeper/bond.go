package keeper

import (
	"fmt"
	"time"

	"github.com/netcloth/netcloth-chain/modules/ipal/types"
	sdk "github.com/netcloth/netcloth-chain/types"
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
	moduleFunds := k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetCoins().AmountOf(sdk.NativeTokenName)
	unbondingAmount := sdk.NewInt(0)
	moduleFundsErr := false

	store := ctx.KVStore(k.storeKey)
	itr := k.UnBondingQueueIterator(ctx, ctx.BlockHeader().Time)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		tMatureUnBondings := types.UnBondings{}
		v := itr.Value()
		k.cdc.MustUnmarshalBinaryLengthPrefixed(v, &tMatureUnBondings)
		for _, v := range tMatureUnBondings {
			if unbondingAmount.LT(moduleFunds) {
				unbondingAmount = unbondingAmount.Add(v.Amount.Amount)
			} else {
				moduleFundsErr = true
				break
			}
		}

		if moduleFundsErr == true {
			ctx.Logger().Error(fmt.Sprintf("module[%s] funds[%v] insufficient", types.ModuleName, moduleFunds.String()))
			break
		}

		matureUnBondings = append(matureUnBondings, tMatureUnBondings...)
		store.Delete(itr.Key())
	}

	return matureUnBondings
}

func (k Keeper) DoUnbond(ctx sdk.Context, unBonding types.UnBonding) error {
	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, unBonding.AccountAddress, sdk.NewCoins(unBonding.Amount))
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("DoUnbond failed, err:", err.Error()))
	}
	return err
}

func (k Keeper) toUnbondingQueue(ctx sdk.Context, aa sdk.AccAddress, amt sdk.Coin) {
	endTime := ctx.BlockHeader().Time.Add(k.GetUnbondingTime(ctx))
	unBonding := types.NewUnBonding(aa, amt, endTime)
	k.InsertUnBondingQueue(ctx, unBonding, endTime)
}
