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
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "IPAL transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		SendTxCmd(cdc),
	)
	return txCmd
}

func SendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ipal [from_key_or_address] [user_address] [server_ip]",
		Short: "Create and sign a IPALClaim tx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

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
