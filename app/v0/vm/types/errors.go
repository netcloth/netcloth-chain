package types

import (
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

var (
	ErrNoPayload                = sdkerrors.New(ModuleName, 1, "no payload")
	ErrOutOfGas                 = sdkerrors.New(ModuleName, 2, "out of gas")
	ErrCodeStoreOutOfGas        = sdkerrors.New(ModuleName, 3, "contract creation code storage out of gas")
	ErrDepth                    = sdkerrors.New(ModuleName, 4, "max call depth exceeded")
	ErrTraceLimitReached        = sdkerrors.New(ModuleName, 5, "the number of logs reached the specified limit")
	ErrNoCompatibleInterpreter  = sdkerrors.New(ModuleName, 6, "no compatible interpreter")
	ErrEmptyInputs              = sdkerrors.New(ModuleName, 7, "empty input")
	ErrInsufficientBalance      = sdkerrors.New(ModuleName, 8, "insufficient balance for transfer")
	ErrContractAddressCollision = sdkerrors.New(ModuleName, 9, "contract address collision")
	ErrNoCodeExist              = sdkerrors.New(ModuleName, 10, "no code exists")
	ErrMaxCodeSizeExceeded      = sdkerrors.New(ModuleName, 11, "evm: max code size exceeded")
	ErrWriteProtection          = sdkerrors.New(ModuleName, 12, "vm: write protection")
	ErrReturnDataOutOfBounds    = sdkerrors.New(ModuleName, 13, "evm: return data out of bounds")
	ErrExecutionReverted        = sdkerrors.New(ModuleName, 14, "evm: execution reverted")
	ErrInvalidJump              = sdkerrors.New(ModuleName, 15, "evm: invalid jump destination")
	ErrGasUintOverflow          = sdkerrors.New(ModuleName, 16, "gas uint64 overflow")
)
