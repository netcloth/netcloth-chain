package v0

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

func NewDefaultGenesisState() sdk.GenesisState {
	return ModuleBasics.DefaultGenesis()
}
