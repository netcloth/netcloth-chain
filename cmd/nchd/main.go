package main

import (
	"encoding/json"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/netcloth/netcloth-chain/app"
	"github.com/netcloth/netcloth-chain/app/v0/genaccounts"
	genaccscli "github.com/netcloth/netcloth-chain/app/v0/genaccounts/client/cli"
	genutilcli "github.com/netcloth/netcloth-chain/app/v0/genutil/client/cli"
	"github.com/netcloth/netcloth-chain/app/v0/staking"
	"github.com/netcloth/netcloth-chain/client"
	"github.com/netcloth/netcloth-chain/server"
	"github.com/netcloth/netcloth-chain/store"
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	flagMinGasPrices   = "minimum-gas-prices"
	flagInvCheckPeriod = "inv-check-period"
)

var invCheckPeriod uint

func main() {
	cdc := app.MakeLatestCodec()

	config := sdk.GetConfig()
	config.Seal()

	ctx := server.NewDefaultContext()

	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "nchd",
		Short:             "nch Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(genutilcli.InitCmd(ctx, cdc, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.CollectGenTxsCmd(ctx, cdc, genaccounts.AppModuleBasic{}, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.GenTxCmd(ctx, cdc, staking.AppModuleBasic{}, genaccounts.AppModuleBasic{}, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(genutilcli.ValidateGenesisCmd(ctx, cdc))
	rootCmd.AddCommand(genaccscli.AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(client.NewCompletionCmd(rootCmd, true))
	rootCmd.AddCommand(replayCmd())
	rootCmd.AddCommand(client.LineBreak)

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	executor := cli.PrepareBaseCmd(rootCmd, "NCH", app.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod, 0, "Assert registered invariants every N blocks")
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	minGasPrices := viper.GetString(flagMinGasPrices)
	return app.NewNCHApp(
		logger, db, traceStore, true, invCheckPeriod,
		app.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))), app.SetMinGasPrices(minGasPrices),
	)
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	if height != -1 {
		nchApp := app.NewNCHApp(logger, db, traceStore, false, uint(1))
		err := nchApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return nchApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	nchApp := app.NewNCHApp(logger, db, traceStore, true, uint(1))
	return nchApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
