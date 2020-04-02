package gov

import (
	"fmt"
	"strconv"

	"github.com/netcloth/netcloth-chain/app/v0/gov/types"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

// Handle all "gov" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgDeposit:
			return handleMsgDeposit(ctx, keeper, msg)

		case MsgSubmitProposal:
			return handleMsgSubmitProposal(ctx, keeper, msg)

		case MsgSoftwareUpgradeProposal:
			return handleMsgSubmitSoftwareUpgradeProposal(ctx, keeper, msg)

		case MsgVote:
			return handleMsgVote(ctx, keeper, msg)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleMsgSubmitProposal(ctx sdk.Context, keeper Keeper, msg MsgSubmitProposal) (*sdk.Result, error) {
	proposal, err := keeper.SubmitProposal(ctx, msg.Content)
	if err != nil {
		return nil, err
	}

	votingStarted, err := keeper.AddDeposit(ctx, proposal.ProposalID, msg.Proposer, msg.InitialDeposit)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Proposer.String()),
		),
	)

	if votingStarted {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeSubmitProposal,
				sdk.NewAttribute(types.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", proposal.ProposalID)),
			),
		)
	}

	return &sdk.Result{
		Data:   keeper.cdc.MustMarshalBinaryLengthPrefixed(proposal.ProposalID),
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleMsgSubmitSoftwareUpgradeProposal(ctx sdk.Context, keeper Keeper, msg MsgSoftwareUpgradeProposal) (*sdk.Result, error) {
	//proposalLevel := GetProposalLevelByProposalKind(msg.Proposal.Type)
	//if num, ok := keeper.HasReachedTheMaxProposalNum(ctx, proposalLevel); ok {
	//	return ErrMoreThanMaxProposal(keeper.codespace, num, proposalLevel.string()).Result()
	//}
	//
	//if !keeper.protocolKeeper.IsValidVersion(ctx, msg.Version) {
	//	return ErrCodeInvalidVersion(keeper.codespace, msg.Version).Result()
	//}
	//
	//if uint64(ctx.BlockHeight()) > msg.SwitchHeight {
	//	return ErrCodeInvalidSwitchHeight(keeper.codespace, uint64(ctx.BlockHeight()), msg.SwitchHeight).Result()
	//}
	//_, found := keeper.guardianKeeper.GetProfiler(ctx, msg.Proposer)
	//if !found {
	//	return ErrNotProfiler(keeper.codespace, msg.Proposer).Result()
	//}
	//
	//if _, ok := keeper.protocolKeeper.GetUpgradeConfig(ctx); ok {
	//	return ErrSwitchPeriodInProcess(keeper.codespace).Result()
	//}
	//
	//proposal := keeper.NewSoftwareUpgradeProposal(ctx, msg)
	//
	//err, votingStarted := keeper.AddInitialDeposit(ctx, proposal, msg.Proposer, msg.InitialDeposit)
	//if err != nil {
	//	return err.Result()
	//}
	//proposalIDBytes := []byte(strconv.FormatUint(proposal.GetProposalID(), 10))
	//
	//resTags := sdk.NewTags(
	//	tags.Proposer, []byte(msg.Proposer.String()),
	//	tags.ProposalID, proposalIDBytes,
	//)
	//
	//if votingStarted {
	//	resTags = resTags.AppendTag(tags.VotingPeriodStart, proposalIDBytes)
	//}
	//
	//keeper.AddProposalNum(ctx, proposal)
	//return sdk.Result{
	//	Data: proposalIDBytes,
	//	Tags: resTags,
	//}
}

func handleMsgDeposit(ctx sdk.Context, keeper Keeper, msg MsgDeposit) (*sdk.Result, error) {
	votingStarted, err := keeper.AddDeposit(ctx, msg.ProposalID, msg.Depositor, msg.Amount)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Depositor.String()),
		),
	)

	if votingStarted {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeProposalDeposit,
				sdk.NewAttribute(types.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", msg.ProposalID)),
			),
		)
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgVote(ctx sdk.Context, keeper Keeper, msg MsgVote) (*sdk.Result, error) {
	err := keeper.AddVote(ctx, msg.ProposalID, msg.Voter, msg.Option)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Voter.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil

}
