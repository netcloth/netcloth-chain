package upgrade

import (
	"fmt"
	"strconv"

	upgtypes "github.com/netcloth/netcloth-chain/app/v0/upgrade/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func EndBlocker(ctx sdk.Context, keeper Keeper) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "EndBlock").With("module", "nch/upgrade"))

	upgradeConfig, ok := keeper.protocolKeeper.GetUpgradeConfig(ctx)
	if ok {
		ctx.Logger().Info(fmt.Sprintf("----upgrade: new upgradeinfo:%s", upgradeConfig.String()))
		//uk.metrics.SetVersion(upgradeConfig.Protocol.Version)

		validator, found := keeper.sk.GetValidatorByConsAddr(ctx, (sdk.ConsAddress)(ctx.BlockHeader().ProposerAddress))
		if !found {
			panic(fmt.Sprintf("validator with consensus-address %s not found", (sdk.ConsAddress)(ctx.BlockHeader().ProposerAddress).String()))
		}

		if ctx.BlockHeader().Version.App == upgradeConfig.Protocol.Version {
			//uk.metrics.SetSignal(validator.GetOperator().String(), upgradeConfig.Protocol.Version)
			keeper.SetSignal(ctx, upgradeConfig.Protocol.Version, validator.ConsAddress().String())
			ctx.Logger().Info("Validator has downloaded the latest software", ", validator", validator.GetOperator().String(), ", version", upgradeConfig.Protocol.Version)
		} else {
			//uk.metrics.DeleteSignal(validator.GetOperator().String(), upgradeConfig.Protocol.Version)
			ok := keeper.DeleteSignal(ctx, upgradeConfig.Protocol.Version, validator.ConsAddress().String())
			if ok {
				ctx.Logger().Info("Validator has restarted the old software", ", validator", validator.GetOperator().String(), ", version", upgradeConfig.Protocol.Version)
			}
		}

		curHeight := uint64(ctx.BlockHeight())
		if curHeight == upgradeConfig.Protocol.Height {
			success := tally(ctx, upgradeConfig.Protocol.Version, keeper, upgradeConfig.Protocol.Threshold)
			if success {
				ctx.Logger().Info("Software Upgrade is successful, ", "version", upgradeConfig.Protocol.Version)
				keeper.protocolKeeper.SetCurrentVersion(ctx, upgradeConfig.Protocol.Version)

				ctx.EventManager().EmitEvent(sdk.NewEvent(
					sdk.AppVersionEvent,
					sdk.NewAttribute(sdk.AppVersionEvent, strconv.FormatUint(keeper.protocolKeeper.GetCurrentVersion(ctx), 10)),
				))
			} else {
				ctx.Logger().Info("Software Upgrade is failure, ", "version", upgradeConfig.Protocol.Version)
				keeper.protocolKeeper.SetLastFailedVersion(ctx, upgradeConfig.Protocol.Version)
			}

			keeper.AddNewVersionInfo(ctx, upgtypes.NewVersionInfo(upgradeConfig, success))
			keeper.protocolKeeper.ClearUpgradeConfig(ctx)
		}

		if curHeight > upgradeConfig.Protocol.Height {
			ctx.Logger().Info(fmt.Sprintf("current height[%d] is big than switch height[%d], failed to switch", ctx.BlockHeight(), upgradeConfig.Protocol.Height))
			keeper.AddNewVersionInfo(ctx, upgtypes.NewVersionInfo(upgradeConfig, false))
			keeper.protocolKeeper.ClearUpgradeConfig(ctx)
		}
	} else {
		//uk.metrics.DeleteVersion()
		ctx.Logger().Debug("----upgrade: no upgradeinfo")
	}
}
