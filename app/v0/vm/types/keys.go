package types

import (
	"github.com/netcloth/netcloth-chain/baseapp/protocol"
)

const (
	ModuleName    = protocol.VmModuleName
	StoreKey      = ModuleName
	CodeKey       = StoreKey + "_code"
	StoreDebugKey = StoreKey + "_debug"
	QuerierRoute  = ModuleName
	RouterKey     = ModuleName
)
