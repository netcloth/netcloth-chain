package cli

import (
	"github.com/NetCloth/netcloth-chain/x/nch"
	"github.com/NetCloth/netcloth-chain/x/token"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagTo     = "to"
	flagAmount = "amount"
)

// SendTxCmd will create a send tx and sign it with the given key.
func SendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "send [to_address] [amount]",
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
			msg := nch.NewMsgSend(cliCtx.GetFromAddress(), to, coins)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTo, "", "Bech32 encoding address to receive coins")
	cmd.Flags().String(flagAmount, "", "Amount of coins to send, for instance: 10nch")
	cmd.MarkFlagRequired(flagTo)
	cmd.MarkFlagRequired(flagAmount)

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

// TokenIssueCmd will create a TokenIssue tx
func TokenIssueCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "issue [to_address] [amount]",
		Short:   "Create and sign a tx to issue coins",
		Example: "nchcli issue --from <key name> --to=<account address> --chain-id=<chain-id> --amount=10nch",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			to, err := sdk.AccAddressFromBech32(viper.GetString(flagTo))
			if err != nil {
				return err
			}

			// parse coins trying to be issued
			amount := viper.GetString(flagAmount)
			coins, err := sdk.ParseCoins(amount)
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := token.NewMsgIssue(cliCtx.GetFromAddress(), to, coins[0])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTo, "", "Bech32 encoding address to receive coins")
	cmd.Flags().String(flagAmount, "", "Amount of coins to send, for instance: 10nch")
	cmd.MarkFlagRequired(flagTo)
	cmd.MarkFlagRequired(flagAmount)

	cmd = client.PostCommands(cmd)[0]

	return cmd
}