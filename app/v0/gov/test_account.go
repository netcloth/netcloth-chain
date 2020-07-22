package gov

import (
	"bytes"
	"sort"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"github.com/netcloth/netcloth-chain/types"
)

// AddrKeys Type that combines an Address with the privKey and pubKey to that address
type AddrKeys struct {
	Address types.AccAddress
	PubKey  crypto.PubKey
	PrivKey crypto.PrivKey
}

// NewAddrKeys constructor of AddrKeys
func NewAddrKeys(address types.AccAddress, pubKey crypto.PubKey, privKey crypto.PrivKey) AddrKeys {
	return AddrKeys{
		Address: address,
		PubKey:  pubKey,
		PrivKey: privKey,
	}
}

// AddrKeysSlice implement `Interface` in sort package.
type AddrKeysSlice []AddrKeys

func (b AddrKeysSlice) Len() int {
	return len(b)
}

// Less Sorts lexographically by Address
func (b AddrKeysSlice) Less(i, j int) bool {
	// bytes package already implements Comparable for []byte.
	switch bytes.Compare(b[i].Address.Bytes(), b[j].Address.Bytes()) {
	case -1:
		return true
	case 0, 1:
		return false
	default:
		panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
	}
}

// Swap
func (b AddrKeysSlice) Swap(i, j int) {
	b[j], b[i] = b[i], b[j]
}

// CreateGenAccounts generates genesis accounts loaded with coins, and returns their addresses, pubkeys, and privkeys.
func CreateGenAccounts(numAccs int, genCoins types.Coins) (
	genAccs []auth.Account, addrs []types.AccAddress, pubKeys []crypto.PubKey, privKeys []crypto.PrivKey) {

	addrKeysSlice := AddrKeysSlice{}

	for i := 0; i < numAccs; i++ {
		privKey := secp256k1.GenPrivKey()
		pubKey := privKey.PubKey()
		addr := types.AccAddress(pubKey.Address())

		addrKeysSlice = append(addrKeysSlice, NewAddrKeys(addr, pubKey, privKey))
	}

	sort.Sort(addrKeysSlice)

	for i := range addrKeysSlice {
		addrs = append(addrs, addrKeysSlice[i].Address)
		pubKeys = append(pubKeys, addrKeysSlice[i].PubKey)
		privKeys = append(privKeys, addrKeysSlice[i].PrivKey)
		genAccs = append(genAccs, &auth.BaseAccount{
			Address: addrKeysSlice[i].Address,
			Coins:   genCoins,
		})
	}

	return
}
