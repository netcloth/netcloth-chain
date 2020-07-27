package gov

// DONTCOVER

import (
	"encoding/json"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/app/v0/gov/client"
	"github.com/netcloth/netcloth-chain/app/v0/gov/client/cli"
	"github.com/netcloth/netcloth-chain/app/v0/gov/client/rest"
	"github.com/netcloth/netcloth-chain/app/v0/gov/types"
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

// AppModuleBasic defines the basic application module used by the gov module.
type AppModuleBasic struct {
	proposalHandlers []client.ProposalHandler // proposal handlers which live in governance cli and rest
}

// NewAppModuleBasic creates a new AppModuleBasic object
func NewAppModuleBasic(proposalHandlers ...client.ProposalHandler) AppModuleBasic {
	return AppModuleBasic{
		proposalHandlers: proposalHandlers,
	}
}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the gov module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec registers the gov module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the gov
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the gov module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := types.ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the gov module.
func (a AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	proposalRESTHandlers := make([]rest.ProposalRESTHandler, 0, len(a.proposalHandlers))
	for _, proposalHandler := range a.proposalHandlers {
		proposalRESTHandlers = append(proposalRESTHandlers, proposalHandler.RESTHandler(ctx))
	}

	rest.RegisterRoutes(ctx, rtr, proposalRESTHandlers)
}

// GetTxCmd returns the root tx command for the gov module.
func (a AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	proposalCLIHandlers := make([]*cobra.Command, 0, len(a.proposalHandlers))
	for _, proposalHandler := range a.proposalHandlers {
		proposalCLIHandlers = append(proposalCLIHandlers, proposalHandler.CLIHandler(cdc))
	}

	return cli.GetTxCmd(StoreKey, cdc, proposalCLIHandlers)
}

// GetQueryCmd returns the root query command for the gov module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

//___________________________
// AppModule implements an application module for the gov module.
type AppModule struct {
	AppModuleBasic
	keeper          Keeper
	supplyKeeper    SupplyKeeper
	akForSimulation AccountKeeper // for simulation
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, supplyKeeper SupplyKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		supplyKeeper:   supplyKeeper,
	}
}

func (am *AppModule) WithAccountKeeper(ak AccountKeeper) *AppModule {
	am.akForSimulation = ak
	return am
}

// Name returns the gov module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants registers module invariants
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	RegisterInvariants(ir, am.keeper)
}

// Route returns the message routing key for the gov module.
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler returns an sdk.Handler for the gov module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the gov module's querier route name.
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns no sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the gov module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, am.supplyKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the gov
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock performs a no-op.
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the gov module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}

// for simulation
// GenerateGenesisState creates a randomized GenState of the gov module.
func (am AppModule) GenerateGenesisState(simState *module.SimulationState) {
	RandomizedGenState(simState)
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return WeightedOperations(simState.AppParams, simState.Cdc, am.akForSimulation, am.keeper, simState.Contents)
}

// ProposalContents returns all the gov content functions used to
// simulate governance proposals.
func (am AppModule) ProposalContents(simState module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized gov param changes for the simulator.
func (am AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
	return nil
}
