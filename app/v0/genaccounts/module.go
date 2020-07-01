package genaccounts

import (
	"encoding/json"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/app/v0/auth/exported"
	"github.com/netcloth/netcloth-chain/app/v0/genaccounts/internal/types"
	"github.com/netcloth/netcloth-chain/app/v0/genaccounts/simulation"
	"github.com/netcloth/netcloth-chain/client/context"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/types/module"
	simtypes "github.com/netcloth/netcloth-chain/types/simulation"
)

var (
	_ module.AppModuleGenesis    = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// app module basics object
type AppModuleBasic struct{}

// module name
func (AppModuleBasic) Name() string {
	return ModuleName
}

// register module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {}

// default genesis state
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(GenesisState{})
}

// module validate genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// register rest routes
func (AppModuleBasic) RegisterRESTRoutes(_ context.CLIContext, _ *mux.Router) {}

// get the root tx command of this module
func (AppModuleBasic) GetTxCmd(_ *codec.Codec) *cobra.Command { return nil }

// get the root query command of this module
func (AppModuleBasic) GetQueryCmd(_ *codec.Codec) *cobra.Command { return nil }

// extra function from sdk.AppModuleBasic
// iterate the genesis accounts and perform an operation at each of them
// - to used by other modules
func (AppModuleBasic) IterateGenesisAccounts(cdc *codec.Codec, appGenesis map[string]json.RawMessage,
	iterateFn func(exported.Account) (stop bool)) {

	genesisState := GetGenesisStateFromAppState(cdc, appGenesis)
	for _, genAcc := range genesisState {
		acc := genAcc.ToAccount()
		if iterateFn(acc) {
			break
		}
	}
}

//___________________________
// app module
type AppModule struct {
	AppModuleBasic
	accountKeeper types.AccountKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(accountKeeper types.AccountKeeper) module.AppModule {
	return module.NewGenesisOnlyAppModule(AppModule{
		AppModuleBasic: AppModuleBasic{},
		accountKeeper:  accountKeeper,
	})
}

func NewSimAppModule(accountKeeper types.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		accountKeeper:  accountKeeper,
	}
}

// module init-genesis
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, ModuleCdc, am.accountKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// module export genesis
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.accountKeeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// simulation
func (am AppModule) GenerateGenesisState(simState *module.SimulationState) {
	genesisAccs := simulation.RandomGenesisAccounts(simState)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesisAccs)
}

func (am AppModule) ProposalContents(simState module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

func (am AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
	return nil
}

func (am AppModule) RegisterStoreDecoder(registry sdk.StoreDecoderRegistry) {
}

func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return nil
}
