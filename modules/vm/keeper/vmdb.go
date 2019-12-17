package keeper

import (
	"fmt"
	"math/big"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/netcloth/netcloth-chain/modules/vm/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

var DefaultVMDenom = "unch"

func CoinsFromBigInt(v *big.Int) sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin(DefaultVMDenom, sdk.NewInt(v.Int64())))
}

var _ StateDB = (*VMDB)(nil)

type VMDB struct {
	ctx sdk.Context
	ak  types.AccountKeeper
	bk  types.BankKeeper
	vk  Keeper

	stateObjects      map[string]*stateObject
	stateObjectsDirty map[string]struct{}

	dbErr error
}

func NewVMDB() *VMDB {
	return &VMDB{}
}

func (db *VMDB) WithContext(ctx sdk.Context) *VMDB {
	db.ctx = ctx
	return db
}

func (db *VMDB) getStateObject(addr sdk.AccAddress) (stateObject *stateObject) {
	if so := db.stateObjects[addr.String()]; so != nil {
		if so.deleted {
			return nil
		}

		return so
	}

	acc := db.ak.GetAccount(db.ctx, addr)
	if acc == nil {
		db.setError(fmt.Errorf("no account found for address: %s", addr.String()))
		return nil
	}

	code := db.vk.GetCode(addr)

	so := newObject(acc, db)
	csdb.setStateObject(so)

	return so
}

// createObject creates a new state object. If there is an existing account with
// the given address, it is overwritten and returned as the second return value.
func (db *VMDB) createObject(addr sdk.AccAddress) (newObj, prevObj *stateObject) {
	prevObj = db.getStateObject(addr)

	acc := db.ak.NewAccountWithAddress(csdb.Ctx, addr)
	newObj = newObject(acc, csdb)
	newObj.SetNonce(0)

	if prevObj == nil {
		csdb.journal.append(createObjectChange{account: &addr})
	} else {
		csdb.journal.append(resetObjectChange{prev: prevObj})
	}
	csdb.setStateObject(newObj)
	return newObj, prevObj

}

func (db *VMDB) CreateAccount(addr sdk.AccAddress) {
	newObj, prev := db.createObject(addr)
	if prev != nil {
		newObj.setBalance(prev.account.Balance())
	}
}

func (db *VMDB) SubBalance(addr sdk.AccAddress, amt *big.Int) {
	db.bk.SubtractCoins(db.ctx, addr, CoinsFromBigInt(amt))
}

func (db *VMDB) AddBalance(addr sdk.AccAddress, amt *big.Int) {
	db.bk.AddCoins(db.ctx, addr, CoinsFromBigInt(amt))
}

func (db *VMDB) GetBalance(addr sdk.AccAddress) *big.Int {
	acc := db.ak.GetAccount(db.ctx, addr)
	if acc != nil {
		for _, coin := range acc.GetCoins() {
			if coin.Denom == DefaultVMDenom {
				return coin.Amount.BigInt()
			}
		}
	}

	return sdk.NewInt(0).BigInt()
}

func (db *VMDB) GetNonce(addr sdk.AccAddress) uint64 {
	acc := db.ak.GetAccount(db.ctx, addr)
	if acc != nil {
		return acc.GetSequence()
	}
	return 0 //TODO check
}

func (db *VMDB) SetNonce(addr sdk.AccAddress, nonce uint64) {
	acc := db.ak.GetAccount(db.ctx, addr)
	if acc != nil {
		if acc.GetSequence()+1 == nonce {
			acc.SetSequence(nonce)
			db.ak.SetAccount(db.ctx, acc)
		}
	}
}

func (db *VMDB) GetCodeHash(addr sdk.AccAddress) sdk.Hash {
	acc := db.ak.GetAccount(db.ctx, addr)
	if acc != nil {
		return sdk.BytesToHash(acc.GetCodeHash())
	}
	return sdk.Hash{} // TODO check: use empty hash?
}

func (db *VMDB) GetCode(addr sdk.AccAddress) []byte {
	return nil
}

func (db *VMDB) SetCode(addr sdk.AccAddress, code []byte) {
	panic("implement me")
}

func (db *VMDB) GetCodeSize(addr sdk.AccAddress) int {
	panic("implement me")
}

func (db *VMDB) AddRefund(uint64) {
	panic("implement me")
}

func (db *VMDB) SubRefund(uint64) {
	panic("implement me")
}

func (db *VMDB) GetRefund() uint64 {
	panic("implement me")
}

func (db *VMDB) GetCommittedState(addr sdk.AccAddress, s sdk.Hash) sdk.Hash {
	panic("implement me")
}

func (db *VMDB) GetState(addr sdk.AccAddress, k sdk.Hash) sdk.Hash {
	panic("implement me")
}

func (db *VMDB) SetState(addr sdk.AccAddress, k sdk.Hash, v sdk.Hash) {
	panic("implement me")
}

func (db *VMDB) Suicide(addr sdk.AccAddress) bool {
	panic("implement me")
}

func (db *VMDB) HasSuicided(addr sdk.AccAddress) bool {
	panic("implement me")
}

func (db *VMDB) Exist(addr sdk.AccAddress) bool {
	panic("implement me")
}

func (db *VMDB) Empty(addr sdk.AccAddress) bool {
	panic("implement me")
}

func (db *VMDB) RevertToSnapshot(int) {
	panic("implement me")
}

func (db *VMDB) Snapshot() int {
	panic("implement me")
}

func (db *VMDB) AddLog(*ethtypes.Log) {
	panic("implement me")
}

func (db *VMDB) AddPreimage(sdk.Hash, []byte) {
	panic("implement me")
}

func (db *VMDB) ForEachStorage(sdk.AccAddress, func(sdk.Hash, sdk.Hash) bool) error {
	panic("implement me")
}

func (db *VMDB) setError(err error) {
	if db.dbErr == nil {
		db.dbErr = err
	}
}
