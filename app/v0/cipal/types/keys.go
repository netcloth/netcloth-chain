package types

import (
	"github.com/netcloth/netcloth-chain/baseapp/protocol"
)

const (
	ModuleName   = protocol.CIpalModuleName
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

var (
	CIPALObjectKey = []byte{0x11}
)

func GetCIPALObjectKey(addr string) []byte {
	return append(CIPALObjectKey, []byte(addr)...)
}
