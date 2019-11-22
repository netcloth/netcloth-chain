package types

import (
	"bytes"
	"fmt"
	"math/big"

	authexported "github.com/netcloth/netcloth-chain/modules/auth/exported"
	"github.com/netcloth/netcloth-chain/modules/auth/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

var (
	emptyCodeHash = sdk.Hash{}
)

type (
	// StateObject interface for interacting with state object
	StateObject interface {
		GetCommittedState(key sdk.Hash) sdk.Hash
		GetState(key sdk.Hash)
		SetState(key, value sdk.Hash)

		Code() []byte
		SetCode(codeHash sdk.Hash, code []byte)
		CodeHash() []byte // codeHash = crypto.Sha256(Code)

		AddBalance(amount *big.Int)
		SubBalance(amount *big.Int)
		SetBalance(amount *big.Int)

		Balance() *big.Int
		ReturnGas(gas *big.Int)
		Address() sdk.AccAddress
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

func newObject(accProto authexported.Account) *stateObject {
	acc, ok := accProto.(*types.BaseAccount)
	if !ok {
		panic(fmt.Sprintf("invalid account type for state object: %T", accProto))
	}

	if acc.CodeHash == nil {
		acc.CodeHash = emptyCodeHash
	}

	return &stateObject{
		account:       acc,
		address:       acc.Address,
		originStorage: make(sdk.Storage),
		dirtyStorage:  make(sdk.Storage),
	}
}

// ----------------------------------------------------------------------------
// Setters
// ----------------------------------------------------------------------------

// SetCode
func (so *stateObject) SetCode(codeHash sdk.Hash, code []byte) {
	prevCode := so.Code(nil)
}

// CodeHash returns the state object's code hash
func (so *stateObject) CodeHash() sdk.Hash {
	return so.account.CodeHash
}

// setError remembers the first non-nil error it is called with.
func (so *stateObject) setError(err error) {
	if so.dbErr == nil {
		so.dbErr = err
	}
}

// Code returns the contract code associated with this object
func (so *stateObject) Code() []byte {
	if so.code != nil {
		return so.code
	}

	if bytes.Equal(so.CodeHash(), emptyCodeHash) {
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

// ----------------------------------------------------------------------------
// Getters
// ----------------------------------------------------------------------------
