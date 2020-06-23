package simulation

import (
	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"github.com/netcloth/netcloth-chain/app/v0/genaccounts/internal/types"
	"github.com/netcloth/netcloth-chain/types/module"
)

func RandomGenesisAccounts(simState *module.SimulationState) (genesisAccs types.GenesisAccounts) {
	for _, acc := range simState.Accounts {
		bacc := auth.NewBaseAccountWithAddress(acc.Address)
		gacc := types.NewGenesisAccount(&bacc)
		genesisAccs = append(genesisAccs, gacc)
	}

	return genesisAccs
}
