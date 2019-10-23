package keeper

import (
	"fmt"
	"github.com/NetCloth/netcloth-chain/modules/params"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/NetCloth/netcloth-chain/codec"
	"github.com/NetCloth/netcloth-chain/modules/ipal/types"
	sdk "github.com/NetCloth/netcloth-chain/types"
)

const ModuleAccount = "ipal"

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey     sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc          *codec.Codec // The wire codec for binary encoding/decoding.
	supplyKeeper types.SupplyKeeper

	paramstore   params.Subspace

	// codespace
	codespace sdk.CodespaceType
}

// NewKeeper creates new instances of the nch Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, supplyKeeper types.SupplyKeeper, paramstore params.Subspace, codespace sdk.CodespaceType) Keeper {
	// ensure ipal module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:     storeKey,
		cdc:          cdc,
		supplyKeeper: supplyKeeper,
		paramstore:   paramstore.WithKeyTable(ParamKeyTable()),
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

func (k Keeper) SetIPALObject(ctx sdk.Context, obj types.IPALObject) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalIPALObject(k.cdc, obj)
	store.Set(types.GetIPALObjectKey(obj.UserAddress), bz)
	ctx.Logger().Info(string(types.GetIPALObjectKey(obj.UserAddress)))
}

func (k Keeper) GetServerNodeObject(ctx sdk.Context, operator sdk.AccAddress) (obj types.ServerNodeObject, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetServerNodeObjectKey(operator))
	if value == nil {
		return obj, false
	}

	obj = types.MustUnmarshalServerNodeObject(k.cdc, value)
	return obj, true
}

func (k Keeper) setServerNodeObject(ctx sdk.Context, obj types.ServerNodeObject) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalServerNodeObject(k.cdc, obj)
	store.Set(types.GetServerNodeObjectKey(obj.OperatorAddress), bz)
}

func (k Keeper) delServerNodeObject(ctx sdk.Context, accAddress sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetServerNodeObjectKey(accAddress))
}

func (k Keeper) setServerNodeObjectByStakeShares(ctx sdk.Context, obj types.ServerNodeObject) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetServerNodeObjectByStakeSharesKey(obj), obj.OperatorAddress)
}

func (k Keeper) delServerNodeObjectByStakeShares(ctx sdk.Context, obj types.ServerNodeObject) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetServerNodeObjectByStakeSharesKey(obj))
}

func (k Keeper) updateServerNode(ctx sdk.Context, old types.ServerNodeObject, new types.MsgServiceNodeClaim) {
	u := types.NewServerNodeObject(new.OperatorAddress, new.Moniker, new.Website, new.ServerEndPoint, new.Details, new.StakeShares)
	k.setServerNodeObject(ctx, u)

	k.delServerNodeObjectByStakeShares(ctx, old)
	k.setServerNodeObjectByStakeShares(ctx, u)
}

func (k Keeper) UnStake(ctx sdk.Context, accountAddress sdk.AccAddress, amt sdk.Coin) {
	completionTime := ctx.BlockHeader().Time.Add(k.GetUnbondingTime(ctx))
	unstaking := k.SetUnstakingEntry(ctx, accountAddress, ctx.BlockHeight(), completionTime, amt)
	k.InsertUnStakingQueue(ctx, unstaking, completionTime)
}

/*
	founded {
		stakeShares >= minStakeShares {
			stakeShares > curStakeShares {
				addStake(stakeShares - curStakeShares)
			} else (stakeShares < curStakeShares) {
			    subStake(curStakeShares - stakeShares)
		    } else {
		    }
			updateNode
		} else {
		    unStake
		    deleteNode
	    }
	} else {
		stakeShares >= minStakeShares {
			createNode
			doStake
		} else {
			return err
		}
	}
*/
func (k Keeper) DoServerNodeClaim(ctx sdk.Context, msg types.MsgServiceNodeClaim) (err sdk.Error) {
	minStakeShares := k.GetMinStakingShares(ctx)
	old, found := k.GetServerNodeObject(ctx, msg.OperatorAddress)
	if found {
		if msg.StakeShares.IsGTE(minStakeShares) {
			if old.StakeShares.IsLT(msg.StakeShares) {
				k.supplyKeeper.SendCoinsFromAccountToModule(ctx, msg.OperatorAddress, ModuleAccount, sdk.NewCoins(msg.StakeShares.Sub(old.StakeShares)))
			} else if msg.StakeShares.IsLT(old.StakeShares) {
				k.UnStake(ctx, msg.OperatorAddress, old.StakeShares.Sub(msg.StakeShares))
			} else {
			}
			k.updateServerNode(ctx, old, msg)
		} else {
			k.UnStake(ctx, msg.OperatorAddress, old.StakeShares)

			k.delServerNodeObject(ctx, old.OperatorAddress)
			k.delServerNodeObjectByStakeShares(ctx, old)
		}
	} else {
		if msg.StakeShares.IsGTE(minStakeShares) {
			serverNode := types.NewServerNodeObject(msg.OperatorAddress, msg.Moniker, msg.Website, msg.ServerEndPoint, msg.Details, msg.StakeShares)
			k.setServerNodeObject(ctx, serverNode)
			k.setServerNodeObjectByStakeShares(ctx, serverNode)
			k.supplyKeeper.SendCoinsFromAccountToModule(ctx, msg.OperatorAddress, ModuleAccount, sdk.NewCoins(msg.StakeShares.Sub(msg.StakeShares)))
		} else {
			return types.ErrStakeSharesInsufficient(fmt.Sprintf("stakeShares insufficient, min: %v", minStakeShares.String()))
		}
	}

	return nil
}

func (k Keeper) GetAllServerNodes(ctx sdk.Context) (serverNodes types.ServerNodeObjects) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ServerNodeObjectKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator := types.MustUnmarshalServerNodeObject(k.cdc, iterator.Value())
		serverNodes = append(serverNodes, validator)
	}
	return serverNodes
}