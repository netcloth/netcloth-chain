package types

import (
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

var (
	ErrEmptyInputs                    = sdkerrors.New(ModuleName, 1, "empty input")
	ErrStringTooLong                  = sdkerrors.New(ModuleName, 2, "input string to long")
	ErrInvalidSignature               = sdkerrors.New(ModuleName, 3, "CIPAL invalid user_request signature")
	ErrIPALClaimUserRequestExpired    = sdkerrors.New(ModuleName, 4, "CIPAL user_request time expired")
	ErrCIPALClaimUserRequestSigVerify = sdkerrors.New(ModuleName, 5, "CIPAL user_request signature verify failed")
)
