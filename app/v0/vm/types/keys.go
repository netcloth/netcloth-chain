package types

import (
	"github.com/netcloth/netcloth-chain/app/protocol"
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	ModuleName   = protocol.VMModuleName
	StoreKey     = ModuleName
	QuerierRoute = ModuleName
	RouterKey    = ModuleName
)

var (
	LogIndexKey = []byte("logIndexKey")
)

// KVStore key prefixes
var (
	KeyPrefixLogs      = []byte{0x01}
	KeyPrefixLogsIndex = []byte{0x02}
	KeyPrefixCode      = []byte{0x03}
	KeyPrefixStorage   = []byte{0x04}
)

// AddressStoragePrefix returns a prefix to iterate over a given account storage.
func AddressStoragePrefix(address sdk.Address) []byte {
	return append(KeyPrefixStorage, address.Bytes()...)
}
