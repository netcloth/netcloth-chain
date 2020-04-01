package main

import (
	v0 "github.com/netcloth/netcloth-chain/app/v0"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/netcloth/netcloth-chain/app"
	"github.com/netcloth/netcloth-chain/client"
	"github.com/netcloth/netcloth-chain/client/keys"
	"github.com/netcloth/netcloth-chain/client/lcd"
	"github.com/netcloth/netcloth-chain/client/rpc"
	"github.com/netcloth/netcloth-chain/modules/auth"
	authcmd "github.com/netcloth/netcloth-chain/modules/auth/client/cli"
	authrest "github.com/netcloth/netcloth-chain/modules/auth/client/rest"
	"github.com/netcloth/netcloth-chain/modules/bank"
	bankcmd "github.com/netcloth/netcloth-chain/modules/bank/client/cli"
	cipalcli "github.com/netcloth/netcloth-chain/modules/cipal/client/cli"
	ipalcli "github.com/netcloth/netcloth-chain/modules/ipal/client/cli"
	vmcli "github.com/netcloth/netcloth-chain/modules/vm/client/cli"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/version"
)

func main() {
	cobra.EnableCommandSorting = false

	cdc := app.MakeLatestCodec()

	config := sdk.GetConfig()
	config.Seal()

	rootCmd := &cobra.Command{
		Use:   "nchcli",
		Short: "NCHNetwork Client",
	}

	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	rootCmd.AddCommand(
		bankcmd.SendTxCmd(cdc),
		cipalcli.CIPALCmd(cdc),
		ipalcli.IPALCmd(cdc),
		vmcli.VMCmd(cdc),
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
	v0.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
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
		client.LineBreak,
	)

	v0.ModuleBasics.AddQueryCommands(queryCmd, cdc)

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

	v0.ModuleBasics.AddTxCommands(txCmd, cdc)

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
