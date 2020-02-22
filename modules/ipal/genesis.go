package ipal

import (
	"fmt"

	"github.com/netcloth/netcloth-chain/modules/ipal/keeper"
	"github.com/netcloth/netcloth-chain/modules/ipal/types"
	sdk "github.com/netcloth/netcloth-chain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)

	for _, serviceNode := range data.ServiceNodes {
		fmt.Println("set serviceNode ")
		serviceNode.Bond = sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(0))
		fmt.Println(serviceNode)
		keeper.CreateServiceNode(ctx, serviceNode)
	}
	return []abci.ValidatorUpdate{}
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
