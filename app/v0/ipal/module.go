package ipal

// DONTCOVER

import (
	"encoding/json"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/app/v0/ipal/client/cli"
	"github.com/netcloth/netcloth-chain/app/v0/ipal/client/rest"
	"github.com/netcloth/netcloth-chain/app/v0/ipal/keeper"
	"github.com/netcloth/netcloth-chain/app/v0/ipal/simulation"
	"github.com/netcloth/netcloth-chain/app/v0/ipal/types"
	"github.com/netcloth/netcloth-chain/client/context"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/types/module"
	simtypes "github.com/netcloth/netcloth-chain/types/simulation"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

type AppModuleBasic struct{}

func (a AppModuleBasic) Name() string {
	return ModuleName
}

func (a AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

func (a AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (a AppModuleBasic) ValidateGenesis(value json.RawMessage) error {
	var data types.GenesisState
	return ModuleCdc.UnmarshalJSON(value, &data)
}

func (a AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

func (a AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

func (a AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

type AppModule struct {
	AppModuleBasic
	keeper          Keeper
	akForSimulation keeper.AccountKeeper // for simulation
}

func NewAppModule(keeper Keeper) AppModule {
	return AppModule{keeper: keeper}
}

func (am *AppModule) WithAccountKeeper(ak keeper.AccountKeeper) *AppModule {
	am.akForSimulation = ak
	return am
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)

	return InitGenesis(ctx, am.keeper, genesisState)
}

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

func (am AppModule) RegisterInvariants(sdk.InvariantRegistry) {
	panic("implement me")
}

func (am AppModule) Route() string {
	return RouterKey
}

func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

func (am AppModule) QuerierRoute() string {
	return QuerierRoute
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// BeginBlock returns the begin blocker for the ipal module.
func (am AppModule) BeginBlock(sdk.Context, abci.RequestBeginBlock) {
}

func (am AppModule) EndBlock(ctx sdk.Context, end abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlocker(ctx, am.keeper)
}

func (am AppModule) GenerateGenesisState(input *module.SimulationState) {
}

func (am AppModule) ProposalContents(simState module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

func (am AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
	return nil
}

func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return simulation.WeightedOperations(simState.AppParams, simState.Cdc, am.akForSimulation, am.keeper)
}
