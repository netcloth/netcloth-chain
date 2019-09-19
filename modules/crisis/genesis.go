package crisis

import (
	sdk "github.com/NetCloth/netcloth-chain/types"
	"github.com/NetCloth/netcloth-chain/modules/crisis/internal/keeper"
	"github.com/NetCloth/netcloth-chain/modules/crisis/internal/types"
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
