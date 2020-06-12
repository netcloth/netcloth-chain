package types

import (
	"encoding/binary"
	"time"

	"github.com/netcloth/netcloth-chain/app/protocol"
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	ModuleName   = protocol.IpalModuleName
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

var (
	IPALNodeKey          = []byte{0x10}
	IPALNodeByBondKey    = []byte{0x11}
	IPALNodeByMonikerKey = []byte{0x12}
	UnBondingKey         = []byte{0x13}
)

func GetIPALNodeKey(addr sdk.AccAddress) []byte {
	return append(IPALNodeKey, addr...)
}

func GetIPALNodeByBondKey(obj IPALNode) []byte {
	bond := obj.Bond.Amount.Int64()
	bondBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bondBytes[:], uint64(bond))

	bondBytesLen := len(bondBytes)

	key := make([]byte, 1+bondBytesLen+sdk.AddrLen)

	key[0] = IPALNodeByBondKey[0]
	copy(key[1:1+bondBytesLen], bondBytes)
	addr := sdk.CopyBytes(obj.OperatorAddress)
	for i, b := range addr {
		addr[i] = ^b
	}

	copy(key[1+bondBytesLen:], addr)
	return key
}

func GetIPALNodeByMonikerKey(moniker string) []byte {
	return append(IPALNodeByMonikerKey, []byte(moniker)...)
}

func GetUnBondingKey(timestamp time.Time) []byte {
	v := sdk.FormatTimeBytes(timestamp)
	return append(UnBondingKey, v...)
}
