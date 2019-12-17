package keeper

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// StateDB is an VM database for full state querying.
type StateDB interface {
	CreateAccount(sdk.AccAddress)

	SubBalance(sdk.AccAddress, *big.Int)
	AddBalance(sdk.AccAddress, *big.Int)
	GetBalance(sdk.AccAddress) *big.Int

	GetNonce(sdk.AccAddress) uint64
	SetNonce(sdk.AccAddress, uint64)

	GetCodeHash(sdk.AccAddress) sdk.Hash
	GetCode(sdk.AccAddress) []byte
	SetCode(sdk.AccAddress, []byte)
	GetCodeSize(sdk.AccAddress) int

	AddRefund(uint64)
	SubRefund(uint64)
	GetRefund() uint64

	GetCommittedState(sdk.AccAddress, sdk.Hash) sdk.Hash
	GetState(sdk.AccAddress, sdk.Hash) sdk.Hash
	SetState(sdk.AccAddress, sdk.Hash, sdk.Hash)

	Suicide(sdk.AccAddress) bool
	HasSuicided(sdk.AccAddress) bool

	// Exist reports whether the given account exists in state.
	// Notably this should also return true for suicided accounts.
	Exist(sdk.AccAddress) bool
	// Empty returns whether the given account is empty. Empty
	// is defined according to EIP161 (balance = nonce = code = 0).
	Empty(sdk.AccAddress) bool

	RevertToSnapshot(int)
	Snapshot() int

	AddLog(*types.Log)
	AddPreimage(sdk.Hash, []byte)

	ForEachStorage(sdk.AccAddress, func(sdk.Hash, sdk.Hash) bool) error
}
