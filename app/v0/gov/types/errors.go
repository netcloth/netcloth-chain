//nolint
package types

import (
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

var (
	ErrUnknownProposal                      = sdkerrors.New(ModuleName, 1, "unknown proposal")
	ErrInactiveProposal                     = sdkerrors.New(ModuleName, 2, "inactive proposal")
	ErrAlreadyActiveProposal                = sdkerrors.New(ModuleName, 3, "proposal already active")
	ErrInvalidProposalContent               = sdkerrors.New(ModuleName, 4, "invalid proposal content")
	ErrInvalidProposalType                  = sdkerrors.New(ModuleName, 5, "invalid proposal type")
	ErrInvalidVote                          = sdkerrors.New(ModuleName, 6, "invalid vote option")
	ErrInvalidGenesis                       = sdkerrors.New(ModuleName, 7, "invalid genesis state")
	ErrNoProposalHandlerExists              = sdkerrors.New(ModuleName, 8, "no handler exists for proposal type")
	ErrSoftwareUpgradeProposalExist         = sdkerrors.New(ModuleName, 9, "software upgrade proposal exist")
	ErrSoftwareUpgradeInvalidVersion        = sdkerrors.New(ModuleName, 10, "invalid software upgrade version")
	ErrSoftwareUpgradeInvalidSwitchHeight   = sdkerrors.New(ModuleName, 11, "invalid software upgrade switch height")
	ErrSoftwareUpgradeInvalidProfiler       = sdkerrors.New(ModuleName, 12, "invalid software upgrade profiler")
	ErrSoftwareUpgradeSwitchPeriodInProcess = sdkerrors.New(ModuleName, 13, "software upgrade already in switch period")
	ErrSoftwareUpgradeInvalidThreshold      = sdkerrors.New(ModuleName, 14, "software upgrade Threshold should be in range [0.8, 1.0]")
)
