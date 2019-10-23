package keeper

import (
	"time"
	"github.com/NetCloth/netcloth-chain/modules/ipal/types"
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

func (k Keeper) GetMinStakingShares(ctx sdk.Context) (res sdk.Coin) {
	k.paramstore.Get(ctx, types.KeyMinStakeShares, &res)
	return
}

func (k Keeper) SetUnbondingTime(ctx sdk.Context, unbondingTime time.Duration) {
	k.paramstore.Set(ctx, types.KeyUnbondingTime, unbondingTime)
}

func (k Keeper) SetMinStakShares(ctx sdk.Context, minStakeShares sdk.Coin) {
	k.paramstore.Set(ctx, types.KeyMinStakeShares, minStakeShares)
}

func (k Keeper) GetParams(ctx sdk.Context) (res types.Params) {
	return types.NewParams(
		k.GetUnbondingTime(ctx),
		k.GetMinStakingShares(ctx))
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}