package types

import (
	"fmt"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

const (
	CodeInvalidOperator       = 100
	CodeProfilerExists        = 101
	CodeProfilerNotExists     = 102
	CodeInvalidDescription    = 105
	CodeDeleteGenesisProfiler = 106
	CodeInvalidGuardian       = 108
	CodeAddressEmpty          = 120
	CodeAddedByEmpty          = 121
	CodeDeletedByEmpty        = 122
)

func ErrInvalidOperator(operator sdk.AccAddress) error {
	return sdkerrors.New(ModuleName, CodeInvalidOperator, fmt.Sprintf("%s is not a valid operator", operator))
}

func ErrProfilerNotExists(profiler sdk.AccAddress) error {
	return sdkerrors.New(ModuleName, CodeProfilerNotExists, fmt.Sprintf("profiler %s is not existed", profiler))
}

func ErrDeleteGenesisProfiler(profiler sdk.AccAddress) error {
	return sdkerrors.New(ModuleName, CodeDeleteGenesisProfiler, fmt.Sprintf("can't delete profiler %s that in genesis", profiler))
}

func ErrProfilerExists(profiler sdk.AccAddress) error {
	return sdkerrors.New(ModuleName, CodeProfilerExists, fmt.Sprintf("profiler %s already exists", profiler))
}

func ErrInvalidDescription() error {
	return sdkerrors.New(ModuleName, CodeInvalidDescription, fmt.Sprintf("description is invalid, length should be in range 1 to %d", MaxDescLenght))
}

func ErrAddressEmpty() error {
	return sdkerrors.New(ModuleName, CodeAddressEmpty, "address is empty")
}

func ErrAddedByEmpty() error {
	return sdkerrors.New(ModuleName, CodeAddedByEmpty, "added_by is empty")
}

func ErrDeletedByEmpty() error {
	return sdkerrors.New(ModuleName, CodeDeletedByEmpty, "deleted_by is empty")
}
