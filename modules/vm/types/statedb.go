package types

import (
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/netcloth/netcloth-chain/modules/auth"
	sdk "github.com/netcloth/netcloth-chain/types"
)

var (
	zeroBalance = sdk.ZeroInt().BigInt()
)

type revision struct {
	id           int
	journalIndex int
}

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

	thash, bhash sdk.Hash
	txIndex      int
	logs         map[sdk.Hash][]*types.Log
	logSize      uint
	preimages    map[sdk.Hash][]byte

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memo-ized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// Journal of state modifications. This is the backbone of
	// Snapshot and RevertToSnapshot.
	journal        *journal
	validRevisions []revision
	nextRevisionID int

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
		preimages:         make(map[sdk.Hash][]byte),
		journal:           newJournal(),
	}
}

// WithContext returns a Database with an updated sdk context
func (csdb *CommitStateDB) WithContext(ctx sdk.Context) *CommitStateDB {
	csdb.ctx = ctx
	return csdb
}

// ----------------------------------------------------------------------------
// Setters
// ----------------------------------------------------------------------------

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

func (csdb *CommitStateDB) SetState(addr sdk.AccAddress, key, value sdk.Hash) {
	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		so.SetState(key, value)
	}
}

func (csdb *CommitStateDB) SetCode(addr sdk.AccAddress, code []byte) {
	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		codeHash := sdk.BytesToHash(crypto.Sha256(code))
		so.SetCode(codeHash, code)
	}
}

func (csdb *CommitStateDB) AddPreimage(hash sdk.Hash, preimage []byte) {
	if _, ok := csdb.preimages[hash]; !ok {
		csdb.journal.append(addPreimageChange{hash: hash})

		pi := make([]byte, len(preimage))
		copy(pi, preimage)
		csdb.preimages[hash] = pi
	}
}

func (csdb *CommitStateDB) AddRefund(gas uint64) {
	csdb.journal.append(refundChange{prev: csdb.refund})
	csdb.refund += gas
}

func (csdb *CommitStateDB) Suicide(addr sdk.AccAddress) bool {
	so := csdb.getStateObject(addr)
	if so == nil {
		return false
	}

	csdb.journal.append(suicideChange{
		account:     &addr,
		prev:        so.suicided,
		prevBalance: sdk.NewIntFromBigInt(so.Balance()),
	})

	so.markSuicided()
	//TODO: set balance 0
	so.account.Coins = sdk.Coins{sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(0))}

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

func (csdb *CommitStateDB) GetCodeHash(addr sdk.AccAddress) sdk.Hash {
	so := csdb.getStateObject(addr)
	if so == nil {
		return sdk.Hash{}
	}

	return sdk.BytesToHash(so.CodeHash())
}

func (csdb *CommitStateDB) GetCommittedState(addr sdk.AccAddress, hash sdk.Hash) sdk.Hash {
	so := csdb.getStateObject(addr)
	if so != nil {
		return so.GetCommittedState(hash)
	}

	return sdk.Hash{}
}

func (csdb *CommitStateDB) GetRefund() uint64 {
	return csdb.refund
}

func (csdb *CommitStateDB) Preimages() map[sdk.Hash][]byte {
	return csdb.preimages
}

func (csdb *CommitStateDB) HasSuicide(addr sdk.AccAddress) bool {
	so := csdb.getStateObject(addr)
	if so != nil {
		return so.suicided
	}

	return false
}

// ----------------------------------------------------------------------------
// Persistence
// ----------------------------------------------------------------------------

func (csdb *CommitStateDB) Commit(deleteEmptyObjects bool) (root sdk.Hash, err error) {
	defer csdb.clearJournalAndRefund()

	// remove dirty state object entries based on the journal
	for addr := range csdb.journal.dirties {
		csdb.stateObjectsDirty[addr] = struct{}{}
	}

	// set the state objects
	for addr, so := range csdb.stateObjects {
		_, isDrity := csdb.stateObjectsDirty[addr]

		switch {
		case so.suicided || (isDrity && deleteEmptyObjects && so.empty()):
			csdb.deleteStateObject(so)

		case isDrity:
			if so.code != nil && so.dirtyCode {
				so.commitCode()
				so.dirtyCode = false
			}

			// update the object in the KVStore
			csdb.updateStateObject(so)
		}
		delete(csdb.stateObjectsDirty, addr)
	}

	return
}

// ClearStateObjects clears cache of state objects to handle account changes outside of the EVM
func (csdb *CommitStateDB) ClearStateObjects() {
	csdb.stateObjects = make(map[string]*stateObject)
	csdb.stateObjectsDirty = make(map[string]struct{})
}

func (csdb *CommitStateDB) updateStateObject(so *stateObject) {
	csdb.ak.SetAccount(csdb.ctx, so.account)
}

func (csdb *CommitStateDB) deleteStateObject(so *stateObject) {
	so.deleted = true
	csdb.ak.RemoveAccount(csdb.ctx, so.account)
}

func (csdb *CommitStateDB) clearJournalAndRefund() {
	csdb.journal = newJournal()
	csdb.validRevisions = csdb.validRevisions[:0]
	csdb.refund = 0
}

// ----------------------------------------------------------------------------
// Snapshotting
// ----------------------------------------------------------------------------

// Snapshot returns an identifier for the current revision of the state.
func (csdb *CommitStateDB) Snapshot() int {
	id := csdb.nextRevisionID
	csdb.nextRevisionID++

	csdb.validRevisions = append(
		csdb.validRevisions,
		revision{
			id:           id,
			journalIndex: csdb.journal.length(),
		},
	)

	return id
}

// RevertToSnapshot reverts all state changes made since the given revision.
func (csdb *CommitStateDB) RevertToSnapshot(revID int) {
	idx := sort.Search(len(csdb.validRevisions), func(i int) bool {
		return csdb.validRevisions[i].id >= revID
	})

	if idx == len(csdb.validRevisions) || csdb.validRevisions[idx].id != revID {
		panic(fmt.Errorf("revision ID %v cannot be reverted", revID))
	}

	snapshot := csdb.validRevisions[idx].journalIndex

	// replay the journal to undo changes and remove invalidated snapshots
	csdb.journal.revert(csdb, snapshot)
	csdb.validRevisions = csdb.validRevisions[:idx]
}

// ----------------------------------------------------------------------------
// Auxiliary
// ----------------------------------------------------------------------------

func (csdb *CommitStateDB) Empty(addr sdk.AccAddress) bool {
	so := csdb.getStateObject(addr)
	return so == nil || so.empty()
}

func (csdb *CommitStateDB) Exist(addr sdk.AccAddress) bool {
	return csdb.getStateObject(addr) != nil
}

// Error returns the first non-nil error the StateDB encountered.
func (csdb *CommitStateDB) Error() error {
	return csdb.dbErr
}
