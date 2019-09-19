package genutil

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/NetCloth/netcloth-chain/codec"
	"github.com/NetCloth/netcloth-chain/modules/genutil/types"
	"github.com/NetCloth/netcloth-chain/client/context"
	sdk "github.com/NetCloth/netcloth-chain/types"
	"github.com/NetCloth/netcloth-chain/types/module"
)

var (
	_ module.AppModuleGenesis = AppModule{}
	_ module.AppModuleBasic   = AppModuleBasic{}
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

//___________________________
// app module
type AppModule struct {
	AppModuleBasic
	accountKeeper types.AccountKeeper
	stakingKeeper types.StakingKeeper
	deliverTx     deliverTxfn
}

// NewAppModule creates a new AppModule object
func NewAppModule(accountKeeper types.AccountKeeper,
	stakingKeeper types.StakingKeeper, deliverTx deliverTxfn) module.AppModule {

	return module.NewGenesisOnlyAppModule(AppModule{
		AppModuleBasic: AppModuleBasic{},
		accountKeeper:  accountKeeper,
		stakingKeeper:  stakingKeeper,
		deliverTx:      deliverTx,
	})
}

// module init-genesis
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, ModuleCdc, am.stakingKeeper, am.deliverTx, genesisState)
}

// module export genesis
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return nil
}
