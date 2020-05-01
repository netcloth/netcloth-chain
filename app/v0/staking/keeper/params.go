package keeper

import (
	"time"

	"github.com/netcloth/netcloth-chain/app/v0/params"
	"github.com/netcloth/netcloth-chain/app/v0/staking/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// Default parameter namespace
const (
	DefaultParamspace = types.ModuleName
)

// ParamTable for staking module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// UnbondingTime
func (k Keeper) UnbondingTime(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.KeyUnbondingTime, &res)
	return
}

// MaxValidators - Maximum number of validators
func (k Keeper) MaxValidators(ctx sdk.Context) (res uint16) {
	k.paramstore.Get(ctx, types.KeyMaxValidators, &res)
	return
}

func (k Keeper) MaxValidatorsExtending(ctx sdk.Context) (res uint16) {
	k.paramstore.Get(ctx, types.KeyMaxValidatorsExtending, &res)
	return
}

func (k Keeper) MaxValidatorsExtendingSpeed(ctx sdk.Context) (res uint16) {
	k.paramstore.Get(ctx, types.KeyMaxValidatorsExtendingSpeed, &res)
	return
}

func (k Keeper) NextExtendingTime(ctx sdk.Context) (res int64) {
	k.paramstore.Get(ctx, types.KeyNextExtendingTime, &res)
	return
}

// MaxEntries - Maximum number of simultaneous unbonding
// delegations or redelegations (per pair/trio)
func (k Keeper) MaxEntries(ctx sdk.Context) (res uint16) {
	k.paramstore.Get(ctx, types.KeyMaxEntries, &res)
	return
}

// BondDenom - Bondable coin denomination
func (k Keeper) BondDenom(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyBondDenom, &res)
	return
}

// MaxLever - max delegation lever
// total user delegation / self delegation
func (k Keeper) MaxLever(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyMaxLever, &res)
	return
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.UnbondingTime(ctx),
		k.MaxValidators(ctx),
		k.MaxValidatorsExtending(ctx),
		k.MaxValidatorsExtendingSpeed(ctx),
		k.NextExtendingTime(ctx),
		k.MaxEntries(ctx),
		k.BondDenom(ctx),
		k.MaxLever(ctx),
	)
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
