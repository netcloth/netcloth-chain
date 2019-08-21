package app

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type (
	// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
	GenesisState struct {
		AuthData auth.GenesisState   `json:"auth"`
		BankData bank.GenesisState   `json:"bank"`
		Accounts []*auth.BaseAccount `json:"accounts"`
	}
)
