package distribution

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/modules/distribution/keeper"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// set the proposer for determining distribution during endblock
// and distribute rewards for the previous block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	// determine the total power signing the block
	var previousTotalPower, sumPreviousPrecommitPower int64
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		previousTotalPower += voteInfo.Validator.Power
		if voteInfo.SignedLastBlock {
			sumPreviousPrecommitPower += voteInfo.Validator.Power
		}
	}

	// TODO this is Tendermint-dependent
	if ctx.BlockHeight() > 1 {
		previousProposer := k.GetPreviousProposerConsAddr(ctx)
		k.AllocateTokens(ctx, sumPreviousPrecommitPower, previousTotalPower, previousProposer, req.LastCommitInfo.GetVotes())
	}

	// record the proposer for when we payout on the next block
	consAddr := sdk.ConsAddress(req.Header.ProposerAddress)
	k.SetPreviousProposerConsAddr(ctx, consAddr)
}
