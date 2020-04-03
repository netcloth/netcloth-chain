package types

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	ModuleName   = "guardian"
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
