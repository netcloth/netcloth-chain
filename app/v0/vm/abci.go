package vm

import (
	"github.com/netcloth/netcloth-chain/app/v0/vm/keeper"
	sdk "github.com/netcloth/netcloth-chain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func EndBlocker(ctx sdk.Context, keeper keeper.Keeper) []abci.ValidatorUpdate {
	// Gas costs are handled within msg handler so costs should be ignored
	ebCtx := ctx.WithBlockGasMeter(sdk.NewInfiniteGasMeter())

	// Update account balances before committing other parts of state
	keeper.StateDB.UpdateAccounts()

	// Commit state objects to KV store
	_, err := keeper.StateDB.WithContext(ebCtx).Commit(true)
	if err != nil {
		panic(err)
	}

	// Clear accounts cache after account data has been committed
	keeper.StateDB.ClearStateObjects()

	return []abci.ValidatorUpdate{}
}
