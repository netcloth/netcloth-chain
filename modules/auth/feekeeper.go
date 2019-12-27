package auth

import (
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
)

var (
	collectedFeesKey = []byte("collectedFees")
	feeAuthKey       = []byte("feeAuth")
)

type FeeKeeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

func NewFeeKeeper(cdc *codec.Codec, key sdk.StoreKey) FeeKeeper {
	return FeeKeeper{
		storeKey: key,
		cdc:      cdc,
	}
}

func (fk FeeKeeper) GetCollectedFees(ctx sdk.Context) sdk.Coins {
	store := ctx.KVStore(fk.storeKey)
	bz := store.Get(collectedFeesKey)
	if bz == nil {
		return sdk.Coins{}
	}

	feePool := &(sdk.Coins{})
	fk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, feePool)
	return *feePool
}

func (fk FeeKeeper) setCollectedFees(ctx sdk.Context, coins sdk.Coins) {
	bz := fk.cdc.MustMarshalBinaryLengthPrefixed(coins)
	store := ctx.KVStore(fk.storeKey)
	store.Set(collectedFeesKey, bz)
}

func (fk FeeKeeper) AddCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins {
	newCoins := fk.GetCollectedFees(ctx).Add(coins)
	fk.setCollectedFees(ctx, newCoins)

	return newCoins
}

func (fk FeeKeeper) RefundCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins {
	newCoins := fk.GetCollectedFees(ctx).Sub(coins)
	if newCoins.IsAnyNegative() {
		panic("fee collector contains negative coins")
	}
	fk.setCollectedFees(ctx, newCoins)
	return newCoins
}

func (fk FeeKeeper) ClearCollectedFees(ctx sdk.Context) {
	fk.setCollectedFees(ctx, sdk.Coins{})
}

func (fk FeeKeeper) GetFeeAuth(ctx sdk.Context) (feeAuth FeeAuth) {
	store := ctx.KVStore(fk.storeKey)
	b := store.Get(feeAuthKey)
	if b == nil {
		panic("stored fee pool should not be nil")
	}
	fk.cdc.MustUnmarshalBinaryLengthPrefixed(b, &feeAuth)
	return
}

func (fk FeeKeeper) SetFeeAuth(ctx sdk.Context, feeAuth FeeAuth) {
	store := ctx.KVStore(fk.storeKey)
	b := fk.cdc.MustMarshalBinaryLengthPrefixed(feeAuth)
	store.Set(feeAuthKey, b)
}
