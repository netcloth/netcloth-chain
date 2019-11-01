package keeper

import (
    abci "github.com/tendermint/tendermint/abci/types"

    "github.com/NetCloth/netcloth-chain/modules/ipal/types"
    sdk "github.com/NetCloth/netcloth-chain/types"
)

func NewQuerier(k Keeper) sdk.Querier {
    return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
        switch path[0] {
        case types.QueryIPAL:
            return nil, sdk.ErrUnknownRequest("unknown ipal query endpoint") //TODO fixme
        default:
            return nil, sdk.ErrUnknownRequest("unknown ipal query endpoint")
        }
    }
}
