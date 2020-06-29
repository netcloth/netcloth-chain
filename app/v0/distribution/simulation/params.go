package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/netcloth/netcloth-chain/app/v0/distribution/types"
	"github.com/netcloth/netcloth-chain/modules/simulation"
	simtypes "github.com/netcloth/netcloth-chain/types/simulation"
)

const (
	keyCommunityTax        = "communitytax"
	keyBaseProposerReward  = "baseproposerreward"
	keyBonusProposerReward = "bonusproposerreward"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keyCommunityTax,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenCommunityTax(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyBaseProposerReward,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenBaseProposerReward(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyBonusProposerReward,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenBonusProposerReward(r))
			},
		),
	}
}
