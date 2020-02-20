package ipal

import (
	"github.com/netcloth/netcloth-chain/modules/ipal/keeper"
	"github.com/netcloth/netcloth-chain/modules/ipal/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) {
	//TODO
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)
	serviceNodes := keeper.GetAllServiceNodes(ctx)

	return types.GenesisState{
		Params:       params,
		ServiceNodes: serviceNodes,
	}
}
