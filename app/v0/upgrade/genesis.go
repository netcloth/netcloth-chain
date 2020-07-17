package upgrade

import (
	"github.com/netcloth/netcloth-chain/app/v0/upgrade/types"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/version"
)

// GenesisState contains all upgrade state that must be provided at genesis
type GenesisState struct {
	GenesisVersion types.VersionInfo `json:"genesis_version" yaml:"genesis_version"`
}

// DefaultGenesisState returns default raw genesis raw message
func DefaultGenesisState() GenesisState {
	return GenesisState{
		GenesisVersion: types.NewVersionInfo(sdk.DefaultUpgradeConfig(0, "https://github.com/netcloth/netcloth-chain/releases/tag/v"+version.Version), true),
	}
}

// InitGenesis builds the genesis version for first version
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	genesisVersion := data.GenesisVersion

	k.AddNewVersionInfo(ctx, genesisVersion)
	k.protocolKeeper.ClearUpgradeConfig(ctx)
	k.protocolKeeper.SetCurrentVersion(ctx, genesisVersion.UpgradeInfo.Protocol.Version)
}

// ExportGenesis outputs genesis state
func ExportGenesis() GenesisState {
	return GenesisState{
		GenesisVersion: types.NewVersionInfo(sdk.DefaultUpgradeConfig(version.AppVersion, "https://github.com/netcloth/netcloth-chain/releases/tag/v"+version.Version), true),
	}
}
