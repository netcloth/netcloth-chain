package vm

import (
	sdk "github.com/netcloth/netcloth-chain/types"
	"math/big"
)

type StateDB interface {
	SubBalance(sdk.AccAddress, *big.Int)
	AddBalance(sdk.AccAddress, *big.Int)
	GetBalance(sdk.AccAddress) *big.Int

	AddPreimage(sdk.Hash, []byte)
}
