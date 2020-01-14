package types

import (
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

var (
	ErrEmptyInputs                    = sdkerrors.Register(ModuleName, 1, "empty input")
	ErrStringTooLong                  = sdkerrors.Register(ModuleName, 2, "input string to long")
	ErrInvalidSignature               = sdkerrors.Register(ModuleName, 3, "CIPAL invalid user_request signature")
	ErrIPALClaimUserRequestExpired    = sdkerrors.Register(ModuleName, 4, "CIPAL user_request time expired")
	ErrCIPALClaimUserRequestSigVerify = sdkerrors.Register(ModuleName, 5, "CIPAL user_request signature verify failed")
)
