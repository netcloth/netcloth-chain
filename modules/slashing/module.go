package slashing

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/client/context"
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/slashing/client/cli"
	"github.com/netcloth/netcloth-chain/modules/slashing/client/rest"
	"github.com/netcloth/netcloth-chain/modules/slashing/internal/types"
	stakingkeeper "github.com/netcloth/netcloth-chain/modules/staking/keeper"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/types/module"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the slashing module.
type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the slashing module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec registers the slashing module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the slashing
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the slashing module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	if err := ModuleCdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the slashing module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns the root tx command for the slashing module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

// GetQueryCmd returns no root query command for the slashing module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

//____________________________________________________________________________

// AppModule implements an application module for the slashing module.
type AppModule struct {
	AppModuleBasic

	keeper        Keeper
	accountKeeper types.AccountKeeper
	stakingKeeper stakingkeeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, accountKeeper types.AccountKeeper, stakingKeeper stakingkeeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		accountKeeper:  accountKeeper,
		stakingKeeper:  stakingKeeper,
	}
}

// Name returns the slashing module's name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the slashing module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the slashing module.
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler returns an sdk.Handler for the slashing module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the slashing module's querier route name.
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns the slashing module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the slashing module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, am.stakingKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the slashing
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the slashing module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req, am.keeper)
}

// EndBlock returns the end blocker for the slashing module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
