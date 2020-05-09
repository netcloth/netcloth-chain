package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"github.com/netcloth/netcloth-chain/app/v0/params"
	"github.com/netcloth/netcloth-chain/app/v0/vm/types"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
)

type Keeper struct {
	Cdc        *codec.Codec
	paramstore params.Subspace
	StateDB    *types.CommitStateDB
}

func NewKeeper(cdc *codec.Codec, storeKey, codeKey, storageDebugKey sdk.StoreKey, paramstore params.Subspace, ak auth.AccountKeeper) Keeper {
	return Keeper{
		Cdc:        cdc,
		paramstore: paramstore.WithKeyTable(ParamKeyTable()),
		StateDB:    types.NewCommitStateDB(ak, storeKey, codeKey, storageDebugKey),
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("modules/%s", types.ModuleName))
}

func (k Keeper) GetState(ctx sdk.Context, addr sdk.AccAddress, hash sdk.Hash) sdk.Hash {
	return k.StateDB.WithContext(ctx).GetState(addr, hash)
}

func (k *Keeper) GetCode(ctx sdk.Context, addr sdk.AccAddress) []byte {
	return k.StateDB.WithContext(ctx).GetCode(addr)
}

func (k *Keeper) GetLogs(ctx sdk.Context, hash sdk.Hash) []*types.Log {
	return k.StateDB.WithContext(ctx).GetLogs(hash)
}
