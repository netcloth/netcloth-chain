package types

import (
	"encoding/binary"
	sdk "github.com/NetCloth/netcloth-chain/types"
	"time"
)

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "ipal"

	// StoreKey is the store key string for ipal
	StoreKey = ModuleName

	// RouterKey is the message route for ipal
	RouterKey = ModuleName

	// QuerierRoute is the querier route for ipal
	QuerierRoute = ModuleName
)

var (
	IPALObjectKey                    = []byte{0x11} // prefix for each key to an ipal object

	ServerNodeObjectKey              = []byte{0x21} // prefix for each key to a ServerNode object
	ServerNodeObjectByStakeSharesKey = []byte{0x22}

	UnStakingKey                     = []byte{0x31}

	UnStakingTODOKey                 = []byte{0x41}
)

func GetIPALObjectKey(addr string) []byte {
	return append(IPALObjectKey, []byte(addr)...)
}

func GetServerNodeObjectKey(addr sdk.AccAddress) []byte {
	return append(ServerNodeObjectKey, addr...)
}

func GetServerNodeObjectByStakeSharesKey(obj ServerNodeObject) []byte {
	stakeShares := obj.StakeShares.Amount.Int64()
	stakeSharesBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(stakeSharesBytes[:], uint64(stakeShares))

	stakeSharesBytesLen := len(stakeSharesBytes)

	key := make([]byte, 1+stakeSharesBytesLen+sdk.AddrLen)

	key[0] = ServerNodeObjectByStakeSharesKey[0]
	copy(key[1:1+stakeSharesBytesLen], stakeSharesBytes)
	addr := sdk.CopyBytes(obj.OperatorAddress)
	for i, b := range addr {
		addr[i] = ^b
	}

	copy(key[1+stakeSharesBytesLen:], addr)
	return key
}

func GetUnStakingKey(accountAddress sdk.AccAddress) []byte {
	return append(UnStakingKey, accountAddress.Bytes()...)
}

func GetUnstakingTimeKey(timestamp time.Time) []byte {
	v := sdk.FormatTimeBytes(timestamp)
	return append(UnStakingTODOKey, v...)
}