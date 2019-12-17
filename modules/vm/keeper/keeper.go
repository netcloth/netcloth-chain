package keeper

import (
	"fmt"
	"math/big"
	"os"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/netcloth/netcloth-chain/modules/auth/exported"

	tmcrypto "github.com/tendermint/tendermint/crypto"

	"github.com/netcloth/netcloth-chain/modules/params"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/vm/common"
	"github.com/netcloth/netcloth-chain/modules/vm/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

var _ VMDB = (*Keeper)(nil)

type Keeper struct {
	ctx        sdk.Context
	storeKey   sdk.StoreKey
	storeTKey  sdk.StoreKey
	cdc        *codec.Codec
	paramstore params.Subspace
	ak         types.AccountKeeper
	bk         types.BankKeeper

	codespace sdk.CodespaceType
}

var VMDenom = "unch"

func CoinsFromBigInt(v *big.Int) sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin(VMDenom, sdk.NewInt(v.Int64())))
}

func (k Keeper) CreateAccount(addr sdk.AccAddress) {
	k.ak.NewAccountWithAddress(k.ctx, addr)
}

func (k Keeper) SubBalance(addr sdk.AccAddress, amt *big.Int) {
	k.bk.SubtractCoins(k.ctx, addr, CoinsFromBigInt(amt))
}

func (k Keeper) AddBalance(addr sdk.AccAddress, amt *big.Int) {
	k.bk.AddCoins(k.ctx, addr, CoinsFromBigInt(amt))
}

func (k Keeper) GetBalance(addr sdk.AccAddress) *big.Int {
	coins := k.bk.GetCoins(k.ctx, addr)
	for _, coin := range coins {
		if coin.Denom == VMDenom {
			return coin.Amount.BigInt()
		}
	}

	return sdk.NewInt(0).BigInt() //TODO
}

func (k Keeper) GetNonce(addr sdk.AccAddress) uint64 {
	acc := k.ak.GetAccount(k.ctx, addr)
	if acc != nil {
		return acc.GetSequence()
	}
	return 0 //TODO
}

func (k Keeper) SetNonce(addr sdk.AccAddress, nonce uint64) {
	acc := k.ak.GetAccount(k.ctx, addr)
	if acc != nil {
		if acc.GetSequence()+1 == nonce {
			acc.SetSequence(nonce)
			k.ak.SetAccount(k.ctx, acc)
		}
	}
}

func (k Keeper) GetCodeHash(sdk.AccAddress) sdk.Hash {

}

func (k Keeper) GetCode(sdk.AccAddress) []byte {
	panic("implement me")
}

func (k Keeper) SetCode(sdk.AccAddress, []byte) {
	panic("implement me")
}

func (k Keeper) GetCodeSize(sdk.AccAddress) int {
	panic("implement me")
}

func (k Keeper) AddRefund(uint64) {
	panic("implement me")
}

func (k Keeper) SubRefund(uint64) {
	panic("implement me")
}

func (k Keeper) GetRefund() uint64 {
	panic("implement me")
}

func (k Keeper) GetCommittedState(sdk.AccAddress, sdk.Hash) sdk.Hash {
	panic("implement me")
}

func (k Keeper) GetState(sdk.AccAddress, sdk.Hash) sdk.Hash {
	panic("implement me")
}

func (k Keeper) SetState(sdk.AccAddress, sdk.Hash, sdk.Hash) {
	panic("implement me")
}

func (k Keeper) Suicide(sdk.AccAddress) bool {
	panic("implement me")
}

func (k Keeper) HasSuicided(sdk.AccAddress) bool {
	panic("implement me")
}

func (k Keeper) Exist(sdk.AccAddress) bool {
	panic("implement me")
}

func (k Keeper) Empty(sdk.AccAddress) bool {
	panic("implement me")
}

func (k Keeper) RevertToSnapshot(int) {
	panic("implement me")
}

func (k Keeper) Snapshot() int {
	panic("implement me")
}

func (k Keeper) AddLog(*ethtypes.Log) {
	panic("implement me")
}

func (k Keeper) AddPreimage(sdk.Hash, []byte) {
	panic("implement me")
}

func (k Keeper) ForEachStorage(sdk.AccAddress, func(sdk.Hash, sdk.Hash) bool) error {
	panic("implement me")
}

func (k Keeper) GetVMObject(ctx sdk.Context, addr sdk.AccAddress) exported.Account {
	store := ctx.KVStore(ak.key)
	bz := store.Get(types.AddressStoreKey(addr))
	if bz == nil {
		return nil
	}
	acc := ak.decodeAccount(bz)
	return acc
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(cdc *codec.Codec, key, tkey sdk.StoreKey, codespace sdk.CodespaceType, paramstore params.Subspace, ak types.AccountKeeper, bk types.BankKeeper) Keeper {
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

func (k *Keeper) WithContext(ctx sdk.Context) *Keeper {
	k.ctx = ctx
	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("modules/%s", types.ModuleName))
}

func (k Keeper) GetCodeByCodeHash(ctx sdk.Context, codeHash sdk.Hash) (code []byte, found bool) {
	return nil, true
}

func (k Keeper) GetCodeByAccount(ctx sdk.Context, acc sdk.AccAddress) (code []byte, found bool) {
	return nil, true
}

func (k Keeper) GetVMState(ctx sdk.Context, key sdk.Hash) sdk.Hash {
	return sdk.Hash{}
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
