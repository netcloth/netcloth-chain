package cipal

import (
	"github.com/netcloth/netcloth-chain/modules/cipal/keeper"
	"github.com/netcloth/netcloth-chain/modules/cipal/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) {
	for _, obj := range data.CIPALObjs {
		keeper.SetCIPALObject(ctx, obj)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	cipals := keeper.GetAllCIPALObjects(ctx)
	return types.GenesisState{
		CIPALObjs: cipals,
	}
}
