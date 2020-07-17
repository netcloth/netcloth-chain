package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/netcloth/netcloth-chain/app/v0/bank/internal/types"
	"github.com/netcloth/netcloth-chain/app/v0/simulation"
	simtypes "github.com/netcloth/netcloth-chain/types/simulation"
)

const keySendEnabled = "sendenabled"

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keySendEnabled,
			func(r *rand.Rand) string {
				return fmt.Sprintf("%v", GenSendEnabled(r))
			},
		),
	}
}
