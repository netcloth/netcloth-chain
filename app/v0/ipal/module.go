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

// AppModuleBasic defines the basic application module used by the ipal module.
type AppModuleBasic struct{}

// Name returns the ipal module's name.
func (a AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the ipal module's types for the given codec.
func (a AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the ipal
// module.
func (a AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the ipal module.
func (a AppModuleBasic) ValidateGenesis(value json.RawMessage) error {
	var data types.GenesisState
	return ModuleCdc.UnmarshalJSON(value, &data)
}

// RegisterRESTRoutes registers the REST routes for the ipal module.
func (a AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns no root tx command for the ipal module.
func (a AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

// GetQueryCmd returns the root query command for the ipal module.
func (a AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

// AppModule implements an application module for the ipal module.
type AppModule struct {
	AppModuleBasic
	keeper          Keeper
	akForSimulation keeper.AccountKeeper // for simulation
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{keeper: keeper}
}

func (am *AppModule) WithAccountKeeper(ak keeper.AccountKeeper) *AppModule {
	am.akForSimulation = ak
	return am
}

// InitGenesis performs genesis initialization for the ipal module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)

	return InitGenesis(ctx, am.keeper, genesisState)
}

// ExportGenesis returns the exported genesis state as raw bytes for the ipal
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// RegisterInvariants registers the ipal module invariants.
func (am AppModule) RegisterInvariants(sdk.InvariantRegistry) {
	panic("implement me")
}

// Route returns the message routing key for the ipal module.
func (am AppModule) Route() string {
	return RouterKey
}

// NewHandler returns an sdk.Handler for the ipal module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the ipal module's querier route name.
func (am AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns the ipal module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// BeginBlock returns the begin blocker for the ipal module.
func (am AppModule) BeginBlock(sdk.Context, abci.RequestBeginBlock) {
}

// EndBlock returns the end blocker for the ipal module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, end abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlocker(ctx, am.keeper)
}

// GenerateGenesisState creates a randomized GenState of the ipal module.
func (am AppModule) GenerateGenesisState(input *module.SimulationState) {
}

// ProposalContents doesn't return any content functions for governance proposals.
func (am AppModule) ProposalContents(simState module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized mint param changes for the simulator.
func (am AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
	return nil
}

// WeightedOperations doesn't return any mint module operation.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return simulation.WeightedOperations(simState.AppParams, simState.Cdc, am.akForSimulation, am.keeper)
}
