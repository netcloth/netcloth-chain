package types

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeOutOfGas                sdk.CodeType = 101
	CodeStoreOutOfGas           sdk.CodeType = 102
	CodeDepth                   sdk.CodeType = 103
	CodeTraceLimitReached       sdk.CodeType = 104
	CodeNoCompatibleInterpreter sdk.CodeType = 105
	CodeEmptyInputs             sdk.CodeType = 106
	CodeInsufficientBalance     sdk.CodeType = 107
	CodeContractExist           sdk.CodeType = 108
	CodeNoCodeExist             sdk.CodeType = 109
	CodeMaxCodeSizeExceeded     sdk.CodeType = 110
	CodeWriteProtection         sdk.CodeType = 111
	CodeReturnDataOutOfBounds   sdk.CodeType = 112
	CodeExecutionReverted       sdk.CodeType = 113
	CodeInvalidJump             sdk.CodeType = 114
	CodeGasUintOverflow         sdk.CodeType = 115
)

func ErrOutOfGas() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeOutOfGas, "out of gas")
}

func ErrCodeStoreOutOfGas() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeStoreOutOfGas, "contract creation code storage out of gas")
}

func ErrDepth() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeDepth, "max call depth exceeded")
}

func ErrTraceLimitReached() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeTraceLimitReached, "the number of logs reached the specified limit")
}

func ErrNoCompatibleInterpreter() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeNoCompatibleInterpreter, "no compatible interpreter")
}

func ErrEmptyInputs() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeEmptyInputs, "empty input")
}

func ErrInsufficientBalance() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInsufficientBalance, "insufficient balance for transfer")
}

func ErrContractAddressCollision() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeContractExist, "contract address collision")
}

func ErrNoCodeExist() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeNoCodeExist, "code exists")
}

func ErrMaxCodeSizeExceeded() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeMaxCodeSizeExceeded, "evm: max code size exceeded")
}

func ErrWriteProtection() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeWriteProtection, "vm: write protection")
}

func ErrReturnDataOutOfBounds() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeReturnDataOutOfBounds, "evm: return data out of bounds")
}

func ErrExecutionReverted() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeExecutionReverted, "evm: execution reverted")
}

func ErrInvalidJump() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidJump, "evm: invalid jump destination")
}

func ErrGasUintOverflow() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeGasUintOverflow, "gas uint64 overflow")
}
