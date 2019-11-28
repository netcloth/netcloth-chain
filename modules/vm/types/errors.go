package types

import (
	"errors"

	sdk "github.com/netcloth/netcloth-chain/types"
)

var (
	ErrOutOfGas                 = errors.New("out of gas")
	ErrCodeStoreOutOfGas        = errors.New("contract creation code storage out of gas")
	ErrDepth                    = errors.New("max call depth exceeded")
	ErrTraceLimitReached        = errors.New("the number of logs reached the specified limit")
	ErrInsufficientBalance      = errors.New("insufficient balance for transfer")
	ErrContractAddressCollision = errors.New("contract address collision")
	ErrNoCompatibleInterpreter  = errors.New("no compatible interpreter")
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeEmptyInputs sdk.CodeType = 100
)

func ErrEmptyInputs(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeEmptyInputs, msg)
}
