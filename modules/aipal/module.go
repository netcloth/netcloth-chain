package aipal

import (
    "encoding/json"
    "github.com/NetCloth/netcloth-chain/client/context"
    "github.com/NetCloth/netcloth-chain/codec"
    "github.com/NetCloth/netcloth-chain/modules/aipal/client/cli"
    "github.com/NetCloth/netcloth-chain/modules/aipal/client/rest"
    "github.com/NetCloth/netcloth-chain/modules/aipal/types"
    sdk "github.com/NetCloth/netcloth-chain/types"
    "github.com/NetCloth/netcloth-chain/types/module"
    "github.com/gorilla/mux"
    "github.com/spf13/cobra"
    abci "github.com/tendermint/tendermint/abci/types"
)

var (
    _ module.AppModule = AppModule{}
    _ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct {}

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

var _ module.AppModuleBasic = AppModuleBasic{}

type AppModule struct {
    AppModuleBasic
    k Keeper
}

func NewAppModule(keeper Keeper) AppModule{
    return AppModule{k:keeper}
}

func (a AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
    var genesisState types.GenesisState
    ModuleCdc.MustUnmarshalJSON(data, &genesisState)
    a.k.SetParams(ctx, genesisState.Params)

    return nil
}

func (a AppModule) ExportGenesis(sdk.Context) json.RawMessage {
    panic("implement me")
}

func (a AppModule) RegisterInvariants(sdk.InvariantRegistry) {
    panic("implement me")
}

func (a AppModule) Route() string {
    return RouterKey
}

func (a AppModule) NewHandler() sdk.Handler {
    return NewHandler(a.k)
}

func (a AppModule) QuerierRoute() string {
    return QuerierRoute
}

func (a AppModule) NewQuerierHandler() sdk.Querier {
    return NewQuerier(a.k)
}

func (a AppModule) BeginBlock(sdk.Context, abci.RequestBeginBlock) {
}

func (a AppModule) EndBlock(ctx sdk.Context, end abci.RequestEndBlock) []abci.ValidatorUpdate {
    return EndBlocker(ctx, a.k)
}
