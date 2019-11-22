package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/cipal/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryCIPAL:
			return queryCIPAL(ctx, req, k)
		case types.QueryCIPALs:
			return queryCIPALs(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown cipal query endpoint")
		}
	}
}

func queryCIPAL(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var queryParams types.QueryCIPALParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &queryParams)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse accAddr: %s", err))
	}

	cipal, found := k.GetCIPALObject(ctx, queryParams.AccAddr)
	if found {
		bz, err := codec.MarshalJSONIndent(types.ModuleCdc, cipal)
		if err != nil {
			return []byte{}, sdk.ErrInternal(err.Error())
		}
		return bz, nil
	}

	return nil, sdk.ErrInternal("not found")
}

func queryCIPALs(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryCIPALsParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse accAddrs: %s", err))
	}

	cipals := types.CIPALObjects{}
	for _, addr := range params.AccAddrs {
		cipal, found := k.GetCIPALObject(ctx, addr)
		if found {
			cipals = append(cipals, cipal)
		}
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, cipals)
	if err != nil {
		return []byte{}, sdk.ErrInternal(err.Error())
	}
	return bz, nil
}
