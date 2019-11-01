package types

import (
    sdk "github.com/NetCloth/netcloth-chain/types"
)

const (
    DefaultCodespace sdk.CodespaceType = ModuleName

    CodeEmptyInputs                    sdk.CodeType = 110
    CodeStringTooLong                  sdk.CodeType = 111
    CodeInvalidIPALClaimUserRequestSig sdk.CodeType = 112
    CodeIPALClaimUserRequestExpired    sdk.CodeType = 113
)

func ErrEmptyInputs(msg string) sdk.Error {
    return sdk.NewError(DefaultCodespace, CodeEmptyInputs, msg)
}

func ErrStringTooLong(msg string) sdk.Error {
    return sdk.NewError(DefaultCodespace, CodeStringTooLong, msg)
}

func ErrInvalidSignature(msg string) sdk.Error {
    return sdk.NewError(DefaultCodespace, CodeInvalidIPALClaimUserRequestSig, msg)
}

func ErrIPALClaimUserRequestExpired(msg string) sdk.Error {
    return sdk.NewError(DefaultCodespace, CodeIPALClaimUserRequestExpired, msg)
}
