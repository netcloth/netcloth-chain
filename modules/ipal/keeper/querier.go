package keeper

import (
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/ipal/types"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParameters:
			return queryParameters(ctx, k)
		case types.QueryServiceNodeList:
			return queryServiceNodeList(ctx, k)
		case types.QueryServiceNode:
			return queryServiceNode(ctx, req, k)
		case types.QueryServiceNodes:
			return queryServiceNodes(ctx, req, k)
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

func queryServiceNodeList(ctx sdk.Context, k Keeper) ([]byte, error) {
	serviceNodes := k.GetAllServiceNodes(ctx)
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, serviceNodes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryServiceNode(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var queryParams types.QueryServiceNodeParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &queryParams)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	serviceNode, found := k.GetServiceNode(ctx, queryParams.AccAddr)
	if found {
		bz, err := codec.MarshalJSONIndent(types.ModuleCdc, serviceNode)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return bz, nil
	}

	return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "not found")
}

func queryServiceNodes(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryServiceNodesParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	servcieNodes := types.ServiceNodes{}
	for _, accAddr := range params.AccAddrs {
		serviceNode, found := k.GetServiceNode(ctx, accAddr)
		if found {
			servcieNodes = append(servcieNodes, serviceNode)
		}
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, servcieNodes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}
