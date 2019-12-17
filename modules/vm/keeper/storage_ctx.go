package keeper

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

type StorageContext struct {
	ctx      *sdk.Context
	storeKey *sdk.StoreKey
}

func NewStorageContext(ctx *sdk.Context, storeKey *sdk.StoreKey) StorageContext {
	return StorageContext{
		ctx:      ctx,
		storeKey: storeKey,
	}
}
