package keeper

import (
	"fmt"

	"github.com/NetCloth/netcloth-chain/codec"
	"github.com/NetCloth/netcloth-chain/modules/aipal/types"
	sdk "github.com/NetCloth/netcloth-chain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryParameters:
			return queryParameters(ctx, k)
		case types.QueryServiceNodeList:
			return queryServiceNodeList(ctx, k)
		case types.QueryServiceNode:
			return queryServiceNode(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown ipal query endpoint")
		}
	}
}

func queryParameters(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return res, nil
}

func queryServiceNodeList(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	serviceNodes := k.GetAllServiceNodes(ctx)
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, serviceNodes)
	if err != nil {
		return []byte{}, sdk.ErrInternal(err.Error())
	}
	return bz, nil
}

func queryServiceNode(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var queryParams types.QueryServiceNodeParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &queryParams)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse accAddr: %s", err))
	}

	serviceNode, found := k.GetServiceNode(ctx, queryParams.AccAddr)
	if found {
		ctx.Logger().Error("found")
		bz, err := codec.MarshalJSONIndent(types.ModuleCdc, serviceNode)
		if err != nil {
			return []byte{}, sdk.ErrInternal(err.Error())
		}
		return bz, nil
	}

	return nil, sdk.ErrInternal("not found")
}
