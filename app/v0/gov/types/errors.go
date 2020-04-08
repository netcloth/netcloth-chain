//nolint
package types

import (
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

var (
	ErrUnknownProposal                      = sdkerrors.Register(ModuleName, 1, "unknown proposal")
	ErrInactiveProposal                     = sdkerrors.Register(ModuleName, 2, "inactive proposal")
	ErrAlreadyActiveProposal                = sdkerrors.Register(ModuleName, 3, "proposal already active")
	ErrInvalidProposalContent               = sdkerrors.Register(ModuleName, 4, "invalid proposal content")
	ErrInvalidProposalType                  = sdkerrors.Register(ModuleName, 5, "invalid proposal type")
	ErrInvalidVote                          = sdkerrors.Register(ModuleName, 6, "invalid vote option")
	ErrInvalidGenesis                       = sdkerrors.Register(ModuleName, 7, "invalid genesis state")
	ErrNoProposalHandlerExists              = sdkerrors.Register(ModuleName, 8, "no handler exists for proposal type")
	ErrSoftwareUpgradeProposalExist         = sdkerrors.Register(ModuleName, 9, "software upgrade proposal exist")
	ErrSoftwareUpgradeInvalidVersion        = sdkerrors.Register(ModuleName, 10, "invalid software upgrade version")
	ErrSoftwareUpgradeInvalidSwitchHeight   = sdkerrors.Register(ModuleName, 11, "invalid software upgrade switch height")
	ErrSoftwareUpgradeInvalidProfiler       = sdkerrors.Register(ModuleName, 12, "invalid software upgrade profiler")
	ErrSoftwareUpgradeSwitchPeriodInProcess = sdkerrors.Register(ModuleName, 13, "software upgrade already in switch period")
)
