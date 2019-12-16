package keeper

import (
	"fmt"
	"os"

	"github.com/netcloth/netcloth-chain/modules/auth/exported"

	tmcrypto "github.com/tendermint/tendermint/crypto"

	"github.com/netcloth/netcloth-chain/modules/params"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/vm/common"
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

	CSDB *types.CommitStateDB

	// codespace
	codespace sdk.CodespaceType
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(cdc *codec.Codec, key, tkey sdk.StoreKey, codespace sdk.CodespaceType, paramstore params.Subspace, ak types.AccountKeeper, bk types.BankKeeper, csdb *types.CommitStateDB) Keeper {
	return Keeper{
		storeKey:   key,
		storeTKey:  tkey,
		cdc:        cdc,
		paramstore: paramstore.WithKeyTable(ParamKeyTable()),
		codespace:  codespace,
		ak:         ak,
		bk:         bk,
		CSDB:       csdb,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("modules/%s", types.ModuleName))
}

func (k Keeper) GetContractCode(ctx sdk.Context, codeHash []byte) (code []byte, found bool) {
	store := ctx.KVStore(k.storeKey)
	code = store.Get(types.GetContractCodeKey(codeHash))
	return code, code != nil
}

func (k Keeper) setContractCode(ctx sdk.Context, codeHash, code []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetContractCodeKey(codeHash), code)
}

func (k Keeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	acc := k.ak.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.NewCoins()
	}
	return acc.GetCoins()
}

func (k Keeper) GetAccount(ctx sdk.Context, addr sdk.AccAddress) exported.Account {
	return k.ak.GetAccount(ctx, addr)
}

func (k Keeper) GetState(ctx sdk.Context, addr sdk.AccAddress, hash sdk.Hash) sdk.Hash {
	return k.CSDB.WithContext(ctx).GetState(addr, hash)
}

// GetCode calls CommitStateDB.GetCode using the passed in context
func (k *Keeper) GetCode(ctx sdk.Context, addr sdk.AccAddress) []byte {
	return k.CSDB.WithContext(ctx).GetCode(addr)
}

func (k *Keeper) GetLogs(ctx sdk.Context, hash sdk.Hash) []*types.Log {
	return k.CSDB.WithContext(ctx).GetLogs(hash)
}

func (k Keeper) DoContractCreate(ctx sdk.Context, msg types.MsgContractCreate) (err sdk.Error) {
	acc := k.ak.GetAccount(ctx, msg.From)
	if acc == nil {
		return sdk.ErrInvalidAddress(fmt.Sprintf("account %s does not exist", msg.From.String()))
	}

	contractAddr := common.CreateAddress(msg.From, acc.GetSequence())
	fmt.Fprintf(os.Stderr, fmt.Sprintf("contractAddr = %v\n", contractAddr.String()))
	contractAcc := k.ak.GetAccount(ctx, contractAddr)
	if contractAcc != nil {
		return types.ErrContractAddressCollision()
	}

	balanceEnough := false
	coins := acc.GetCoins()
	for _, coin := range coins {
		if coin.IsGTE(msg.Amount) {
			balanceEnough = true
		}
	}

	if balanceEnough == false {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("balace not enouth, amount=%v, account'balance=%v", msg.Amount, acc.GetCoins()))
	}

	codeHash := tmcrypto.Sha256(msg.Code)

	// create account
	contractAcc = k.ak.NewAccountWithAddress(ctx, contractAddr.Bytes())
	contractAcc.SetCodeHash(codeHash)
	k.ak.SetAccount(ctx, contractAcc)

	// transfer
	k.bk.SendCoins(ctx, msg.From, contractAddr.Bytes(), sdk.NewCoins(msg.Amount))

	// store code
	_, found := k.GetContractCode(ctx, codeHash)
	if !found {
		k.setContractCode(ctx, codeHash, msg.Code)
	}

	return nil
}
