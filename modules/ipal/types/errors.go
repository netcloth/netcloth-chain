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
	CodeBadDenom                       sdk.CodeType = 114
	CodeStakeSharesInsufficient        sdk.CodeType = 115
)

// ErrEmptyInputs is an error
func ErrEmptyInputs(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeEmptyInputs, msg)
}

// ErrStringTooLon is an error
func ErrStringTooLong(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeStringTooLong, msg)
}

// ErrInvalidSignature is an error
func ErrInvalidSignature(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidIPALClaimUserRequestSig, msg)
}

// ErrIPALClaimUserRequestExpired is an error
func ErrIPALClaimUserRequestExpired(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeIPALClaimUserRequestExpired, msg)
}

func ErrBadDenom(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeBadDenom, msg)
}

func ErrStakeSharesInsufficient(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeStakeSharesInsufficient, msg)
}