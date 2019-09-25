package main

import (
	"encoding/json"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/NetCloth/netcloth-chain/app"
	"github.com/NetCloth/netcloth-chain/baseapp"
	"github.com/NetCloth/netcloth-chain/client"
	"github.com/NetCloth/netcloth-chain/modules/genaccounts"
	genaccscli "github.com/NetCloth/netcloth-chain/modules/genaccounts/client/cli"
	genutilcli "github.com/NetCloth/netcloth-chain/modules/genutil/client/cli"
	"github.com/NetCloth/netcloth-chain/modules/staking"
	"github.com/NetCloth/netcloth-chain/server"
	"github.com/NetCloth/netcloth-chain/store"
	sdk "github.com/NetCloth/netcloth-chain/types"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const (
	flagOverwrite = "overwrite"
)

// nchd custom flags
const flagInvCheckPeriod = "inv-check-period"

var invCheckPeriod uint

func main() {
	cdc := app.CreateCodec()

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	app.SetBech32AddressPrefixes(config)
	config.Seal()

	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "nchd",
		Short:             "nch Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.CollectGenTxsCmd(ctx, cdc, genaccounts.AppModuleBasic{}, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.GenTxCmd(ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{},
		genaccounts.AppModuleBasic{}, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics))
	rootCmd.AddCommand(genaccscli.AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(client.NewCompletionCmd(rootCmd, true))
	rootCmd.AddCommand(replayCmd())
	rootCmd.AddCommand(client.LineBreak)

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "NCH", app.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		0, "Assert registered invariants every N blocks")
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewNCHApp(
		logger, db, traceStore, true, invCheckPeriod,
		baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
	)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
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
