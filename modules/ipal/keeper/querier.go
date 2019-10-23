package keeper

import (
	"github.com/NetCloth/netcloth-chain/codec"
	"github.com/NetCloth/netcloth-chain/modules/ipal/types"
	sdk "github.com/NetCloth/netcloth-chain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryParameters:
			return queryParameters(ctx, k)
		case types.QueryServerNode:
			return queryServerNode(ctx, k)
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

func queryServerNode(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	serverNodes := k.GetAllServerNodes(ctx)
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, serverNodes)
	if err != nil {
		return []byte{}, sdk.ErrInternal(err.Error())
	}
	return bz, nil
}