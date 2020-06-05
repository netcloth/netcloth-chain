package types

import (
	"github.com/netcloth/netcloth-chain/baseapp/protocol"
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	ModuleName   = protocol.GuardianModuleName
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

var (
	profilerKey = []byte{0x00}
)

func GetProfilerKey(addr sdk.AccAddress) []byte {
	return append(profilerKey, addr.Bytes()...)
}

func GetProfilersSubspaceKey() []byte {
	return profilerKey
}
