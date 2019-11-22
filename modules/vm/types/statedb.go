package types

import (
	"crypto"
	"github.com/netcloth/netcloth-chain/modules/auth"
	sdk "github.com/netcloth/netcloth-chain/types"
	"sync"
)

type CommitStateDB struct {
	ctx sdk.Context

	ak         auth.AccountKeeper
	storageKey sdk.StoreKey
	codeKey    sdk.StoreKey

	//stateObjects      map[sdk.AccAddress]*stateObject
	//stateObjectsDirty map[sdk.AccAddress]struct{}

	refund uint64

	thash, bhash crypto.Hash
	txIndex      int
	// logs
	logSize   uint
	preimages map[crypto.Hash][]byte

	dbErr error

	lock sync.Mutex
}
