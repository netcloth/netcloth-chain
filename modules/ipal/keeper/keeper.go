package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/NetCloth/netcloth-chain/codec"
	"github.com/NetCloth/netcloth-chain/modules/ipal/types"
	sdk "github.com/NetCloth/netcloth-chain/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey     sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc          *codec.Codec // The wire codec for binary encoding/decoding.
	supplyKeeper types.SupplyKeeper

	// codespace
	codespace sdk.CodespaceType
}

// NewKeeper creates new instances of the nch Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, supplyKeeper types.SupplyKeeper, codespace sdk.CodespaceType) Keeper {
	// ensure ipal module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:     storeKey,
		cdc:          cdc,
		supplyKeeper: supplyKeeper,
		codespace:    codespace,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// get a single ipal object
func (k Keeper) GetIPALObject(ctx sdk.Context, userAddress, serverIP string) (obj types.IPALObject, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetIPALObjectKey(userAddress))
	ctx.Logger().Info(string(types.GetIPALObjectKey(userAddress)))
	if value == nil {
		return obj, false
	}

	obj = types.MustUnmarshalIPALObject(k.cdc, value)
	return obj, true
}

// set ipal object
func (k Keeper) SetIPALObject(ctx sdk.Context, obj types.IPALObject) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalIPALObject(k.cdc, obj)
	store.Set(types.GetIPALObjectKey(obj.UserAddress), bz)
	ctx.Logger().Info(string(types.GetIPALObjectKey(obj.UserAddress)))
}

// get a single ServerNode object
func (k Keeper) GetServerNodeObject(ctx sdk.Context, operator sdk.AccAddress) (obj types.ServerNodeObject, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetServerNodeObjectKey(operator))
	ctx.Logger().Info(string(types.GetServerNodeObjectKey(operator)))
	if value == nil {
		return obj, false
	}

	obj = types.MustUnmarshalServerNodeObject(k.cdc, value)
	return obj, true
}

// set ServerNode object
func (k Keeper) SetServerNodeObject(ctx sdk.Context, obj types.ServerNodeObject) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalServerNodeObject(k.cdc, obj)
	store.Set(types.GetServerNodeObjectKey(obj.OperatorAddress), bz)
	ctx.Logger().Info(string(types.GetServerNodeObjectKey(obj.OperatorAddress)))
}
