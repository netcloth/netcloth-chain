package cli

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

type SoftwareUpgradeProposalJSON struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Type         string   `json:"type"`
	Deposit      sdk.Coin `json:"deposit"`
	Version      uint64   `json:"version"`
	Software     string   `json:"software"`
	SwitchHeight uint64   `json:"switch_height"`
	Threshold    sdk.Dec  `json:"threshold"`
}
