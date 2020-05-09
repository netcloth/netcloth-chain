package guardian

import (
	"errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryProfilers:
			return queryProfilers(ctx, k)
		default:
			return nil, errors.New("unknown guardian query endpoint")
		}
	}
}

func queryProfilers(ctx sdk.Context, k Keeper) ([]byte, error) {
	profilersIterator := k.ProfilersIterator(ctx)
	defer profilersIterator.Close()

	var profilers []Guardian
	for ; profilersIterator.Valid(); profilersIterator.Next() {
		var profiler Guardian
		k.cdc.MustUnmarshalBinaryLengthPrefixed(profilersIterator.Value(), &profiler)
		profilers = append(profilers, profiler)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, profilers)
	if err != nil {
		return nil, err
	}
	return bz, nil
}
