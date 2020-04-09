// nolint
package types

import (
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

var (
	ErrEmptyDelegatorAddr      = sdkerrors.New(ModuleName, 1, "delegator address is empty")
	ErrEmptyWithdrawAddr       = sdkerrors.New(ModuleName, 2, "withdraw address is empty")
	ErrEmptyValidatorAddr      = sdkerrors.New(ModuleName, 3, "validator address is empty")
	ErrEmptyDelegationDistInfo = sdkerrors.New(ModuleName, 4, "no delegation distribution info")
	ErrNoValidatorDistInfo     = sdkerrors.New(ModuleName, 5, "no validator distribution info")
	ErrNoValidatorCommission   = sdkerrors.New(ModuleName, 6, "no validator commission to withdraw")
	ErrSetWithdrawAddrDisabled = sdkerrors.New(ModuleName, 7, "set withdraw address disabled")
	ErrBadDistribution         = sdkerrors.New(ModuleName, 8, "community pool does not have sufficient coins to distribute")
	ErrInvalidProposalAmount   = sdkerrors.New(ModuleName, 9, "invalid community pool spend proposal amount")
	ErrEmptyProposalRecipient  = sdkerrors.New(ModuleName, 10, "invalid community pool spend proposal recipient")
	ErrNoValidatorExists       = sdkerrors.New(ModuleName, 11, "validator does not exist")
	ErrNoDelegationExists      = sdkerrors.New(ModuleName, 12, "delegation does not exist")
)
