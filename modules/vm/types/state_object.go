package types

import (
	"math/big"

	"github.com/netcloth/netcloth-chain/modules/auth/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

type (
	StateObject interface {
		Code() []byte
		SetCode(codeHash sdk.Hash, code []byte)
		CodeHash() []byte // codeHash = crypto.Sha256(Code)

		AddBalance(amount *big.Int)
		SubBalance(amount *big.Int)
		SetBalance(amount *big.Int)

		Balance() *big.Int
		ReturnGas(gas *big.Int)
		Address() sdk.AccAddress
	}

	stateObject struct {
		address sdk.AccAddress
		stateDB *CommitStateDB
		account *types.BaseAccount

		dbErr error

		code sdk.Code
	}
)
