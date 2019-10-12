package types

import (
	sdk "github.com/NetCloth/netcloth-chain/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeEmptyInputs sdk.CodeType = 110
	CodeStringTooLong sdk.CodeType = 111
	CodeInvalidIPALClaimUserRequestSig = 112
	CdoeIPALClaimUserRequestExpired = 113
)


// ErrEmptyInputs is an error
func ErrEmptyInputs(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyInputs, "empty input to ipal transaction")
}

// ErrStringTooLon is an error
func ErrStringTooLong(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeStringTooLong, "string size exceeds limit")
}

// ErrInvalidSignature is an error
func ErrInvalidSignature(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidIPALClaimUserRequestSig, "invalid IPALClaim user request signature")
}

// ErrIPALClaimExpired is an error
func ErrIPALClaimExpired(codespaceType sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespaceType, CdoeIPALClaimUserRequestExpired, "IPALClaim user request expired")
}