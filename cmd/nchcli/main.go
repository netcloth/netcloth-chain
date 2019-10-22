package main

import (
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/NetCloth/netcloth-chain/app"
	"github.com/NetCloth/netcloth-chain/client"
	"github.com/NetCloth/netcloth-chain/client/keys"
	"github.com/NetCloth/netcloth-chain/client/lcd"
	"github.com/NetCloth/netcloth-chain/client/rpc"
	"github.com/NetCloth/netcloth-chain/modules/auth"
	authcmd "github.com/NetCloth/netcloth-chain/modules/auth/client/cli"
	authrest "github.com/NetCloth/netcloth-chain/modules/auth/client/rest"
	"github.com/NetCloth/netcloth-chain/modules/bank"
	bankcmd "github.com/NetCloth/netcloth-chain/modules/bank/client/cli"
	ipalcli "github.com/NetCloth/netcloth-chain/modules/ipal/client/cli"
	sdk "github.com/NetCloth/netcloth-chain/types"
	"github.com/NetCloth/netcloth-chain/version"
)

const (
	storeAcc = "acc"
	storeNCH = "nch"

	flagFrom   = "from"
	flagTo     = "to"
	flagAmount = "amount"
)

func main() {
	// Configure cobra to sort commands
	cobra.EnableCommandSorting = false

	// Instantiate the codec for the command line application
	cdc := app.CreateCodec()

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	app.SetBech32AddressPrefixes(config)
	config.Seal()

	rootCmd := &cobra.Command{
		Use:   "nchcli",
		Short: "NCHNetwork Client",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		bankcmd.SendTxCmd(cdc),
		ipalcli.IPALCmd(cdc),
		rpc.StatusCommand(),
		client.ConfigCmd(app.DefaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		client.LineBreak,
		lcd.ServeCommand(cdc, registerRoutes),
		client.LineBreak,
		keys.Commands(),
		client.LineBreak,
		version.Cmd,
		client.NewCompletionCmd(rootCmd, true),
	)

	executor := cli.PrepareMainCmd(rootCmd, "NCH", app.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func registerRoutes(rs *lcd.RestServer) {
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
}

func queryCmd(cdc *amino.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subCommands",
	}

	queryCmd.AddCommand(
		authcmd.GetAccountCmd(cdc),
		client.LineBreak,
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		authcmd.QueryTxCmd(cdc),
		client.LineBreak,
		ipalcli.GetQueryCmd(cdc),
		ipalcli.GetServerNodeCmd(cdc),
		client.LineBreak,
	)
	// add modules' query commands
	app.ModuleBasics.AddQueryCommands(queryCmd, cdc)

	return queryCmd
}

func txCmd(cdc *amino.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		bankcmd.SendTxCmd(cdc),
		client.LineBreak,
		authcmd.GetSignCommand(cdc),
		authcmd.GetMultiSignCommand(cdc),
		client.LineBreak,
		authcmd.GetBroadcastCommand(cdc),
		client.LineBreak,
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(txCmd, cdc)

	// remove auth and bank commands as they're mounted under the root tx command
	var cmdsToRemove []*cobra.Command

	for _, cmd := range txCmd.Commands() {
		if cmd.Use == auth.ModuleName || cmd.Use == bank.ModuleName {
			cmdsToRemove = append(cmdsToRemove, cmd)
		}
	}

	txCmd.RemoveCommand(cmdsToRemove...)

	return txCmd
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	if err := viper.BindPFlag(client.FlagChainID, cmd.PersistentFlags().Lookup(client.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}