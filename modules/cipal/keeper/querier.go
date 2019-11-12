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
		case types.QueryIPAL:
			return queryIPAL(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown cipal query endpoint")
		}
	}
}

func queryIPAL(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var queryParams types.QueryCIPALParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &queryParams)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse accAddr: %s", err))
	}

	ipal, found := k.GetIPALObject(ctx, queryParams.AccAddr)
	if found {
		ctx.Logger().Error("found")
		bz, err := codec.MarshalJSONIndent(types.ModuleCdc, ipal)
		if err != nil {
			return []byte{}, sdk.ErrInternal(err.Error())
		}
		return bz, nil
	}

	return nil, sdk.ErrInternal("not found")
}
