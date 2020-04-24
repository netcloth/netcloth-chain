package tests

import (
	"bytes"
	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"math/rand"
	"sort"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/netcloth/netcloth-chain/types"
)

// Account contains a privkey, pubkey, address tuple
// eventually more useful data can be placed in here.
// (e.g. number of coins)
type Account struct {
	PrivKey crypto.PrivKey
	PubKey  crypto.PubKey
	Address sdk.AccAddress
}

// are two accounts equal
func (acc Account) Equals(acc2 Account) bool {
	return acc.Address.Equals(acc2.Address)
}

// RandomAcc pick a random account from an array
func RandomAcc(r *rand.Rand, accs []Account) Account {
	return accs[r.Intn(
		len(accs),
	)]
}

// RandomAccounts generates n random accounts
func RandomAccounts(r *rand.Rand, n int) []Account {
	accs := make([]Account, n)
	for i := 0; i < n; i++ {
		// don't need that much entropy for simulation
		privkeySeed := make([]byte, 15)
		r.Read(privkeySeed)
		useSecp := r.Int63()%2 == 0
		if useSecp {
			accs[i].PrivKey = secp256k1.GenPrivKeySecp256k1(privkeySeed)
		} else {
			accs[i].PrivKey = ed25519.GenPrivKeyFromSecret(privkeySeed)
		}

		accs[i].PubKey = accs[i].PrivKey.PubKey()
		accs[i].Address = sdk.AccAddress(accs[i].PubKey.Address())
	}

	return accs
}

// Type that combines an Address with the privKey and pubKey to that address
type AddrKeys struct {
	Address sdk.AccAddress
	PubKey  crypto.PubKey
	PrivKey crypto.PrivKey
}

func NewAddrKeys(address sdk.AccAddress, pubKey crypto.PubKey,
	privKey crypto.PrivKey) AddrKeys {

	return AddrKeys{
		Address: address,
		PubKey:  pubKey,
		PrivKey: privKey,
	}
}

// implement `Interface` in sort package.
type AddrKeysSlice []AddrKeys

func (b AddrKeysSlice) Len() int {
	return len(b)
}

// Sorts lexographically by Address
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

func (b AddrKeysSlice) Swap(i, j int) {
	b[j], b[i] = b[i], b[j]
}

// CreateGenAccounts generates genesis accounts loaded with coins, and returns
// their addresses, pubkeys, and privkeys.
func CreateGenAccounts(numAccs int, genCoins sdk.Coins) (genAccs []auth.Account,
	addrs []sdk.AccAddress, pubKeys []crypto.PubKey, privKeys []crypto.PrivKey) {

	addrKeysSlice := AddrKeysSlice{}

	for i := 0; i < numAccs; i++ {
		privKey := secp256k1.GenPrivKey()
		pubKey := privKey.PubKey()
		addr := sdk.AccAddress(pubKey.Address())

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
