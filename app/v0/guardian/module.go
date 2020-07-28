package guardian

// DONTCOVER

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

// AppModuleBasic defines the basic application module used by the guardian module.
type AppModuleBasic struct{}

// Name returns the guardian module's name.
func (a AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the guardian module's types for the given codec.
func (a AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	guardiantypes.RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the guardian
// module.
func (a AppModuleBasic) DefaultGenesis() json.RawMessage {
	d, _ := json.Marshal(DefaultGenesisState())
	return d
}

// ValidateGenesis performs genesis state validation
func (a AppModuleBasic) ValidateGenesis(d json.RawMessage) error {
	var gs GenesisState
	return json.Unmarshal(d, &gs)
}

// RegisterRESTRoutes registers the REST routes
func (a AppModuleBasic) RegisterRESTRoutes(context.CLIContext, *mux.Router) {
}

// GetTxCmd returns the root tx command
func (a AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GuardianCmd(cdc)
}

// GetQueryCmd returns the root query command for the guardian module.
func (a AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(cdc)
}

// AppModule implements an application module for the guardian module.
type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{keeper: keeper}
}

// InitGenesis performs genesis initialization for the guardian module. It returns
// no validator updates.
func (a AppModule) InitGenesis(ctx sdk.Context, d json.RawMessage) []types.ValidatorUpdate {
	var gs GenesisState
	err := json.Unmarshal(d, &gs)
	if err != nil {
		panic(err)
	}
	InitGenesis(ctx, a.keeper, gs)
	return nil
}

// ExportGenesis returns the exported genesis state as raw bytes for the guardian
// module.
func (a AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	d, err := json.Marshal(ExportGenesis(ctx, a.keeper))
	if err != nil {
		panic(err)
	}
	return d
}

// RegisterInvariants registers module invariants
func (a AppModule) RegisterInvariants(sdk.InvariantRegistry) {
	panic("implement me")
}

// Route returns the message routing key for the guardian module.
func (a AppModule) Route() string {
	return RouterKey
}

// NewHandler returns an sdk.Handler for the guardian module.
func (a AppModule) NewHandler() sdk.Handler {
	return NewHandler(a.keeper)
}

// QuerierRoute returns the guardian module's querier route name.
func (a AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns no sdk.Querier.
func (a AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(a.keeper)
}

// BeginBlock returns the begin blocker for the guardian module.
func (a AppModule) BeginBlock(sdk.Context, types.RequestBeginBlock) {
	panic("implement me")
}

// EndBlock returns the end blocker for the guardian module. It returns no validator
// updates.
func (a AppModule) EndBlock(ctx sdk.Context, b types.RequestEndBlock) []types.ValidatorUpdate {
	return nil
}
