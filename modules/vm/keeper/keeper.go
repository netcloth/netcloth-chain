package keeper

import (
	"fmt"
	"os"

	"github.com/netcloth/netcloth-chain/modules/vm/common"

	tmcrypto "github.com/tendermint/tendermint/crypto"

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

func (k Keeper) GetContractCode(ctx sdk.Context, contractAddr, codeHash []byte) (code []byte, found bool) {
	store := ctx.KVStore(k.storeKey)
	code = store.Get(types.GetContractCodeKey(contractAddr, codeHash))
	return code, code != nil
}

func (k Keeper) setContractCode(ctx sdk.Context, contractAddr, codeHash, code []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetContractCodeKey(contractAddr, codeHash), code)
}

func (k Keeper) DoContractCreate(ctx sdk.Context, msg types.MsgContractCreate) (err sdk.Error) {
	acc := k.ak.GetAccount(ctx, msg.From)
	if acc == nil {
		return sdk.ErrInvalidAddress(fmt.Sprintf("account %s does not exist", msg.From.String()))
	}

	contractAddr := common.CreateAddress(msg.From, acc.GetSequence())
	fmt.Fprintf(os.Stderr, fmt.Sprintf("contractAddr = %v\n", contractAddr.String()))

	codeHash := tmcrypto.Sha256(msg.Code)
	_, found := k.GetContractCode(ctx, contractAddr.Bytes(), codeHash)
	if found {
		return types.ErrContractExist("contract already exist")
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

	// create account
	contractAcc := k.ak.NewAccountWithAddress(ctx, contractAddr.Bytes())
	contractAcc.SetCodeHash(codeHash)
	k.ak.SetAccount(ctx, contractAcc)

	// transfer
	k.bk.SendCoins(ctx, msg.From, contractAddr.Bytes(), sdk.NewCoins(msg.Amount))

	// store code
	k.setContractCode(ctx, contractAddr.Bytes(), codeHash, msg.Code)

	return nil
}
