package crisis

import (
	"github.com/netcloth/netcloth-chain/modules/crisis/internal/keeper"
	"github.com/netcloth/netcloth-chain/modules/crisis/internal/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// new crisis genesis
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) {
	keeper.SetConstantFee(ctx, data.ConstantFee)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) types.GenesisState {
	constantFee := keeper.GetConstantFee(ctx)
	return types.NewGenesisState(constantFee)
}
