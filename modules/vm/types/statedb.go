package types

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/tendermint/tendermint/crypto"

	"github.com/netcloth/netcloth-chain/modules/auth"
	"github.com/netcloth/netcloth-chain/modules/vm/common"
	sdk "github.com/netcloth/netcloth-chain/types"
)

var (
	zeroBalance = sdk.ZeroInt().BigInt()
)

type CommitStateDB struct {
	ctx sdk.Context

	ak         auth.AccountKeeper
	storageKey sdk.StoreKey
	codeKey    sdk.StoreKey

	// maps that hold 'live' objects, which will get modified while processing a
	// state transition
	stateObjects      map[string]*stateObject
	stateObjectsDirty map[string]struct{}

	// The refund counter, also used by state transitioning.
	refund uint64

	thash, bhash common.Hash
	txIndex      int
	// logs
	logSize   uint
	preimages map[common.Hash][]byte

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memo-ized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	lock sync.Mutex
}

// NewCommitStateDB returns a reference to a newly initialized CommitStateDB
// which implements Geth's state.StateDB interface.
//
// CONTRACT: Stores used for state must be cache-wrapped as the ordering of the
// key/value space matters in determining the merkle root.
func NewCommitStateDB(ctx sdk.Context, ak auth.AccountKeeper, storageKey, codeKey sdk.StoreKey) *CommitStateDB {
	return &CommitStateDB{
		ctx:               ctx,
		ak:                ak,
		storageKey:        storageKey,
		codeKey:           codeKey,
		stateObjects:      make(map[string]*stateObject),
		stateObjectsDirty: make(map[string]struct{}),
		preimages:         make(map[common.Hash][]byte),
	}
}

func (csdb *CommitStateDB) SetBalance(addr sdk.AccAddress, amount *big.Int) {
	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		so.SetBalance(amount)
	}
}

func (csdb *CommitStateDB) AddBalance(addr sdk.AccAddress, amount *big.Int) {
	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		so.AddBalance(amount)
	}
}

func (csdb *CommitStateDB) SubBalance(addr sdk.AccAddress, amount *big.Int) {
	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		so.SubBalance(amount)
	}
}

func (csdb *CommitStateDB) SetNonce(addr sdk.AccAddress, nonce uint64) {
	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		so.SetNonce(nonce)
	}

}

func (csdb *CommitStateDB) SetState(addr sdk.AccAddress, key, value common.Hash) {
	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		so.SetState(key, value)
	}
}

func (csdb *CommitStateDB) SetCode(addr sdk.AccAddress, code []byte) {
	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		codeHash := common.BytesToHash(crypto.Sha256(code))
		so.SetCode(codeHash, code)
	}
}

func (csdb *CommitStateDB) Suicide(addr sdk.AccAddress) bool {
	so := csdb.getStateObject(addr)
	if so == nil {
		return false
	}

	so.markSuicided()
	//TODO: set balance 0
	so.account.Coins = sdk.Coins{sdk.NewCoin(sdk.NativeTokenName, 0)}

	return true
}

func (csdb *CommitStateDB) GetOrNewStateObject(addr sdk.AccAddress) StateObject {
	so := csdb.getStateObject(addr)
	if so == nil || so.deleted {
		so, _ = csdb.createObject(addr)
	}

	return so
}

// createObject creates a new state object. If there is an existing account with
// the given address, it is overwritten and returned as the second return value.
func (csdb *CommitStateDB) createObject(addr sdk.AccAddress) (newObj, prevObj *stateObject) {
	prevObj = csdb.getStateObject(addr)

	acc := csdb.ak.NewAccountWithAddress(csdb.ctx, addr)
	newObj = newObject(acc)
	newObj.SetNonce(0)

	if prevObj == nil {
		// TODO
	} else {
		// TODO
	}
	csdb.setStateObject(newObj)
	return newObj, prevObj

}

// getStateObject attempts to retrieve a state object given by the address.
// Returns nil and sets an error if not found.
func (csdb *CommitStateDB) getStateObject(addr sdk.AccAddress) (stateObject *stateObject) {
	// prefer "live" (cached) objects
	if so := csdb.stateObjects[addr.String()]; so != nil {
		if so.deleted {
			return nil
		}

		return so
	}

	// otherwise, attempt to fetch the account from the account mapper
	acc := csdb.ak.GetAccount(csdb.ctx, addr.Bytes())
	if acc == nil {
		csdb.setError(fmt.Errorf("no account found for address: %s", addr.String()))
		return nil
	}

	// insert the state object into the live set
	so := newObject(acc)
	csdb.setStateObject(so)

	return so
}

// CreateAccount explicitly creates a state object. If a state object with the address
// already exists the balance is carried over to the new account.
//
// CreateAccount is called during the EVM CREATE operation. The situation might arise that
// a contract does the following:
//
//   1. sends funds to sha(account ++ (nonce + 1))
//   2. tx_create(sha(account ++ nonce)) (note that this gets the address of 1)
//
// Carrying over the balance ensures that Ether doesn't disappear.
func (csdb *CommitStateDB) CreateAccount(addr sdk.AccAddress) {
	newObj, prev := csdb.createObject(addr)
	if prev != nil {
		newObj.setBalance(prev.account.Balance())
	}
}

// WithContext returns a Database with an updated sdk context
func (csdb *CommitStateDB) WithContext(ctx sdk.Context) *CommitStateDB {
	csdb.ctx = ctx
	return csdb
}

func (csdb *CommitStateDB) setStateObject(so *stateObject) {
	csdb.stateObjects[so.Address().String()] = so
}

// setError remembers the first non-nil error it is called with.
func (csdb *CommitStateDB) setError(err error) {
	if csdb.dbErr == nil {
		csdb.dbErr = err
	}
}

// ----------------------------------------------------------------------------
// Getters
// ----------------------------------------------------------------------------

func (csdb *CommitStateDB) GetBalance(addr sdk.AccAddress) *big.Int {
	so := csdb.getStateObject(addr)
	if so != nil {
		return so.Balance()
	}
	return zeroBalance
}

func (csdb *CommitStateDB) GetNonce(addr sdk.AccAddress) uint64 {
	so := csdb.getStateObject(addr)
	if so != nil {
		return so.Nonce()
	}
	return 0
}

func (csdb *CommitStateDB) TxIndex() int {
	return csdb.txIndex
}

func (csdb *CommitStateDB) GetCode(addr sdk.AccAddress) []byte {
	so := csdb.getStateObject(addr)
	if so != nil {
		return so.Code()
	}

	return nil
}

func (csdb *CommitStateDB) GetCodeSize(addr sdk.AccAddress) int {
	so := csdb.getStateObject(addr)
	if so == nil {
		return 0
	}

	if so.code != nil {
		return len(so.code)
	}

	return len(so.Code())
}

func (csdb *CommitStateDB) GetCodeHash(addr sdk.AccAddress) common.Hash {
	so := csdb.getStateObject(addr)
	if so == nil {
		return common.Hash{}
	}

	return common.BytesToHash(so.CodeHash())
}

///////////////////
func (csdb *CommitStateDB) Empty(addr sdk.AccAddress) bool {
	so := csdb.getStateObject(addr)
	return so == nil || so.empty()
}

func (csdb *CommitStateDB) Exist(addr sdk.AccAddress) bool {
	return csdb.getStateObject(addr) != nil
}
