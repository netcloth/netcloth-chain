package types

import (
	"github.com/netcloth/netcloth-chain/app/protocol"
)

const (
	ModuleName    = protocol.VMModuleName
	StoreKey      = ModuleName
	CodeKey       = StoreKey + "_code"
	LogKey        = StoreKey + "_log"
	StoreDebugKey = StoreKey + "_debug"
	QuerierRoute  = ModuleName
	RouterKey     = ModuleName
)

var (
	LogIndexKey = []byte("logIndexKey")
)
