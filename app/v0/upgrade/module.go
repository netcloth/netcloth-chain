package upgrade

// DONTCOVER

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/app/v0/upgrade/client/cli"
	"github.com/netcloth/netcloth-chain/app/v0/upgrade/client/rest"
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

func (a AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
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
	EndBlocker(ctx, a.keeper)
	return nil
}
