package gov

import (
	"fmt"

	"github.com/netcloth/netcloth-chain/app/v0/gov/types"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

// SubmitProposal create new proposal given a content
func (keeper Keeper) SubmitProposal(ctx sdk.Context, content Content, proposer sdk.AccAddress) (Proposal, error) {
	if !keeper.router.HasRoute(content.ProposalRoute()) {
		return types.Proposal{}, types.ErrNoProposalHandlerExists
	}

	if ProposalTypeSoftwareUpgrade == content.ProposalType() {
		if keeper.SoftwareUpgradeProposalExist(ctx) {
			return types.Proposal{}, types.ErrSoftwareUpgradeProposalExist
		}
	}

	proposalID, err := keeper.GetProposalID(ctx)
	if err != nil {
		return Proposal{}, err
	}

	// Execute the proposal content in a cache-wrapped context to validate the
	// actual parameter changes before the proposal proceeds through the
	// governance process. State is not persisted.
	cacheCtx, _ := ctx.CacheContext()
	handler := keeper.router.GetRoute(content.ProposalRoute())
	if err := handler(cacheCtx, content, proposalID, proposer); err != nil {
		return types.Proposal{}, err
	}

	submitTime := ctx.BlockHeader().Time
	depositPeriod := keeper.GetDepositParams(ctx).MaxDepositPeriod

	proposal := NewProposal(content, proposalID, submitTime, submitTime.Add(depositPeriod), proposer)

	keeper.SetProposal(ctx, proposal)
	keeper.InsertInactiveProposalQueue(ctx, proposalID, proposal.DepositEndTime)
	keeper.setProposalID(ctx, proposalID+1)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitProposal,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	if ProposalTypeSoftwareUpgrade == content.ProposalType() {
		keeper.SoftwareUpgradeSet(ctx)
	}

	return proposal, nil
}

// GetProposal get Proposal from store by ProposalID
func (keeper Keeper) GetProposal(ctx sdk.Context, proposalID uint64) (proposal Proposal, ok bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(ProposalKey(proposalID))
	if bz == nil {
		return
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &proposal)
	return proposal, true
}

// SetProposal set a proposal to store
func (keeper Keeper) SetProposal(ctx sdk.Context, proposal Proposal) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposal)
	store.Set(ProposalKey(proposal.ProposalID), bz)
}

// DeleteProposal deletes a proposal from store
func (keeper Keeper) DeleteProposal(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		panic(fmt.Sprintf("couldn't find proposal with id#%d", proposalID))
	}
	keeper.RemoveFromInactiveProposalQueue(ctx, proposalID, proposal.DepositEndTime)
	keeper.RemoveFromActiveProposalQueue(ctx, proposalID, proposal.VotingEndTime)
	store.Delete(ProposalKey(proposalID))
}

// GetProposals returns all the proposals from store
func (keeper Keeper) GetProposals(ctx sdk.Context) (proposals Proposals) {
	keeper.IterateProposals(ctx, func(proposal types.Proposal) bool {
		proposals = append(proposals, proposal)
		return false
	})
	return
}

// GetProposalsFiltered get Proposals from store by ProposalID
// voterAddr will filter proposals by whether or not that address has voted on them
// depositorAddr will filter proposals by whether or not that address has deposited to them
// status will filter proposals by status
// numLatest will fetch a specified number of the most recent proposals, or 0 for all proposals
func (keeper Keeper) GetProposalsFiltered(ctx sdk.Context, voterAddr sdk.AccAddress, depositorAddr sdk.AccAddress, status ProposalStatus, numLatest uint64) []Proposal {

	maxProposalID, err := keeper.GetProposalID(ctx)
	if err != nil {
		return []Proposal{}
	}

	matchingProposals := []Proposal{}

	if numLatest == 0 {
		numLatest = maxProposalID
	}

	for proposalID := maxProposalID - numLatest; proposalID < maxProposalID; proposalID++ {
		if voterAddr != nil && len(voterAddr) != 0 {
			_, found := keeper.GetVote(ctx, proposalID, voterAddr)
			if !found {
				continue
			}
		}

		if depositorAddr != nil && len(depositorAddr) != 0 {
			_, found := keeper.GetDeposit(ctx, proposalID, depositorAddr)
			if !found {
				continue
			}
		}

		proposal, ok := keeper.GetProposal(ctx, proposalID)
		if !ok {
			continue
		}

		if ValidProposalStatus(status) && proposal.Status != status {
			continue
		}

		matchingProposals = append(matchingProposals, proposal)
	}
	return matchingProposals
}

// GetProposalID gets the highest proposal ID
func (keeper Keeper) GetProposalID(ctx sdk.Context) (proposalID uint64, err error) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(ProposalIDKey)
	if bz == nil {
		return 0, sdkerrors.Wrap(types.ErrInvalidGenesis, "initial proposal ID hasn't been set")
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &proposalID)
	return proposalID, nil
}

// setProposalID sets the new proposal ID to the store
func (keeper Keeper) setProposalID(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID)
	store.Set(ProposalIDKey, bz)
}

// ActivateVotingPeriod - active voting period
func (keeper Keeper) ActivateVotingPeriod(ctx sdk.Context, proposal Proposal) { //TODO rename to activateVotingPeriod
	proposal.VotingStartTime = ctx.BlockHeader().Time
	votingPeriod := keeper.GetVotingParams(ctx).VotingPeriod
	proposal.VotingEndTime = proposal.VotingStartTime.Add(votingPeriod)
	proposal.Status = StatusVotingPeriod
	keeper.SetProposal(ctx, proposal)

	keeper.RemoveFromInactiveProposalQueue(ctx, proposal.ProposalID, proposal.DepositEndTime)
	keeper.InsertActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndTime)
}
