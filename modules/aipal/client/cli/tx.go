package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/NetCloth/netcloth-chain/client"
	"github.com/NetCloth/netcloth-chain/client/context"
	"github.com/NetCloth/netcloth-chain/codec"
	"github.com/NetCloth/netcloth-chain/modules/aipal/types"
	"github.com/NetCloth/netcloth-chain/modules/auth"
	"github.com/NetCloth/netcloth-chain/modules/auth/client/utils"
	sdk "github.com/NetCloth/netcloth-chain/types"
)

func AIPALCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "IPAL transaction subcommands",
	}
	txCmd.AddCommand(
		ServiceNodeClaimCmd(cdc),
	)
	return txCmd
}

func ServiceNodeClaimCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "claim",
		Short:   "Create and sign a ServiceNodeClaim tx",
		Example: "nchcli aipal claim --from=<user key name> --moniker=<name> --website=<website> --endpoints=<endpoints> --details=<details> --bond=<bond>",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			moniker := viper.GetString(flagMoniker)
			website := viper.GetString(flagWebsite)
			endpointsStr := viper.GetString(flagEndPoints)
			details := viper.GetString(flagDetails)
			stakeAmount := viper.GetString(flagBond)

			coin, err := sdk.ParseCoin(stakeAmount)
			if err != nil {
				return err
			}

			endpoints, err := types.EndpointsFromString(endpointsStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgServiceNodeClaim(cliCtx.GetFromAddress(), moniker, website, details, endpoints, coin)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagMoniker, "", "server node moniker")
	cmd.Flags().String(flagWebsite, "", "server node website")
	cmd.Flags().String(flagEndPoints, "", "server node endpoints, in format: serviceType|endpoint,serviceType|endpoint (e.g. 1|192.168.1.100:10000,2|192.168.1.101:20000)")
	cmd.Flags().String(flagDetails, "", "server node details")
	cmd.Flags().String(flagBond, "", "stake amount (e.g. 1000000unch)")

	cmd.MarkFlagRequired(flagMoniker)
	cmd.MarkFlagRequired(flagEndPoints)
	cmd.MarkFlagRequired(flagBond)

	cmd = client.PostCommands(cmd)[0]

	return cmd
}
