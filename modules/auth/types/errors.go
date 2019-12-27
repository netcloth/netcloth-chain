package types

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName
	CodeInvalidGas   sdk.CodeType      = 110
)

func ErrInvalidGas(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidGas, msg)
}
