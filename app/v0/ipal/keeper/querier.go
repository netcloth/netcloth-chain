package keeper

import (
	"github.com/netcloth/netcloth-chain/app/v0/ipal/types"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParameters:
			return queryParameters(ctx, k)
		case types.QueryIPALNodeList:
			return queryIPALNodeList(ctx, k)
		case types.QueryIPALNode:
			return queryIPALNode(ctx, req, k)
		case types.QueryIPALNodes:
			return queryIPALNodes(ctx, req, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown ipal query path: %s", path[0])
		}
	}
}

func queryParameters(ctx sdk.Context, k Keeper) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryIPALNodeList(ctx sdk.Context, k Keeper) ([]byte, error) {
	ipalNodes := k.GetAllIPALNodes(ctx)
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, ipalNodes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryIPALNode(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var queryParams types.QueryIPALNodeParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &queryParams)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	ipalNode, found := k.GetIPALNode(ctx, queryParams.AccAddr)
	if found {
		bz, err := codec.MarshalJSONIndent(types.ModuleCdc, ipalNode)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return bz, nil
	}

	return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "not found")
}

func queryIPALNodes(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryIPALNodesParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	ipalNodes := types.IPALNodes{}
	for _, accAddr := range params.AccAddrs {
		node, found := k.GetIPALNode(ctx, accAddr)
		if found {
			ipalNodes = append(ipalNodes, node)
		}
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, ipalNodes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}
