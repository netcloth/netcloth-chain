package types

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/tendermint/tendermint/crypto"

	authexported "github.com/netcloth/netcloth-chain/modules/auth/exported"
	"github.com/netcloth/netcloth-chain/modules/auth/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

var (
	_ StateObject = (*stateObject)(nil)

	emptyCodeHash = sdk.Hash{}
)

var ripemd = sdk.HexToAddress("0000000000000000000000000000000000000003")

type (
	// StateObject interface for interacting with state object
	StateObject interface {
		GetCommittedState(key sdk.Hash) sdk.Hash
		GetState(key sdk.Hash) sdk.Hash
		SetState(key, value sdk.Hash)

		Code() []byte
		SetCode(codeHash sdk.Hash, code []byte)
		CodeHash() []byte // codeHash = crypto.Sha256(Code)

		AddBalance(amount *big.Int)
		SubBalance(amount *big.Int)
		SetBalance(amount *big.Int)

		Balance() *big.Int
		//ReturnGas(gas *big.Int)
		Address() sdk.AccAddress

		SetNonce(nonce uint64)
		Nonce() uint64
	}

	// stateObject represents an NCH account which is being modified
	//
	// The usage pattern is as follows:
	// First you need to obtain a state object.
	// Account values can be accessed and modified through the object.
	// Finally, call CommitTrie to write the modified storage trie into a database.
	stateObject struct {
		address sdk.AccAddress
		stateDB *CommitStateDB
		account *types.BaseAccount

		dbErr error

		code sdk.Code // contract bytecode

		originStorage sdk.Storage // Storage cache of original entries to dedup rewrites
		dirtyStorage  sdk.Storage // Storage entries that need to be flushed to disk

		// cache flags
		//
		// When an object is marked suicided, it will be deleted from the trie during the "update" phase of the state transition.
		dirtyCode bool // true if the code was updated
		suicided  bool
		deleted   bool
	}
)

func newObject(accProto authexported.Account, csdb *CommitStateDB) *stateObject {
	acc, ok := accProto.(*types.BaseAccount)
	if !ok {
		panic(fmt.Sprintf("invalid account type for state object: %T", accProto))
	}

	if acc.CodeHash == nil {
		acc.CodeHash = emptyCodeHash.Bytes()
	}

	return &stateObject{
		stateDB:       csdb,
		account:       acc,
		address:       acc.Address,
		originStorage: make(sdk.Storage),
		dirtyStorage:  make(sdk.Storage),
	}
}

// ----------------------------------------------------------------------------
// Setters
// ----------------------------------------------------------------------------

// SetState updates a value in account storage. Note, the key will be prefixed
// with the address of the state object.
func (so *stateObject) SetState(key, value sdk.Hash) {
	prev := so.GetState(key)
	if prev == value {
		return
	}

	prefixKey := so.GetStorageByAddressKey(key.Bytes())

	so.setState(prefixKey, value)
}

func (so *stateObject) setState(key, value sdk.Hash) {
	so.dirtyStorage[key] = value
}

// SetCode
func (so *stateObject) SetCode(codeHash sdk.Hash, code []byte) {
	so.setCode(codeHash, code)
}

func (so *stateObject) setCode(codeHash sdk.Hash, code []byte) {
	so.code = code
	so.account.CodeHash = codeHash.Bytes()
	so.dirtyCode = true
}

// AddBalance adds an amount to a state object's balance. It is used to add
// funds to the destination account of a transfer.
func (so *stateObject) AddBalance(amount *big.Int) {
	amt := sdk.NewIntFromBigInt(amount)

	if amt.Sign() == 0 {
		if so.empty() {
			so.touch()
		}
		return
	}

	newBalance := so.account.Balance().Add(amt)
	so.SetBalance(newBalance.BigInt())
}

// SubBalance removes an amount from the stateObject's balance. It is used to
// remove funds from the origin account of a transfer.
func (so *stateObject) SubBalance(amount *big.Int) {
	amt := sdk.NewIntFromBigInt(amount)

	if amt.Sign() == 0 {
		return
	}

	newBalance := so.account.Balance().Sub(amt)
	so.SetBalance(newBalance.BigInt())
}

// SetBalance sets the state object's balance.
func (so *stateObject) SetBalance(amount *big.Int) {
	amt := sdk.NewIntFromBigInt(amount)

	so.setBalance(amt)
}

func (so *stateObject) setBalance(amount sdk.Int) {
	so.account.SetBalance(amount)
}

// SetNonce sets the state object's nonce (sequence number).
func (so *stateObject) SetNonce(nonce uint64) {
	so.setNonce(nonce)
}

func (so *stateObject) setNonce(nonce uint64) {
	so.account.Sequence = nonce
}

// setError remembers the first non-nil error it is called with.
func (so *stateObject) setError(err error) {
	if so.dbErr == nil {
		so.dbErr = err
	}
}

func (so *stateObject) markSuicided() {
	so.suicided = true
}

func (so *stateObject) commitCode() {
	ctx := so.stateDB.Ctx
	store := ctx.KVStore(so.stateDB.codeKey)
	store.Set(so.CodeHash(), so.code)
}

func (so *stateObject) commitState() {
	ctx := so.stateDB.Ctx
	store := ctx.KVStore(so.stateDB.storageKey)

	for key, value := range so.dirtyStorage {
		delete(so.dirtyStorage, key)

		// skip no-op changes, persist actual changes
		if value == so.originStorage[key] {
			continue
		}

		so.originStorage[key] = value

		// delete empty values
		if (value == sdk.Hash{}) {
			store.Delete(key.Bytes())
			continue
		}

		store.Set(key.Bytes(), value.Bytes())
	}

	// TODO: Set the account (storage) root (but we probably don't need this)
}

// ----------------------------------------------------------------------------
// Getters
// ----------------------------------------------------------------------------

// Address returns the address of the state object
func (so stateObject) Address() sdk.AccAddress {
	return so.address
}

// Balance returns the state object's current balance
func (so *stateObject) Balance() *big.Int {
	return so.account.Balance().BigInt()
}

// CodeHash returns the state object's code hash
func (so *stateObject) CodeHash() []byte {
	return so.account.CodeHash
}

// Nonce returns the state object's current nonce(sequence number)
func (so *stateObject) Nonce() uint64 {
	return so.account.Sequence
}

// Code returns the contract code associated with this object
func (so *stateObject) Code() []byte {
	// TODO
	if so.code != nil {
		return so.code
	}

	if bytes.Equal(so.CodeHash(), emptyCodeHash.Bytes()) {
		return nil
	}

	ctx := so.stateDB.Ctx
	store := ctx.KVStore(so.stateDB.codeKey)
	code := store.Get(so.CodeHash())

	if len(code) == 0 {
		so.setError(fmt.Errorf("failed to get code hash %x for address: %x", so.CodeHash(), so.address))
	}

	so.code = code
	return code
}

// GetState retrieves a value from the account storage trie. Note, the key will be prefixed with the address of the state object
func (so *stateObject) GetState(key sdk.Hash) sdk.Hash {
	prefixKey := so.GetStorageByAddressKey(key.Bytes())

	// if we have a dirty value for this state entry, return it
	value, dirty := so.dirtyStorage[prefixKey]
	if dirty {
		return value
	}

	// otherwise return the entry's original valeu
	return so.GetCommittedState(prefixKey)
}

// GetCommittedState retrieves a value from the committed account storage trie.
// Note, the key will be prefixed with the address of the state object.
func (so *stateObject) GetCommittedState(key sdk.Hash) sdk.Hash {
	prefixKey := so.GetStorageByAddressKey(key.Bytes())

	// if we have the original value cached, return that
	value, cached := so.originStorage[prefixKey]
	if cached {
		return value
	}

	// otherwise load the value from KVStore
	ctx := so.stateDB.Ctx
	store := ctx.KVStore(so.stateDB.storageKey)
	rawValue := store.Get(prefixKey.Bytes())

	if len(rawValue) > 0 {
		value.SetBytes(rawValue)
	}

	so.originStorage[prefixKey] = value
	return value
}

// ----------------------------------------------------------------------------
// Auxiliary
// ----------------------------------------------------------------------------

// empty returns whether the account is considered empty.
func (so *stateObject) empty() bool {
	return so.account.Sequence == 0 &&
		so.account.Balance().Sign() == 0 &&
		bytes.Equal(so.account.CodeHash, emptyCodeHash.Bytes())
}

func (so *stateObject) touch() {
	so.stateDB.journal.append(touchChange{
		account: &so.address,
	})

	if so.address.Equals(ripemd) {
		so.stateDB.journal.dirty(so.address)
	}
}

// GetStorageByAddressKey returns a hash of the composite key for a state
// object's storage prefixed with it's address.
func (so stateObject) GetStorageByAddressKey(key []byte) sdk.Hash {
	prefix := so.Address().Bytes()
	compositeKey := make([]byte, len(prefix)+len(key))

	copy(compositeKey, prefix)
	copy(compositeKey[len(prefix):], key)

	h := sdk.Hash{}
	h.SetBytes(crypto.Sha256(compositeKey))
	return h
}
