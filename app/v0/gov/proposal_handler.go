package gov

import (
	"github.com/netcloth/netcloth-chain/app/v0/gov/types"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

// NewGovProposalHandler returns an sdk.Handler for gov proposal
func NewGovProposalHandler(k Keeper) Handler {
	return func(ctx sdk.Context, content Content, pid uint64, proposer sdk.AccAddress) error {
		switch {
		case ProposalTypeText == content.ProposalType():
			return nil

		case ProposalTypeSoftwareUpgrade == content.ProposalType():
			c, ok := content.(SoftwareUpgradeProposal)
			if !ok {
				sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "proposal type must be SoftwareUpgradeProposal")
			}
			return handleSoftwareUpgradeProposal(ctx, k, c, pid, proposer)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized gov proposal type: %s", content.ProposalType())
		}
	}
}

func handleSoftwareUpgradeProposal(ctx sdk.Context, keeper Keeper, proposalContent SoftwareUpgradeProposal, pid uint64, proposer sdk.AccAddress) error {
	if err := proposalContent.ValidateBasic(); err != nil {
		return err
	}

	if !keeper.pk.IsValidVersion(ctx, proposalContent.Version) {
		return types.ErrSoftwareUpgradeInvalidVersion
	}

	if uint64(ctx.BlockHeight()) > proposalContent.SwitchHeight {
		return types.ErrSoftwareUpgradeInvalidSwitchHeight
	}

	_, found := keeper.gk.GetProfiler(ctx, proposer)
	if !found {
		return types.ErrSoftwareUpgradeInvalidProfiler
	}

	if _, ok := keeper.pk.GetUpgradeConfig(ctx); ok {
		return types.ErrSoftwareUpgradeSwitchPeriodInProcess
	}

	pd := sdk.NewProtocolDefinition(proposalContent.Version, proposalContent.Software, proposalContent.SwitchHeight, proposalContent.Threshold)
	uc := sdk.NewUpgradeConfig(pid, pd)
	keeper.pk.SetUpgradeConfig(ctx, uc)

	return nil
}
