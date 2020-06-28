package simulation

import (
	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"github.com/netcloth/netcloth-chain/app/v0/genaccounts/internal/types"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/types/module"
)

func RandomGenesisAccounts(simState *module.SimulationState) (genesisAccs types.GenesisAccounts) {
	for _, acc := range simState.Accounts {
		bacc := auth.NewBaseAccountWithAddress(acc.Address)
		coins := sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, simState.InitialStake))
		bacc.SetCoins(coins)
		gacc := types.NewGenesisAccount(&bacc)
		genesisAccs = append(genesisAccs, gacc)
	}

	return genesisAccs
}
