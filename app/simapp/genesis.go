package simapp

import (
	v0 "github.com/netcloth/netcloth-chain/app/v0"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() sdk.GenesisState {
	return v0.ModuleBasics.DefaultGenesis()
}
