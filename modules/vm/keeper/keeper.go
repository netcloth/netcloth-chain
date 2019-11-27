package keeper

import (
	"fmt"

	"github.com/netcloth/netcloth-chain/modules/auth/exported"

	"github.com/netcloth/netcloth-chain/modules/params"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/vm/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// keeper of the staking store
type Keeper struct {
	storeKey   sdk.StoreKey
	storeTKey  sdk.StoreKey
	cdc        *codec.Codec
	paramstore params.Subspace
	ak         types.AccountKeeper
	bk         types.BankKeeper

	// codespace
	codespace sdk.CodespaceType
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(cdc *codec.Codec, key, tkey sdk.StoreKey,
	codespace sdk.CodespaceType, paramstore params.Subspace, ak types.AccountKeeper, bk types.BankKeeper) Keeper {

	return Keeper{
		storeKey:   key,
		storeTKey:  tkey,
		cdc:        cdc,
		paramstore: paramstore.WithKeyTable(ParamKeyTable()),
		codespace:  codespace,
		ak:         ak,
		bk:         bk,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("modules/%s", types.ModuleName))
}

func (k Keeper) GetAccount(ctx sdk.Context, address sdk.AccAddress) exported.Account {
	return k.ak.GetAccount(ctx, address)
}

func (k Keeper) Transfer(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amount sdk.Coins) {
	k.bk.SendCoins(ctx, fromAddr, toAddr, amount)
}
