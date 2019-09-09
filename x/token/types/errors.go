//nolint
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "token"

	CodeInvalidMoniker       sdk.CodeType = 100
)

func ErrInvalidMoniker(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMoniker, msg)
}