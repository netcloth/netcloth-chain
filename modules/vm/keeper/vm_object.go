package keeper

import (
	"bytes"
	"fmt"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/netcloth/netcloth-chain/types"
)

var _ StateObject = (*stateObject)(nil)
var emptyCodeHash = sdk.Hash{}

type (
	StateObject interface { // TODO: rename StateObject to VMObject
		Address() sdk.AccAddress

		GetCommittedState(key sdk.Hash) sdk.Hash
		GetState(key sdk.Hash) sdk.Hash
		SetState(key, value sdk.Hash)

		Code() []byte
		SetCode(codeHash sdk.Hash, code []byte)
		CodeHash() []byte
	}

	stateObject struct { // TODO: rename stateObject to vmObject
		sCtx StorageContext

		address sdk.AccAddress

		code     sdk.Code
		codeHash sdk.Hash

		originStorage sdk.Storage
		dirtyStorage  sdk.Storage

		dbErr error

		dirtyCode bool
		suicided  bool
		deleted   bool
	}
)

func newObject(addr sdk.AccAddress, code []byte, codeHash sdk.Hash) *stateObject {
	return &stateObject{
		address:       addr,
		code:          code, //TODO check: code = append(code, code)?
		codeHash:      codeHash,
		originStorage: make(sdk.Storage),
		dirtyStorage:  make(sdk.Storage),
	}
}

func (so *stateObject) WithContext(ctx StorageContext) *stateObject {
	so.sCtx = ctx
	return so
}

func (so *stateObject) GetStore() sdk.KVStore {
	return so.sCtx.ctx.KVStore(*so.sCtx.storeKey)
}

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

func (so *stateObject) SetCode(codeHash sdk.Hash, code []byte) {
	so.setCode(codeHash, code)
}

func (so *stateObject) setCode(codeHash sdk.Hash, code []byte) {
	so.code = code
	so.codeHash = codeHash
	so.dirtyCode = true
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
	if so.dirtyCode {
		store := so.GetStore()
		store.Set(so.CodeHash(), so.code)
	}
}

func (so *stateObject) commitState() {
	store := so.GetStore()
	for key, value := range so.dirtyStorage {
		delete(so.dirtyStorage, key)

		if value == so.originStorage[key] {
			continue
		}

		so.originStorage[key] = value

		if (value == sdk.Hash{}) {
			store.Delete(key.Bytes())
			continue
		}

		store.Set(key.Bytes(), value.Bytes())
	}
}

// ----------------------------------------------------------------------------
// Getters
// ----------------------------------------------------------------------------

func (so stateObject) Address() sdk.AccAddress {
	return so.address
}

func (so *stateObject) CodeHash() []byte {
	return so.codeHash.Bytes()
}

func (so *stateObject) Code() []byte {
	if so.code != nil {
		return so.code
	}

	if bytes.Equal(so.CodeHash(), emptyCodeHash.Bytes()) {
		return nil
	}

	store := so.GetStore()
	code := store.Get(so.CodeHash())

	if len(code) == 0 {
		so.setError(fmt.Errorf("failed to get code hash %x for address: %x", so.CodeHash(), so.address))
	}

	so.code = code
	return code
}

func (so *stateObject) GetState(key sdk.Hash) sdk.Hash {
	prefixKey := so.GetStorageByAddressKey(key.Bytes())

	value, dirty := so.dirtyStorage[prefixKey]
	if dirty {
		return value
	}

	return so.GetCommittedState(prefixKey)
}

func (so *stateObject) GetCommittedState(key sdk.Hash) sdk.Hash {
	prefixKey := so.GetStorageByAddressKey(key.Bytes())

	value, cached := so.originStorage[prefixKey]
	if cached {
		return value
	}

	store := so.GetStore()
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
	return bytes.Equal(so.CodeHash(), sdk.Hash{}.Bytes())
}

func (so stateObject) GetStorageByAddressKey(key []byte) sdk.Hash {
	prefix := so.Address().Bytes()
	compositeKey := make([]byte, len(prefix)+len(key))

	copy(compositeKey, prefix)
	copy(compositeKey[len(prefix):], key)

	return sdk.BytesToHash(crypto.Sha256(compositeKey))
}
