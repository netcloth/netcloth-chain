package vm

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/app/v0/vm/client/cli"
	"github.com/netcloth/netcloth-chain/app/v0/vm/client/rest"
	"github.com/netcloth/netcloth-chain/app/v0/vm/types"
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
	keeper Keeper
}

func NewAppModule(keeper Keeper) AppModule {
	return AppModule{keeper: keeper}
}

func (a AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	a.keeper.SetParams(ctx, genesisState.Params)

	return nil
}

func (a AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	kvs := a.keeper.StateDB.WithContext(ctx).ExportState()
	fmt.Fprintf(os.Stderr, fmt.Sprintf("len(kvs)=%d", len(kvs)))
	return types.ModuleCdc.MustMarshalJSON(kvs)
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

func (a AppModule) BeginBlock(sdk.Context, abci.RequestBeginBlock) {
	// TODO
}

func (a AppModule) EndBlock(ctx sdk.Context, end abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlocker(ctx, a.keeper)

}
