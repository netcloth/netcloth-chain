package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/vm/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryParameters:
			return queryParameters(ctx, k)
		case types.QueryContractCode:
			return queryCode(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown vm query endpoint")
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

func queryCode(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	if len(req.Data) != 20 {
		return nil, sdk.ErrInvalidAddress("address invalid")
	}

	accAddr := sdk.AccAddress(req.Data)
	acc := k.ak.GetAccount(ctx, accAddr)
	if acc == nil {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", accAddr))
	}

	code, found := k.GetContractCode(ctx, acc.GetCodeHash())
	if !found {
		return nil, types.ErrNoCodeExist(fmt.Sprintf("codeHash %s does not have code", acc.GetCodeHash()))
	}

	return code, nil
}
