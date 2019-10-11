package types

import (
	sdk "github.com/NetCloth/netcloth-chain/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeEmptyInputs sdk.CodeType = 110
	CodeStringTooLong sdk.CodeType = 111
)


// ErrEmptyInputs is an error
func ErrEmptyInputs(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyInputs, "empty input to ipal transaction")
}

// ErrStringTooLon is an error
func ErrStringTooLong(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeStringTooLong, "string size exceeds limit")
}