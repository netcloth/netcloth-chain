package keeper

import (
	"github.com/netcloth/netcloth-chain/app/v0/params"
	"github.com/netcloth/netcloth-chain/app/v0/vm/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	DefaultParamspace = types.ModuleName
)

func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

func (k Keeper) GetMaxCodeSize(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyMaxCodeSize, &res)
	return
}

func (k Keeper) SetMaxCodeSize(ctx sdk.Context, maxCodeSize uint64) {
	k.paramstore.Set(ctx, types.KeyMaxCodeSize, maxCodeSize)
}

func (k Keeper) GetVMOpGasParams(ctx sdk.Context) (params [256]uint64) {
	k.paramstore.Get(ctx, types.KeyVMOpGasParams, &params)
	return
}

func (k Keeper) SetVMOpGasParams(ctx sdk.Context, params [256]uint64) {
	k.paramstore.Set(ctx, types.KeyVMOpGasParams, params)
}

func (k Keeper) GetVMCommonGasParams(ctx sdk.Context) (params types.VMCommonGasParams) {
	k.paramstore.Get(ctx, types.KeyVMCommonGasParams, &params)
	return
}

func (k Keeper) SetVMCommonGasParams(ctx sdk.Context, params types.VMCommonGasParams) {
	k.paramstore.Set(ctx, types.KeyVMCommonGasParams, params)
}

func (k Keeper) GetParams(ctx sdk.Context) (res types.Params) {
	return types.NewParams(
		k.GetMaxCodeSize(ctx),
		k.GetVMOpGasParams(ctx),
		k.GetVMCommonGasParams(ctx),
	)
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
