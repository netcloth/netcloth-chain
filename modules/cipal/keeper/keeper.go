package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/NetCloth/netcloth-chain/codec"
	"github.com/NetCloth/netcloth-chain/modules/cipal/types"
	"github.com/NetCloth/netcloth-chain/modules/params"
	sdk "github.com/NetCloth/netcloth-chain/types"
)

type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramstore params.Subspace
	codespace  sdk.CodespaceType
}

func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, paramstore params.Subspace, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
		paramstore: paramstore.WithKeyTable(ParamKeyTable()),
		codespace:  codespace,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

func (k Keeper) GetIPALObject(ctx sdk.Context, userAddress string) (obj types.CIPALObject, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetCIPALObjectKey(userAddress))
	ctx.Logger().Info(string(types.GetCIPALObjectKey(userAddress)))
	if value == nil {
		return obj, false
	}

	obj = types.MustUnmarshalCIPALObject(k.cdc, value)
	return obj, true
}

func (k Keeper) SetIPALObject(ctx sdk.Context, obj types.CIPALObject) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalCIPALObject(k.cdc, obj)
	store.Set(types.GetCIPALObjectKey(obj.UserAddress), bz)
	ctx.Logger().Info(string(types.GetCIPALObjectKey(obj.UserAddress)))
}
