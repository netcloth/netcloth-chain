package vm

// DONTCOVER

import (
	"encoding/json"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/app/v0/vm/client/cli"
	"github.com/netcloth/netcloth-chain/app/v0/vm/client/rest"
	"github.com/netcloth/netcloth-chain/app/v0/vm/simulation"
	"github.com/netcloth/netcloth-chain/app/v0/vm/types"
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
	return types.ModuleName
}

func (a AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

func (a AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (a AppModuleBasic) ValidateGenesis(value json.RawMessage) error {
	var data types.GenesisState
	if err := types.ModuleCdc.UnmarshalJSON(value, &data); err != nil {
		return err
	}

	return ValidateGenesis(data)
}

func (a AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

func (a AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

func (a AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(types.StoreKey, cdc)
}

var _ module.AppModuleBasic = AppModuleBasic{}

type AppModule struct {
	AppModuleBasic
	keeper          Keeper
	akForSimulation AccountKeeper // for simulation
}

func (am *AppModule) WithAccountKeeper(ak AccountKeeper) *AppModule {
	am.akForSimulation = ak
	return am
}

// InitGenesis instantiates the genesis state
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	am.keeper.SetParams(ctx, genesisState.Params)

	am.keeper.StateDB.WithContext(ctx).ImportState(genesisState)

	return nil
}

// ExportGenesis exports the genesis state to be used by daemon
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	vmState := am.keeper.StateDB.WithContext(ctx).ExportState()
	vmState.Params = am.keeper.GetParams(ctx)
	return ModuleCdc.MustMarshalJSON(vmState)
}

func (am AppModule) RegisterInvariants(sdk.InvariantRegistry) {
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

// BeginBlock function for module at start of each block
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
}

// EndBlock function for module at end of block
func (am AppModule) EndBlock(ctx sdk.Context, end abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlocker(ctx, am.keeper)
}

func NewAppModule(keeper Keeper) AppModule {
	return AppModule{keeper: keeper}
}

//____________________________________________________________________________

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the vm module.
func (am AppModule) GenerateGenesisState(input *module.SimulationState) {
}

// ProposalContents doesn't return any content functions for governance proposals.
func (am AppModule) ProposalContents(simState module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized staking param changes for the simulator.
func (am AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
	return nil
}

// WeightedOperations returns the all the staking module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return simulation.WeightedOperations(simState.AppParams, simState.Cdc, am.akForSimulation, am.keeper)
}
