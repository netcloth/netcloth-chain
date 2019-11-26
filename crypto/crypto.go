package crypto

import (
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"

	sdk "github.com/netcloth/netcloth-chain/types"
)

func Keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

func CreateCommonAddress(b common.Address, nonce uint64) common.Address { // from ethereum
	data, _ := rlp.EncodeToBytes([]interface{}{b, nonce})
	return common.BytesToAddress(Keccak256(data)[12:])
}

func CreateAddress(b sdk.Address, nonce uint64) sdk.Address {
	data, _ := rlp.EncodeToBytes([]interface{}{b, nonce})
	return sdk.AccAddress(Keccak256(data)[12:])
}
