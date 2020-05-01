package guardian

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/app/v0/guardian/client/cli"
	guardiantypes "github.com/netcloth/netcloth-chain/app/v0/guardian/types"
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
	return ModuleName
}

func (a AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	guardiantypes.RegisterCodec(cdc)
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
}

func (a AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GuardianCmd(cdc)
}

func (a AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(cdc)
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

func (a AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	d, err := json.Marshal(ExportGenesis(ctx, a.keeper))
	if err != nil {
		panic(err)
	}
	return d
}

func (a AppModule) RegisterInvariants(sdk.InvariantRegistry) {
	panic("implement me")
}

func (a AppModule) Route() string {
	return RouterKey
}

func (a AppModule) NewHandler() sdk.Handler {
	return NewHandler(a.keeper)
}

func (a AppModule) QuerierRoute() string {
	return QuerierRoute
}

func (a AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(a.keeper)
}

func (a AppModule) BeginBlock(sdk.Context, types.RequestBeginBlock) {
	panic("implement me")
}

func (a AppModule) EndBlock(ctx sdk.Context, b types.RequestEndBlock) []types.ValidatorUpdate {
	return nil
}
