package mint

import (
	"fmt"

	"github.com/netcloth/netcloth-chain/modules/mint/internal/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k Keeper) {
	// fetch stored params
	params := k.GetParams(ctx)
	blockHeight := ctx.BlockHeight()

	if blockHeight <= 1 {
		params.NextInflateHeight = params.NextInflateHeight + params.BlocksPerYear
		k.SetParams(ctx, params)
	} else if blockHeight == params.NextInflateHeight {
		//TODO, check inflate ceiling
		params.BlockProvision = params.InflationRate.Mul(params.BlockProvision)
		params.NextInflateHeight = params.NextInflateHeight + params.BlocksPerYear
		k.SetParams(ctx, params)
	}

	// mint coins, update supply
	mintedCoin := sdk.NewCoin(params.MintDenom, params.BlockProvision.TruncateInt())
	mintedCoins := sdk.NewCoins(mintedCoin)
	fmt.Println(fmt.Sprintf("minted coins: %s", mintedCoins.String()))

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
			sdk.NewAttribute(types.AttributeKeyInflation, params.InflationRate.String()),
			sdk.NewAttribute(types.AttributeKeyBlockProvision, params.BlockProvision.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)
}
