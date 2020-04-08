package upgrade

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/app/v0/upgrade/client/cli"
	upgtypes "github.com/netcloth/netcloth-chain/app/v0/upgrade/types"
	"github.com/netcloth/netcloth-chain/client/context"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/types/module"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct{}

func (a AppModuleBasic) Name() string {
	return upgtypes.ModuleName
}

func (a AppModuleBasic) RegisterCodec(*codec.Codec) {
}

func (a AppModuleBasic) DefaultGenesis() json.RawMessage {
	d, _ := json.Marshal(DefaultGenesisState())
	return d
}

func (a AppModuleBasic) ValidateGenesis(d json.RawMessage) error {
	var gs GenesisState
	return json.Unmarshal(d, &gs)
}

func (a AppModuleBasic) RegisterRESTRoutes(context.CLIContext, *mux.Router) {
	panic("implement me")
}

func (a AppModuleBasic) GetTxCmd(*codec.Codec) *cobra.Command {
	return nil
}

func (a AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(upgtypes.ModuleName, cdc)
}

type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

func NewAppModule(keeper Keeper) AppModule {
	return AppModule{keeper: keeper}
}

func (a AppModule) InitGenesis(ctx sdk.Context, d json.RawMessage) []types.ValidatorUpdate {
	var gs GenesisState
	err := json.Unmarshal(d, &gs)
	if err != nil {
		panic(err)
	}
	InitGenesis(ctx, a.keeper, gs)
	return nil
}

func (a AppModule) ExportGenesis(sdk.Context) json.RawMessage {
	d, err := json.Marshal(ExportGenesis())
	if err != nil {
		panic(err)
	}
	return d
}

func (a AppModule) RegisterInvariants(sdk.InvariantRegistry) {
	panic("implement me")
}

func (a AppModule) Route() string {
	return upgtypes.RouterKey
}

func (a AppModule) NewHandler() sdk.Handler {
	return nil
}

func (a AppModule) QuerierRoute() string {
	return upgtypes.QuerierRoute
}

func (a AppModule) NewQuerierHandler() sdk.Querier {
	return nil
}

func (a AppModule) BeginBlock(sdk.Context, types.RequestBeginBlock) {
	panic("implement me")
}

func (a AppModule) EndBlock(ctx sdk.Context, b types.RequestEndBlock) []types.ValidatorUpdate {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "nch/upgrade"))

	upgradeConfig, ok := a.keeper.protocolKeeper.GetUpgradeConfig(ctx)
	if ok {
		ctx.Logger().Error(fmt.Sprintf("----upgrade: yes upgradeinfo:%s", upgradeConfig.String()))
		//uk.metrics.SetVersion(upgradeConfig.Protocol.Version)

		validator, found := a.keeper.sk.GetValidatorByConsAddr(ctx, (sdk.ConsAddress)(ctx.BlockHeader().ProposerAddress))
		if !found {
			panic(fmt.Sprintf("validator with consensus-address %s not found", (sdk.ConsAddress)(ctx.BlockHeader().ProposerAddress).String()))
		}

		if ctx.BlockHeader().Version.App == upgradeConfig.Protocol.Version {
			a.keeper.SetSignal(ctx, upgradeConfig.Protocol.Version, validator.ConsAddress().String())
			//uk.metrics.SetSignal(validator.GetOperator().String(), upgradeConfig.Protocol.Version)

			ctx.Logger().Info("Validator has downloaded the latest software ", "validator", validator.GetOperator().String(), "version", upgradeConfig.Protocol.Version)

		} else {

			ok := a.keeper.DeleteSignal(ctx, upgradeConfig.Protocol.Version, validator.ConsAddress().String())
			//uk.metrics.DeleteSignal(validator.GetOperator().String(), upgradeConfig.Protocol.Version)

			if ok {
				ctx.Logger().Info("Validator has restarted the old software ",
					"validator", validator.GetOperator().String(), "version", upgradeConfig.Protocol.Version)
			}
		}

		if uint64(ctx.BlockHeight())+1 == upgradeConfig.Protocol.Height {
			success := tally(ctx, upgradeConfig.Protocol.Version, a.keeper, upgradeConfig.Protocol.Threshold)

			if success {
				ctx.Logger().Info("Software Upgrade is successful.", "version", upgradeConfig.Protocol.Version)
				a.keeper.protocolKeeper.SetCurrentVersion(ctx, upgradeConfig.Protocol.Version)
			} else {
				ctx.Logger().Info("Software Upgrade is failure.", "version", upgradeConfig.Protocol.Version)
				a.keeper.protocolKeeper.SetLastFailedVersion(ctx, upgradeConfig.Protocol.Version)
			}

			a.keeper.AddNewVersionInfo(ctx, upgtypes.NewVersionInfo(upgradeConfig, success))
			a.keeper.protocolKeeper.ClearUpgradeConfig(ctx)
			a.keeper.gk.SoftwareUpgradeClear(ctx)
		}

		if uint64(ctx.BlockHeight())+1 > upgradeConfig.Protocol.Height {
			ctx.Logger().Info(fmt.Sprintf("current height[%d] is big than switch height[%d], failed to switch", ctx.BlockHeight()+1, upgradeConfig.Protocol.Height))
			a.keeper.AddNewVersionInfo(ctx, upgtypes.NewVersionInfo(upgradeConfig, false))
			a.keeper.protocolKeeper.ClearUpgradeConfig(ctx)
			a.keeper.gk.SoftwareUpgradeClear(ctx)
		}
	} else {
		//uk.metrics.DeleteVersion()
		ctx.Logger().Error("----upgrade: no upgradeinfo")
	}

	e := sdk.NewEvent(
		sdk.AppVersionEvent,
		sdk.NewAttribute(sdk.AppVersionEvent, strconv.FormatUint(a.keeper.protocolKeeper.GetCurrentVersion(ctx), 10)),
	)

	ctx.EventManager().EmitEvent(e)

	return nil
}
