package simulation

import (
	"math/rand"

	"github.com/netcloth/netcloth-chain/modules/gov/types"
	"github.com/netcloth/netcloth-chain/modules/simulation"
	simappparams "github.com/netcloth/netcloth-chain/simapp/params"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// OpWeightSubmitTextProposal app params key for text proposal
const OpWeightSubmitTextProposal = "op_weight_submit_text_proposal"

// ProposalContents defines the module weighted proposals' contents
func ProposalContents() []simulation.WeightedProposalContent {
	return []simulation.WeightedProposalContent{
		{
			AppParamsKey:       OpWeightSubmitTextProposal,
			DefaultWeight:      simappparams.DefaultWeightTextProposal,
			ContentSimulatorFn: SimulateTextProposalContent,
		},
	}
}

// SimulateTextProposalContent returns a random text proposal content.
func SimulateTextProposalContent(r *rand.Rand, _ sdk.Context, _ []simulation.Account) types.Content {
	return types.NewTextProposal(
		simulation.RandStringOfLength(r, 140),
		simulation.RandStringOfLength(r, 5000),
	)
}
