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

// GetMaxCallCreateDepth return MaxCallCreateDepth from store
func (k Keeper) GetMaxCallCreateDepth(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyMaxCallCreateDepth, &res)
	return
}

// SetMaxCallCreateDepth save MaxCallCreateDepth to store
func (k Keeper) SetMaxCallCreateDepth(ctx sdk.Context, maxCallCreateDepth uint64) {
	k.paramstore.Set(ctx, types.KeyMaxCallCreateDepth, maxCallCreateDepth)
}

func (k Keeper) GetVMOpGasParams(ctx sdk.Context) (params [256]uint64) {
	k.paramstore.Get(ctx, types.KeyVMOpGasParams, &params)
	return
}

func (k Keeper) SetVMOpGasParams(ctx sdk.Context, params [256]uint64) {
	k.paramstore.Set(ctx, types.KeyVMOpGasParams, params)
}

// GetVMContractCreationGasParams return VMContractCreationGasParams from store
func (k Keeper) GetVMContractCreationGasParams(ctx sdk.Context) (params types.VMContractCreationGasParams) {
	k.paramstore.Get(ctx, types.KeyVMContractCreationGasParams, &params)
	return
}

// SetVMContractCreationGasParams save VMContractCreationGasParams to store
func (k Keeper) SetVMContractCreationGasParams(ctx sdk.Context, params types.VMContractCreationGasParams) {
	k.paramstore.Set(ctx, types.KeyVMContractCreationGasParams, params)
}

func (k Keeper) GetParams(ctx sdk.Context) (res types.Params) {
	return types.NewParams(
		k.GetMaxCodeSize(ctx),
		k.GetMaxCallCreateDepth(ctx),
		k.GetVMOpGasParams(ctx),
		k.GetVMContractCreationGasParams(ctx),
	)
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
