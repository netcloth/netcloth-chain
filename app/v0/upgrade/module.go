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

// check the implementation of the interface
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic is a struct of app module basics object
type AppModuleBasic struct{}

// Name returns module name
func (a AppModuleBasic) Name() string {
	return upgtypes.ModuleName
}

// RegisterCodec registers module codec
func (a AppModuleBasic) RegisterCodec(*codec.Codec) {
}

// DefaultGenesis returns default genesis state
func (a AppModuleBasic) DefaultGenesis() json.RawMessage {
	d, _ := json.Marshal(DefaultGenesisState())
	return d
}

// ValidateGenesis validates genesis
func (a AppModuleBasic) ValidateGenesis(d json.RawMessage) error {
	var gs GenesisState
	return json.Unmarshal(d, &gs)
}

// RegisterRESTRoutes register rest routes
func (a AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns the transaction commands for this module
func (a AppModuleBasic) GetTxCmd(*codec.Codec) *cobra.Command {
	return nil
}

// GetQueryCmd gets the root query command of this module
func (a AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(upgtypes.ModuleName, cdc)
}

// AppModule is a struct of app module
type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// NewAppModule creates a new AppModule object for upgrade module
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{keeper: keeper}
}

// InitGenesis initializes module genesis
func (a AppModule) InitGenesis(ctx sdk.Context, d json.RawMessage) []types.ValidatorUpdate {
	var gs GenesisState
	err := json.Unmarshal(d, &gs)
	if err != nil {
		panic(err)
	}
	InitGenesis(ctx, a.keeper, gs)
	return nil
}

// ExportGenesis exports module genesis
func (a AppModule) ExportGenesis(sdk.Context) json.RawMessage {
	d, err := json.Marshal(ExportGenesis())
	if err != nil {
		panic(err)
	}
	return d
}

// RegisterInvariants performs a no-op.
func (a AppModule) RegisterInvariants(sdk.InvariantRegistry) {
	return
}

// Route returns module message route name
func (a AppModule) Route() string {
	return upgtypes.RouterKey
}

// nolint
func (a AppModule) NewHandler() sdk.Handler {
	return nil
}

// QuerierRoute returns module querier route name
func (a AppModule) QuerierRoute() string {
	return upgtypes.QuerierRoute
}

// nolint
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
