package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"github.com/netcloth/netcloth-chain/app/v0/auth/client/utils"
	"github.com/netcloth/netcloth-chain/app/v0/ipal/types"
	"github.com/netcloth/netcloth-chain/client"
	"github.com/netcloth/netcloth-chain/client/context"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func IPALCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "IPAL transaction subcommands",
	}
	txCmd.AddCommand(
		IPALNodeClaimCmd(cdc),
	)
	return txCmd
}

func IPALNodeClaimCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "claim",
		Short:   "Create and sign a IPALNodeClaim tx",
		Example: "nchcli ipal claim --from=<user key name> --moniker=<name> --website=<website> --endpoints=<endpoints> --details=<details> --extension=<extension> --bond=<bond>",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			endpointDelimiter := viper.GetString(flagEndpointDelimiter)
			endpointTypeDelimiter := viper.GetString(flagEndpointTypeDelimiter)
			endpointsStr := viper.GetString(flagEndpoints)
			endpoints, err := types.EndpointsFromString(endpointsStr, endpointDelimiter, endpointTypeDelimiter)
			if err != nil {
				return err
			}

			for _, ep := range endpoints {
				fmt.Println("---", ep)
			}

			moniker := viper.GetString(flagMoniker)
			website := viper.GetString(flagWebsite)
			details := viper.GetString(flagDetails)
			extension := viper.GetString(flagExtension)
			stakeAmount := viper.GetString(flagBond)

			coin, err := sdk.ParseCoin(stakeAmount)
			if err != nil {
				return err
			}

			msg := types.NewMsgIPALNodeClaim(cliCtx.GetFromAddress(), moniker, website, details, extension, endpoints, coin)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagMoniker, "", "ipal node moniker")
	cmd.Flags().String(flagWebsite, "", "ipal node website")
	cmd.Flags().String(flagEndpoints, "", "ipal node endpoints, in format: serviceType|endpoint,serviceType|endpoint (e.g. 1|192.168.1.100:10000,2|192.168.1.101:20000)")
	cmd.Flags().String(flagEndpointDelimiter, ",", "endpoints delimiter, e.g. '#' as delimiter: 1|192.168.1.100:10000#2|192.168.1.101:20000")
	cmd.Flags().String(flagEndpointTypeDelimiter, "|", "endpoint delimiter, e.g. '-' as delimiter: 1-192.168.1.100:10000,2-192.168.1.101:20000")
	cmd.Flags().String(flagDetails, "", "ipal node details")
	cmd.Flags().String(flagExtension, "", "extension for future user define")
	cmd.Flags().String(flagBond, "", "stake amount (e.g. 1000000pnch)")

	cmd.MarkFlagRequired(flagMoniker)
	cmd.MarkFlagRequired(flagEndpoints)
	cmd.MarkFlagRequired(flagBond)

	cmd = client.PostCommands(cmd)[0]

	return cmd
}
