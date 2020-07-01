package simulation

// DONTCOVER

import (
	"math/rand"

	"github.com/netcloth/netcloth-chain/app/v0/bank/internal/types"
	"github.com/netcloth/netcloth-chain/types/module"
)

// Simulation parameter constants
const (
	SendEnabled = "send_enabled"
)

// GenSendEnabled randomized SendEnabled
func GenSendEnabled(r *rand.Rand) bool {
	return r.Int63n(101) <= 95 // 95% chance of transfers being enabled
}

// RandomizedGenState generates a random GenesisState for bank
func RandomizedGenState(simState *module.SimulationState) {
	var sendEnabled bool
	simState.AppParams.GetOrGenerate(
		simState.Cdc, SendEnabled, &sendEnabled, simState.Rand,
		func(r *rand.Rand) { sendEnabled = GenSendEnabled(r) },
	)

	type GenesisState struct {
		SendEnabled bool `json:"send_enabled" yaml:"send_enabled"`
	}

	bankGenesis := GenesisState{SendEnabled: sendEnabled}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(bankGenesis)
}
