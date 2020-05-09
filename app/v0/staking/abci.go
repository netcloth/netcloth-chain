package staking

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/app/v0/staking/keeper"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// Called every block, update validator set
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	return k.BlockValidatorUpdates(ctx)
}
