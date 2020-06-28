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

// RandomGenesisAccounts returns a slice of account balances. Each account has
// a balance of simState.InitialStake for sdk.DefaultBondDenom.
//func RandomGenesisBalances(simState *module.SimulationState) []types.Balance {
//	genesisBalances := []types.Balance{}
//
//	for _, acc := range simState.Accounts {
//		genesisBalances = append(genesisBalances, types.Balance{
//			Address: acc.Address,
//			Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(simState.InitialStake))),
//		})
//	}
//
//	return genesisBalances
//}

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
