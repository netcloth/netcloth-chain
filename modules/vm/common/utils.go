package common

import (
	"encoding/binary"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/netcloth/netcloth-chain/types"
)

func CreateAddress(addr sdk.AccAddress, nonce uint64) sdk.AccAddress {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, nonce)
	data := append(addr.Bytes(), b...)
	return BytesToAddress(crypto.Sha256(data)[12:])
}
