package types

import (
	"fmt"
	"math/big"

	"github.com/tendermint/tendermint/crypto"

	authexported "github.com/netcloth/netcloth-chain/app/v0/auth/exported"
	"github.com/netcloth/netcloth-chain/app/v0/auth/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

var (
	_ StateObject = (*stateObject)(nil)
)

type (
	// StateObject interface for interacting with state object
	StateObject interface {
		GetCommittedState(key sdk.Hash) sdk.Hash
		GetState(key sdk.Hash) sdk.Hash
		SetState(key, value sdk.Hash)

		Code() []byte
		SetCode(codeHash sdk.Hash, code []byte)
		CodeHash() []byte

		AddBalance(amount *big.Int)
		SubBalance(amount *big.Int)
		SetBalance(amount *big.Int)

		Balance() *big.Int
		ReturnGas(gas *big.Int)
		Address() sdk.Address

		SetNonce(nonce uint64)
		Nonce() uint64
	}
	// stateObject represents an Ethereum account which is being modified.
	//
	// The usage pattern is as follows:
	// First you need to obtain a state object.
	// Account values can be accessed and modified through the object.
	// Finally, call CommitTrie to write the modified storage trie into a database.
	stateObject struct {
		address sdk.AccAddress
		stateDB *CommitStateDB
		account *types.BaseAccount

		// DB error.
		// State objects are used by the consensus core and VM which are
		// unable to deal with database-level errors. Any error that occurs
		// during a database read is memoized here and will eventually be returned
		// by StateDB.Commit.
		dbErr error

		code sdk.Code // contract bytecode, which gets set when code is loaded

		originStorage sdk.Storage // Storage cache of original entries to dedup rewrites
		dirtyStorage  sdk.Storage // Storage entries that need to be flushed to disk

		// cache flags
		//
		// When an object is marked suicided it will be delete from the trie during
		// the "update" phase of the state transition.
		dirtyCode bool
		suicided  bool
		deleted   bool
	}
)

func newObject(db *CommitStateDB, accProto authexported.Account) *stateObject {
	acc, ok := accProto.(*types.BaseAccount)
	if !ok {
		panic(fmt.Sprintf("invalid account type for state object: %T", accProto))
	}

	return &stateObject{
		stateDB:       db,
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
	// if the new value is the same as old, don't set
	prev := so.GetState(key)
	if prev == value {
		return
	}

	prefixKey := so.GetStorageByAddressKey(key.Bytes())

	// since the new value is different, update and journal the change
	so.stateDB.journal.append(storageChange{
		account:   &so.address,
		key:       prefixKey,
		prevValue: prev,
	})

	so.setState(prefixKey, value)
}

func (so *stateObject) setState(key, value sdk.Hash) {
	so.dirtyStorage[key] = value
}

// SetCode sets the state object's code.
func (so *stateObject) SetCode(codeHash sdk.Hash, code []byte) {
	prevCode := so.Code()

	so.stateDB.journal.append(codeChange{
		account:  &so.address,
		prevHash: so.CodeHash(),
		prevCode: prevCode,
	})

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

	// EIP158: We must check emptiness for the objects such that the account
	// clearing (0,0,0 objects) can take effect.
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

	so.stateDB.journal.append(balanceChange{
		account: &so.address,
		prev:    so.account.Balance(),
	})

	so.setBalance(amt)
}

func (so *stateObject) setBalance(amount sdk.Int) {
	so.account.SetBalance(amount)
}

// SetNonce sets the state object's nonce (sequence number).
func (so *stateObject) SetNonce(nonce uint64) {
	so.stateDB.journal.append(nonceChange{
		account: &so.address,
		prev:    so.account.Sequence,
	})

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

// commitState commits all dirty storage to a KVStore.

type DebugAccKV struct {
	AccAddr sdk.AccAddress `json:"acc_addr"`
	K       sdk.Hash       `json:"k"`
	V       sdk.Hash       `json:"v"`
}

func (dkv *DebugAccKV) String() string {
	return fmt.Sprintf(`{"K": "%s_%s", "V": "%s"}`, dkv.AccAddr.String(), dkv.K.String(), dkv.V.String())
}

func (dkv *DebugAccKV) MarshalJSON() ([]byte, error) {
	return ([]byte)(dkv.String()), nil
}

var DebugKeyPrefix = ([]byte)("DEBUG:")

func (dkv *DebugAccKV) Reset(accAddr sdk.AccAddress, k, v sdk.Hash) *DebugAccKV {
	dkv.AccAddr = accAddr
	dkv.K = k
	dkv.V = v

	return dkv
}

func (dkv *DebugAccKV) DebugAccKVFromKV(k, v []byte) {
	addrByte := k[len(DebugKeyPrefix) : len(DebugKeyPrefix)+20]
	dkv.AccAddr = addrByte
	dkv.K = sdk.BytesToHash(k[len(DebugKeyPrefix)+20:])
	dkv.V = sdk.BytesToHash(v)
}

func (dkv *DebugAccKV) DebugAccKVToKV() (k []byte, v sdk.Hash) {
	k = append(k, DebugKeyPrefix...)
	k = append(k, dkv.AccAddr...)
	k = append(k, dkv.K.Bytes()...)
	v = dkv.V
	return
}

func (so *stateObject) commitState() {
	ctx := so.stateDB.ctx
	store := ctx.KVStore(so.stateDB.storageKey)

	debugStore := (sdk.KVStore)(nil)
	if so.stateDB.debug {
		debugStore = ctx.KVStore(so.stateDB.storageDebugKey)
	}

	var kv DebugAccKV
	for key, value := range so.dirtyStorage {
		delete(so.dirtyStorage, key)

		if value == so.originStorage[key] {
			continue
		}

		so.originStorage[key] = value

		if (value == sdk.Hash{}) {
			store.Delete(key.Bytes())

			if debugStore != nil {
				k, _ := kv.Reset(so.address, key, value).DebugAccKVToKV()
				debugStore.Delete(k)
			}

			continue
		}

		store.Set(key.Bytes(), value.Bytes())

		if debugStore != nil {
			k, v := kv.Reset(so.address, key, value).DebugAccKVToKV()
			debugStore.Set(k, v.Bytes())
		}
	}

}

// commitCode persists the state object's code to the KVStore.
func (so *stateObject) commitCode() {
	ctx := so.stateDB.ctx
	store := ctx.KVStore(so.stateDB.codeKey)
	store.Set(so.CodeHash(), so.code)
}

// ----------------------------------------------------------------------------
// Getters
// ----------------------------------------------------------------------------

// Address returns the address of the state object.
func (so stateObject) Address() sdk.Address {
	return so.address
}

// Balance returns the state object's current balance.
func (so *stateObject) Balance() *big.Int {
	return so.account.Balance().BigInt()
}

// CodeHash returns the state object's code hash.
func (so *stateObject) CodeHash() []byte {
	return so.account.CodeHash
}

// Nonce returns the state object's current nonce (sequence number).
func (so *stateObject) Nonce() uint64 {
	return so.account.Sequence
}

// Code returns the contract code associated with this object, if any.
func (so *stateObject) Code() []byte {
	if so.code != nil {
		return so.code
	}

	if so.CodeHash() == nil {
		return nil
	}

	ctx := so.stateDB.ctx
	store := ctx.KVStore(so.stateDB.codeKey)
	code := store.Get(so.CodeHash())

	if len(code) == 0 {
		so.setError(fmt.Errorf("failed to get code hash %x for address: %x", so.CodeHash(), so.Address()))
	}

	so.code = code
	return code
}

// GetState retrieves a value from the account storage trie. Note, the key will
// be prefixed with the address of the state object.
func (so *stateObject) GetState(key sdk.Hash) sdk.Hash {
	prefixKey := so.GetStorageByAddressKey(key.Bytes())

	// if we have a dirty value for this state entry, return it
	value, dirty := so.dirtyStorage[prefixKey]
	if dirty {
		return value
	}

	// otherwise return the entry's original value
	return so.GetCommittedState(key)
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

	// otherwise load the value from the KVStore
	ctx := so.stateDB.ctx
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

// ReturnGas returns the gas back to the origin. Used by the Virtual machine or
// Closures. It performs a no-op.
func (so *stateObject) ReturnGas(gas *big.Int) {}

func (so *stateObject) deepCopy(db *CommitStateDB) *stateObject {
	newStateObj := newObject(db, so.account)

	newStateObj.code = so.code
	newStateObj.dirtyStorage = so.dirtyStorage.Copy()
	newStateObj.originStorage = so.originStorage.Copy()
	newStateObj.suicided = so.suicided
	newStateObj.dirtyCode = so.dirtyCode
	newStateObj.deleted = so.deleted

	return newStateObj
}

// empty returns whether the account is considered empty.
func (so *stateObject) empty() bool {
	return so.account.Sequence == 0 &&
		so.account.Balance().Sign() == 0 &&
		len(so.account.CodeHash) == 0
}

func (so *stateObject) touch() {
	so.stateDB.journal.append(touchChange{
		account: &so.address,
	})

	//if so.address == ripemd {//TODO check
	//	// Explicitly put it in the dirty-cache, which is otherwise generated from
	//	// flattened journals.
	//	so.stateDB.journal.dirty(so.address)
	//}
}

// GetStorageByAddressKey returns a hash of the composite key for a state
// object's storage prefixed with it's address.
func (so stateObject) GetStorageByAddressKey(key []byte) sdk.Hash {
	prefix := so.Address().Bytes()
	compositeKey := make([]byte, len(prefix)+len(key))

	copy(compositeKey, prefix)
	copy(compositeKey[len(prefix):], key)

	return sdk.BytesToHash(crypto.Sha256(compositeKey))
}

type SO struct {
	Address       sdk.AccAddress    `json:"address" yaml:"address"`
	BaseAccount   types.BaseAccount `json:"base_account" yaml:"base_account"`
	OriginStorage sdk.Storage       `json:"origin_storage" yaml:"origin_storage"`
	DirtyStorage  sdk.Storage       `json:"dirty_storage" yaml:"dirty_storage"`
	DirtyCode     bool              `json:"dirty_code" yaml:"dirty_code"`
	Suicided      bool              `json:"suicided" yaml:"suicided"`
	Deleted       bool              `json:"deleted" yaml:"deleted"`
	Code          sdk.Code          `json:"code" yaml:"code"`
}

type SOs []SO
