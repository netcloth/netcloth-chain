package common

import (
	"encoding/binary"
	//"github.com/ethereum/go-ethereum/rlp"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/netcloth/netcloth-chain/types"
)

func CreateAddress(addr sdk.Address, nonce uint64) sdk.Address {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, nonce)
	data := append(addr.Bytes(), b...)
	//data, _ := rlp.EncodeToBytes([]interface{}{b, nonce})
	return sdk.AccAddress(crypto.Sha256(data)[12:])
}
