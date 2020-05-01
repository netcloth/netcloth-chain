//nolint
package types

import (
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

// x/slashing module sentinel errors
var (
	ErrNoValidatorForAddress        = sdkerrors.New(ModuleName, 1, "address is not associated with any known validator")
	ErrBadValidatorAddr             = sdkerrors.New(ModuleName, 2, "validator does not exist for that address")
	ErrValidatorJailed              = sdkerrors.New(ModuleName, 3, "validator still jailed; cannot be unjailed")
	ErrValidatorNotJailed           = sdkerrors.New(ModuleName, 4, "validator not jailed; cannot be unjailed")
	ErrMissingSelfDelegation        = sdkerrors.New(ModuleName, 5, "validator has no self-delegation; cannot be unjailed")
	ErrSelfDelegationTooLowToUnjail = sdkerrors.New(ModuleName, 6, "validator's self delegation less than minimum; cannot be unjailed")
	ErrNoSigningInfoFound           = sdkerrors.New(ModuleName, 7, "no validator signing info found")
)
