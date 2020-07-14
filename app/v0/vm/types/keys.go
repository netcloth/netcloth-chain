package types

import (
	"github.com/netcloth/netcloth-chain/app/protocol"
)

const (
	ModuleName   = protocol.VMModuleName
	StoreKey     = ModuleName
	CodeKey      = StoreKey + "_code"
	LogKey       = StoreKey + "_log"
	QuerierRoute = ModuleName
	RouterKey    = ModuleName
)

var (
	LogIndexKey = []byte("logIndexKey")
)
