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
	return sdk.BytesToAddress(crypto.Sha256(data)[12:])
}

func CreateAddress2(addr sdk.AccAddress, salt [32]byte, inithash []byte) sdk.AccAddress {
	var bs []byte
	bs = append(bs, []byte{0xff}...)
	bs = append(bs, addr.Bytes()...)
	bs = append(bs, salt[:]...)
	bs = append(bs, inithash...)
	return sdk.BytesToAddress(crypto.Sha256(bs)[12:])
}
