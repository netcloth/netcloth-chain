package guardian

import (
	"github.com/netcloth/netcloth-chain/app/v0/guardian/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

type GenesisState struct {
	Profilers []types.Guardian `json:"profilers"`
}

func NewGenesisState(profilers []types.Guardian) GenesisState {
	return GenesisState{
		Profilers: profilers,
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, profiler := range data.Profilers {
		keeper.AddProfiler(ctx, profiler)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	profilersIterator := k.ProfilersIterator(ctx)
	defer profilersIterator.Close()
	var profilers []types.Guardian
	for ; profilersIterator.Valid(); profilersIterator.Next() {
		var profiler types.Guardian
		k.cdc.MustUnmarshalBinaryLengthPrefixed(profilersIterator.Value(), &profiler)
		profilers = append(profilers, profiler)
	}

	return NewGenesisState(profilers)
}

func DefaultGenesisState() GenesisState {
	guardian := types.Guardian{Description: "genesis", AccountType: types.Genesis}
	return NewGenesisState([]types.Guardian{guardian})
}
