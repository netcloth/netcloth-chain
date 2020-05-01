// nolint
package types

import (
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

var (
	ErrEmptyValidatorAddr              = sdkerrors.New(ModuleName, 1, "empty validator address")
	ErrBadValidatorAddr                = sdkerrors.New(ModuleName, 2, "validator address is invalid")
	ErrNoValidatorFound                = sdkerrors.New(ModuleName, 3, "validator does not exist")
	ErrValidatorOwnerExists            = sdkerrors.New(ModuleName, 4, "validator already exist for this operator address; must use new validator operator address")
	ErrValidatorPubKeyExists           = sdkerrors.New(ModuleName, 5, "validator already exist for this pubkey; must use new validator pubkey")
	ErrValidatorPubKeyTypeNotSupported = sdkerrors.New(ModuleName, 6, "validator pubkey type is not supported")
	ErrValidatorJailed                 = sdkerrors.New(ModuleName, 7, "validator for this address is currently jailed")
	ErrBadRemoveValidator              = sdkerrors.New(ModuleName, 8, "failed to remove validator")
	ErrCommissionNegative              = sdkerrors.New(ModuleName, 9, "commission must be positive")
	ErrCommissionHuge                  = sdkerrors.New(ModuleName, 10, "commission cannot be more than 100%")
	ErrCommissionGTMaxRate             = sdkerrors.New(ModuleName, 11, "commission cannot be more than the max rate")
	ErrCommissionUpdateTime            = sdkerrors.New(ModuleName, 12, "commission cannot be changed more than once in 24h")
	ErrCommissionChangeRateNegative    = sdkerrors.New(ModuleName, 13, "commission change rate must be positive")
	ErrCommissionChangeRateGTMaxRate   = sdkerrors.New(ModuleName, 14, "commission change rate cannot be more than the max rate")
	ErrCommissionGTMaxChangeRate       = sdkerrors.New(ModuleName, 15, "commission cannot be changed more than max change rate")
	ErrSelfDelegationBelowMinimum      = sdkerrors.New(ModuleName, 16, "validator's self delegation must be greater than their minimum self delegation")
	ErrMinSelfDelegationInvalid        = sdkerrors.New(ModuleName, 17, "minimum self delegation must be a positive integer")
	ErrMinSelfDelegationDecreased      = sdkerrors.New(ModuleName, 18, "minimum self delegation cannot be decrease")
	ErrEmptyDelegatorAddr              = sdkerrors.New(ModuleName, 19, "empty delegator address")
	ErrBadDenom                        = sdkerrors.New(ModuleName, 20, "invalid coin denomination")
	ErrBadDelegationAddr               = sdkerrors.New(ModuleName, 21, "invalid address for (address, validator) tuple")
	ErrBadDelegationAmount             = sdkerrors.New(ModuleName, 22, "invalid delegation amount")
	ErrNoDelegation                    = sdkerrors.New(ModuleName, 23, "no delegation for (address, validator) tuple")
	ErrBadDelegatorAddr                = sdkerrors.New(ModuleName, 24, "delegator does not exist with address")
	ErrNoDelegatorForAddress           = sdkerrors.New(ModuleName, 25, "delegator does not contain delegation")
	ErrInsufficientShares              = sdkerrors.New(ModuleName, 26, "insufficient delegation shares")
	ErrDelegationValidatorEmpty        = sdkerrors.New(ModuleName, 27, "cannot delegate to an empty validator")
	ErrNotEnoughDelegationShares       = sdkerrors.New(ModuleName, 28, "not enough delegation shares")
	ErrBadSharesAmount                 = sdkerrors.New(ModuleName, 29, "invalid shares amount")
	ErrBadSharesPercent                = sdkerrors.New(ModuleName, 30, "Invalid shares percent")
	ErrNotMature                       = sdkerrors.New(ModuleName, 31, "entry not mature")
	ErrNoUnbondingDelegation           = sdkerrors.New(ModuleName, 32, "no unbonding delegation found")
	ErrMaxUnbondingDelegationEntries   = sdkerrors.New(ModuleName, 33, "too many unbonding delegation entries for (delegator, validator) tuple")
	ErrBadRedelegationAddr             = sdkerrors.New(ModuleName, 34, "invalid address for (address, src-validator, dst-validator) tuple")
	ErrNoRedelegation                  = sdkerrors.New(ModuleName, 35, "no redelegation found")
	ErrSelfRedelegation                = sdkerrors.New(ModuleName, 36, "cannot redelegate to the same validator")
	ErrTinyRedelegationAmount          = sdkerrors.New(ModuleName, 37, "too few tokens to redelegate (truncates to zero tokens)")
	ErrBadRedelegationDst              = sdkerrors.New(ModuleName, 38, "redelegation destination validator not found")
	ErrTransitiveRedelegation          = sdkerrors.New(ModuleName, 39, "redelegation to this validator already in progress; first redelegation to this validator must complete before next redelegation")
	ErrMaxRedelegationEntries          = sdkerrors.New(ModuleName, 40, "too many redelegation entries for (delegator, src-validator, dst-validator) tuple")
	ErrDelegatorShareExRateInvalid     = sdkerrors.New(ModuleName, 41, "cannot delegate to validators with invalid (zero) ex-rate")
	ErrBothShareMsgsGiven              = sdkerrors.New(ModuleName, 42, "both shares amount and shares percent provided")
	ErrNeitherShareMsgsGiven           = sdkerrors.New(ModuleName, 43, "neither shares amount nor shares percent provided")
	ErrInvalidHistoricalInfo           = sdkerrors.New(ModuleName, 44, "invalid historical info")
	ErrNoHistoricalInfo                = sdkerrors.New(ModuleName, 45, "no historical info found")
	ErrDelegatorShareExceedMaxLever    = sdkerrors.New(ModuleName, 46, "delegation exceed max lever")
)
