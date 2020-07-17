package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/netcloth/netcloth-chain/app/v0/simulation"
	"github.com/netcloth/netcloth-chain/app/v0/staking/types"
	simtypes "github.com/netcloth/netcloth-chain/types/simulation"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyMaxValidators),
			func(r *rand.Rand) string {
				return fmt.Sprintf("%d", GenMaxValidators(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyUnbondingTime),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenUnbondingTime(r))
			},
		),
		//simulation.NewSimParamChange(types.ModuleName, string(types.KeyHistoricalEntries),
		//	func(r *rand.Rand) string {
		//		return fmt.Sprintf("%d", GetHistEntries(r))
		//	},
		//),
	}
}
