package keeper

import (
    "time"
    "github.com/NetCloth/netcloth-chain/modules/ipala/types"
    "github.com/NetCloth/netcloth-chain/modules/params"
    sdk "github.com/NetCloth/netcloth-chain/types"
)

const (
    DefaultParamspace = types.ModuleName
)

func ParamKeyTable() params.KeyTable {
    return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

func (k Keeper) GetUnbondingTime(ctx sdk.Context) (res time.Duration) {
    k.paramstore.Get(ctx, types.KeyUnbondingTime, &res)
    return
}

func (k Keeper) SetUnbondingTime(ctx sdk.Context, unbondingTime time.Duration) {
    k.paramstore.Set(ctx, types.KeyUnbondingTime, unbondingTime)
}

func (k Keeper) GetMinBond(ctx sdk.Context) (res sdk.Coin) {
    k.paramstore.Get(ctx, types.KeyMinBond, &res)
    return
}

func (k Keeper) SetMinBond(ctx sdk.Context, minBond sdk.Coin) {
    k.paramstore.Set(ctx, types.KeyMinBond, minBond)
}

func (k Keeper) GetParams(ctx sdk.Context) (res types.Params) {
    return types.NewParams(
        k.GetUnbondingTime(ctx),
        k.GetMinBond(ctx))
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
    k.paramstore.SetParamSet(ctx, &params)
}
