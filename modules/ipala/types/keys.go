package types

import (
    "encoding/binary"
    sdk "github.com/NetCloth/netcloth-chain/types"
    "time"
)

const (
    ModuleName = "ipala"
    StoreKey = ModuleName
    RouterKey = ModuleName
    QuerierRoute = ModuleName
)

var (
    ServiceNodeKey       = []byte{0x10}
    ServiceNodeByBondKey = []byte{0x11}
    UnBondsKey           = []byte{0x12}
    UnBondingKey         = []byte{0x13}
)

func GetServiceNodeKey(addr sdk.AccAddress) []byte {
    return append(ServiceNodeKey, addr...)
}

func GetServiceNodeByBondKey(obj ServiceNode) []byte {
    bond := obj.Bond.Amount.Int64()
    bondBytes := make([]byte, 8)
    binary.BigEndian.PutUint64(bondBytes[:], uint64(bond))

    bondBytesLen := len(bondBytes)

    key := make([]byte, 1+bondBytesLen+sdk.AddrLen)

    key[0] = ServiceNodeByBondKey[0]
    copy(key[1:1+bondBytesLen], bondBytes)
    addr := sdk.CopyBytes(obj.OperatorAddress)
    for i, b := range addr {
        addr[i] = ^b
    }

    copy(key[1+bondBytesLen:], addr)
    return key
}

func GetUnBondsKey(accountAddress sdk.AccAddress) []byte {
    return append(UnBondsKey, accountAddress.Bytes()...)
}

func GetUnBondingKey(timestamp time.Time) []byte {
    v := sdk.FormatTimeBytes(timestamp)
    return append(UnBondingKey, v...)
}
