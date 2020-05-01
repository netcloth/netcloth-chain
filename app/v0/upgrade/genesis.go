package upgrade

import (
	"github.com/netcloth/netcloth-chain/app/v0/upgrade/types"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/version"
)

var protocolV0 = version.AppVersion

type GenesisState struct {
	GenesisVersion types.VersionInfo `json:genesis_version`
}

func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	genesisVersion := data.GenesisVersion

	k.AddNewVersionInfo(ctx, genesisVersion)
	k.protocolKeeper.ClearUpgradeConfig(ctx)
	k.protocolKeeper.SetCurrentVersion(ctx, genesisVersion.UpgradeInfo.Protocol.Version)
}

func ExportGenesis() GenesisState {
	return GenesisState{
		types.NewVersionInfo(sdk.DefaultUpgradeConfig(protocolV0, "https://github.com/netcloth/netcloth-chain/releases/tag/v"+version.Version), true),
	}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		types.NewVersionInfo(sdk.DefaultUpgradeConfig(protocolV0, "https://github.com/netcloth/netcloth-chain/releases/tag/v"+version.Version), true),
	}
}
