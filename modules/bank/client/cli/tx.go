package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/netcloth/netcloth-chain/client"
	"github.com/netcloth/netcloth-chain/client/context"
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/auth"
	"github.com/netcloth/netcloth-chain/modules/auth/client/utils"
	"github.com/netcloth/netcloth-chain/modules/bank/internal/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	flagTo     = "to"
	flagAmount = "amount"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Bank transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		SendTxCmd(cdc),
	)
	return txCmd
}

// SendTxCmd will create a send tx and sign it with the given key.
func SendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "send [from_key_or_address] [to_address] [amount]",
		Short:   "Create and sign a send tx",
		Example: "nchcli send --from <key name> --to=<account address> --chain-id=<chain-id> --amount=10nch",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			to, err := sdk.AccAddressFromBech32(viper.GetString(flagTo))
			if err != nil {
				return err
			}

			// parse coins trying to be sent
			amount := viper.GetString(flagAmount)
			coins, err := sdk.ParseCoins(amount)
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgSend(cliCtx.GetFromAddress(), to, coins)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTo, "", "Bech32 encoding address to receive coins")
	cmd.Flags().String(flagAmount, "", "Amount of coins to send, for instance: 10pnch")
	cmd.MarkFlagRequired(flagTo)
	cmd.MarkFlagRequired(flagAmount)

	cmd = client.PostCommands(cmd)[0]

	return cmd
}
