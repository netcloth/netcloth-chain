package mint

import (
	"fmt"

	"github.com/netcloth/netcloth-chain/app/v0/mint/internal/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k Keeper) {
	// fetch stored params
	params := k.GetParams(ctx)
	blockHeight := ctx.BlockHeight()
	supply := k.StakingTokenSupply(ctx)

	// check total inflation ceiling
	// if total token supply >= TotalSupplyCeiling, stop inflating
	if supply.GTE(params.TotalSupplyCeiling) {
		ctx.Logger().Info(fmt.Sprintf("current token supply: %s, stop inflating", supply.String()))

		params.BlockProvision = sdk.NewDec(0)
		params.NextInflationCutBackHeight = 0
		k.SetParams(ctx, params)
		return
	}

	if blockHeight <= 1 {
		// update next inflate height at chain startup
		params.NextInflationCutBackHeight = params.NextInflationCutBackHeight + params.BlocksPerYear
		k.SetParams(ctx, params)
	} else if blockHeight == params.NextInflationCutBackHeight {
		// adjust block provision and next inflate height
		params.BlockProvision = params.InflationCutBackRate.Mul(params.BlockProvision)
		params.NextInflationCutBackHeight = params.NextInflationCutBackHeight + params.BlocksPerYear
		k.SetParams(ctx, params)
	}

	// mint coins, update token supply
	mintedCoin := sdk.NewCoin(params.MintDenom, params.BlockProvision.TruncateInt())
	mintedCoins := sdk.NewCoins(mintedCoin)
	err := k.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	// send the minted coins to the fee collector account
	err = k.AddCollectedFees(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyInflation, params.InflationCutBackRate.String()),
			sdk.NewAttribute(types.AttributeKeyBlockProvision, params.BlockProvision.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)
}
