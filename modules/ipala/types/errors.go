package types

import (
    sdk "github.com/NetCloth/netcloth-chain/types"
)

const (
    DefaultCodespace sdk.CodespaceType = ModuleName

    CodeEmptyInputs      sdk.CodeType = 100
    CodeStringTooLong    sdk.CodeType = 101
    CodeBadDenom         sdk.CodeType = 102
    CodeBondInsufficient sdk.CodeType = 103
)

func ErrEmptyInputs(msg string) sdk.Error {
    return sdk.NewError(DefaultCodespace, CodeEmptyInputs, msg)
}

func ErrStringTooLong(msg string) sdk.Error {
    return sdk.NewError(DefaultCodespace, CodeStringTooLong, msg)
}

func ErrBadDenom(msg string) sdk.Error {
    return sdk.NewError(DefaultCodespace, CodeBadDenom, msg)
}

func ErrBondInsufficient(msg string) sdk.Error {
    return sdk.NewError(DefaultCodespace, CodeBondInsufficient, msg)
}
