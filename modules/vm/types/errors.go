package types

import (
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

var (
	ErrNoPayload                = sdkerrors.Register(ModuleName, 1, "no payload")
	ErrOutOfGas                 = sdkerrors.Register(ModuleName, 2, "out of gas")
	ErrCodeStoreOutOfGas        = sdkerrors.Register(ModuleName, 3, "contract creation code storage out of gas")
	ErrDepth                    = sdkerrors.Register(ModuleName, 4, "max call depth exceeded")
	ErrTraceLimitReached        = sdkerrors.Register(ModuleName, 5, "the number of logs reached the specified limit")
	ErrNoCompatibleInterpreter  = sdkerrors.Register(ModuleName, 6, "no compatible interpreter")
	ErrEmptyInputs              = sdkerrors.Register(ModuleName, 7, "empty input")
	ErrInsufficientBalance      = sdkerrors.Register(ModuleName, 8, "insufficient balance for transfer")
	ErrContractAddressCollision = sdkerrors.Register(ModuleName, 9, "contract address collision")
	ErrNoCodeExist              = sdkerrors.Register(ModuleName, 10, "no code exists")
	ErrMaxCodeSizeExceeded      = sdkerrors.Register(ModuleName, 11, "evm: max code size exceeded")
	ErrWriteProtection          = sdkerrors.Register(ModuleName, 12, "vm: write protection")
	ErrReturnDataOutOfBounds    = sdkerrors.Register(ModuleName, 13, "evm: return data out of bounds")
	ErrExecutionReverted        = sdkerrors.Register(ModuleName, 14, "evm: execution reverted")
	ErrInvalidJump              = sdkerrors.Register(ModuleName, 15, "evm: invalid jump destination")
	ErrGasUintOverflow          = sdkerrors.Register(ModuleName, 16, "gas uint64 overflow")
)
