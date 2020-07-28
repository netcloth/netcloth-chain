package keeper

import (
	"fmt"

	"github.com/netcloth/netcloth-chain/app/v0/ipal/types"
	"github.com/netcloth/netcloth-chain/app/v0/params"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

// Keeper defines the ipal store
type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	supplyKeeper types.SupplyKeeper
	paramstore   params.Subspace
}

// NewKeeper creates a new ipal Keeper instance
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, supplyKeeper types.SupplyKeeper, paramstore params.Subspace) Keeper {
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:     storeKey,
		cdc:          cdc,
		supplyKeeper: supplyKeeper,
		paramstore:   paramstore.WithKeyTable(ParamKeyTable()),
	}
}

// GetIPALNode returns a IPAL object by operator address
func (k Keeper) GetIPALNode(ctx sdk.Context, operator sdk.AccAddress) (obj types.IPALNode, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetIPALNodeKey(operator))
	if value == nil {
		return obj, false
	}

	obj = types.MustUnmarshalIPALNode(k.cdc, value)
	return obj, true
}

func (k Keeper) setIPALNode(ctx sdk.Context, obj types.IPALNode) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalIPALNode(k.cdc, obj)
	store.Set(types.GetIPALNodeKey(obj.OperatorAddress), bz)
}

func (k Keeper) delIPALNode(ctx sdk.Context, accAddress sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetIPALNodeKey(accAddress))
}

func (k Keeper) setIPALNodeByBond(ctx sdk.Context, obj types.IPALNode) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalIPALNode(k.cdc, obj)
	store.Set(types.GetIPALNodeByBondKey(obj), bz)
}

func (k Keeper) delIPALNodeByBond(ctx sdk.Context, obj types.IPALNode) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetIPALNodeByBondKey(obj))
}

func (k Keeper) setIPALNodeByMonikerIndex(ctx sdk.Context, obj types.IPALNode) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetIPALNodeByMonikerKey(obj.Moniker), obj.OperatorAddress)
}

func (k Keeper) GetIPALNodeAddByMoniker(ctx sdk.Context, moniker string) (acc sdk.AccAddress, exist bool) {
	store := ctx.KVStore(k.storeKey)
	v := store.Get(types.GetIPALNodeByMonikerKey(moniker))
	return v, v != nil
}

func (k Keeper) delIPALNodeByMonikerIndex(ctx sdk.Context, moniker string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetIPALNodeByMonikerKey(moniker))
}

// CreateIPALNode sets a new IPAL object
func (k Keeper) CreateIPALNode(ctx sdk.Context, node types.IPALNode) {
	k.setIPALNode(ctx, node)
	k.setIPALNodeByBond(ctx, node)
	k.setIPALNodeByMonikerIndex(ctx, node)
}

func (k Keeper) updateIPALNode(ctx sdk.Context, old types.IPALNode, new types.IPALNode) {
	k.setIPALNode(ctx, new)

	k.delIPALNodeByBond(ctx, old)
	k.setIPALNodeByBond(ctx, new)

	k.delIPALNodeByMonikerIndex(ctx, old.Moniker)
	k.setIPALNodeByMonikerIndex(ctx, new)
}

func (k Keeper) deleteIPALNode(ctx sdk.Context, obj types.IPALNode) {
	k.delIPALNode(ctx, obj.OperatorAddress)
	k.delIPALNodeByBond(ctx, obj)
	k.delIPALNodeByMonikerIndex(ctx, obj.Moniker)
}

func (k Keeper) bond(ctx sdk.Context, aa sdk.AccAddress, amt sdk.Coin) error {
	return k.supplyKeeper.SendCoinsFromAccountToModule(ctx, aa, types.ModuleName, sdk.Coins{amt})
}

// DoIPALNodeClaim - updates ipal object and bond coins
func (k Keeper) DoIPALNodeClaim(ctx sdk.Context, m types.MsgIPALNodeClaim) (err error) {
	minBond := k.GetMinBond(ctx)
	n, found := k.GetIPALNode(ctx, m.OperatorAddress)
	if found {
		if m.Bond.IsGTE(minBond) {
			if n.Bond.IsLT(m.Bond) {
				err := k.bond(ctx, m.OperatorAddress, m.Bond.Sub(n.Bond))
				if err != nil {
					return err
				}
			} else if m.Bond.IsLT(n.Bond) {
				k.toUnbondingQueue(ctx, m.OperatorAddress, n.Bond.Sub(m.Bond))
			}

			ipalNode := types.NewIPALNode(m.OperatorAddress, m.Moniker, m.Website, m.Details, m.Extension, m.Endpoints, m.Bond)
			k.updateIPALNode(ctx, n, ipalNode)
		} else {
			k.toUnbondingQueue(ctx, m.OperatorAddress, n.Bond)
			k.deleteIPALNode(ctx, n)
		}
	} else {
		if m.Bond.IsGTE(minBond) {
			err := k.bond(ctx, m.OperatorAddress, m.Bond)
			if err != nil {
				return err
			}

			ipalNode := types.NewIPALNode(m.OperatorAddress, m.Moniker, m.Website, m.Details, m.Extension, m.Endpoints, m.Bond)
			k.CreateIPALNode(ctx, ipalNode)
		} else {
			return sdkerrors.Wrapf(types.ErrBondInsufficient, "bond insufficient, min bond: %s, actual bond: %s", minBond.String(), m.Bond.String())
		}
	}

	return nil
}

// GetAllIPALNodes - lists all ipal objects
func (k Keeper) GetAllIPALNodes(ctx sdk.Context) (ipalNodes types.IPALNodes) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.IPALNodeByBondKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		n := types.MustUnmarshalIPALNode(k.cdc, iterator.Value())
		ipalNodes = append(ipalNodes, n)
	}
	return ipalNodes
}
