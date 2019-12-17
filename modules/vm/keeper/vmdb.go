package keeper

import (
	"fmt"
	"math/big"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/netcloth/netcloth-chain/modules/vm/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

var _ VMDB = (*VM2DB)(nil)

type VM2DB struct {
	ctx sdk.Context
	ak  types.AccountKeeper
	bk  types.BankKeeper
	vk  Keeper

	stateObjects      map[string]*stateObject
	stateObjectsDirty map[string]struct{}

	dbErr error
}

func NewVMDB() *VM2DB {
	return &VM2DB{}
}

func (db *VM2DB) WithContext(ctx sdk.Context) *VM2DB {
	db.ctx = ctx
	return db
}

func (db *VM2DB) getStateObject(addr sdk.AccAddress) (stateObject *stateObject) {
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
func (db *VM2DB) createObject(addr sdk.AccAddress) (newObj, prevObj *stateObject) {
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

func (db *VM2DB) CreateAccount(addr sdk.AccAddress) {
	newObj, prev := db.createObject(addr)
	if prev != nil {
		newObj.setBalance(prev.account.Balance())
	}
}

func (db *VM2DB) SubBalance(addr sdk.AccAddress, amt *big.Int) {
	db.bk.SubtractCoins(db.ctx, addr, CoinsFromBigInt(amt))
}

func (db *VM2DB) AddBalance(addr sdk.AccAddress, amt *big.Int) {
	db.bk.AddCoins(db.ctx, addr, CoinsFromBigInt(amt))
}

func (db *VM2DB) GetBalance(addr sdk.AccAddress) *big.Int {
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

func (db *VM2DB) GetNonce(addr sdk.AccAddress) uint64 {
	acc := db.ak.GetAccount(db.ctx, addr)
	if acc != nil {
		return acc.GetSequence()
	}
	return 0 //TODO check
}

func (db *VM2DB) SetNonce(addr sdk.AccAddress, nonce uint64) {
	acc := db.ak.GetAccount(db.ctx, addr)
	if acc != nil {
		if acc.GetSequence()+1 == nonce {
			acc.SetSequence(nonce)
			db.ak.SetAccount(db.ctx, acc)
		}
	}
}

func (db *VM2DB) GetCodeHash(addr sdk.AccAddress) sdk.Hash {
	acc := db.ak.GetAccount(db.ctx, addr)
	if acc != nil {
		return sdk.BytesToHash(acc.GetCodeHash())
	}
	return sdk.Hash{} // TODO check: use empty hash?
}

func (db *VM2DB) GetCode(addr sdk.AccAddress) []byte {
	return nil
}

func (db *VM2DB) SetCode(addr sdk.AccAddress, code []byte) {
	panic("implement me")
}

func (db *VM2DB) GetCodeSize(addr sdk.AccAddress) int {
	panic("implement me")
}

func (db *VM2DB) AddRefund(uint64) {
	panic("implement me")
}

func (db *VM2DB) SubRefund(uint64) {
	panic("implement me")
}

func (db *VM2DB) GetRefund() uint64 {
	panic("implement me")
}

func (db *VM2DB) GetCommittedState(addr sdk.AccAddress, s sdk.Hash) sdk.Hash {
	panic("implement me")
}

func (db *VM2DB) GetState(addr sdk.AccAddress, k sdk.Hash) sdk.Hash {
	panic("implement me")
}

func (db *VM2DB) SetState(addr sdk.AccAddress, k sdk.Hash, v sdk.Hash) {
	panic("implement me")
}

func (db *VM2DB) Suicide(addr sdk.AccAddress) bool {
	panic("implement me")
}

func (db *VM2DB) HasSuicided(addr sdk.AccAddress) bool {
	panic("implement me")
}

func (db *VM2DB) Exist(addr sdk.AccAddress) bool {
	panic("implement me")
}

func (db *VM2DB) Empty(addr sdk.AccAddress) bool {
	panic("implement me")
}

func (db *VM2DB) RevertToSnapshot(int) {
	panic("implement me")
}

func (db *VM2DB) Snapshot() int {
	panic("implement me")
}

func (db *VM2DB) AddLog(*ethtypes.Log) {
	panic("implement me")
}

func (db *VM2DB) AddPreimage(sdk.Hash, []byte) {
	panic("implement me")
}

func (db *VM2DB) ForEachStorage(sdk.AccAddress, func(sdk.Hash, sdk.Hash) bool) error {
	panic("implement me")
}

func (db *VM2DB) setError(err error) {
	if db.dbErr == nil {
		db.dbErr = err
	}
}
