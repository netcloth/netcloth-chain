package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/NetCloth/netcloth-chain/client"
	"github.com/NetCloth/netcloth-chain/client/context"
	"github.com/NetCloth/netcloth-chain/codec"
	"github.com/NetCloth/netcloth-chain/modules/auth"
	"github.com/NetCloth/netcloth-chain/modules/auth/client/utils"
	"github.com/NetCloth/netcloth-chain/modules/ipal/types"
	sdk "github.com/NetCloth/netcloth-chain/types"
)

const (
	flagUserAddress = "address"
	flagServerIP    = "ip"
)

// GetTxCmd returns the transaction commands for this module
func IPALCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "IPAL transaction subcommands",
	}
	txCmd.AddCommand(
		IPALClaimCmd(cdc),
	)
	return txCmd
}

func IPALClaimCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "claim",
		Short:   "Create and sign a IPALClaim tx",
		Example: "nchcli ipal claim  --from <key name> --address=<address> --ip=<ip>",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			userAddress := viper.GetString(flagUserAddress)
			serverIP := viper.GetString(flagServerIP)

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgIPALClaim(cliCtx.GetFromAddress(), userAddress, serverIP)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagUserAddress, "", "user address")
	cmd.Flags().String(flagServerIP, "", "server ip")
	cmd.MarkFlagRequired(flagUserAddress)
	cmd.MarkFlagRequired(flagServerIP)

	cmd = client.PostCommands(cmd)[0]

	return cmd
}
