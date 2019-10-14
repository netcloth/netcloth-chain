package cli

import (
	"github.com/NetCloth/netcloth-chain/client/keys"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"

	"github.com/NetCloth/netcloth-chain/client"
	"github.com/NetCloth/netcloth-chain/client/context"
	"github.com/NetCloth/netcloth-chain/codec"
	"github.com/NetCloth/netcloth-chain/modules/auth"
	"github.com/NetCloth/netcloth-chain/modules/auth/client/utils"
	"github.com/NetCloth/netcloth-chain/modules/ipal/types"
	sdk "github.com/NetCloth/netcloth-chain/types"
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
		Example: "nchcli ipal claim  --from <key name> --user=<user key name> --ip=<server ip>",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			user, err := sdk.AccAddressFromBech32(viper.GetString(flagUser))
			if err != nil {
				return err
			}
			serverIP := viper.GetString(flagServerIP)

			info, err := txBldr.Keybase().Get(cliCtx.GetFromName())
			if err != nil {
				return err
			}
			userAddress := info.GetAddress().String()

			// build user request signature
			// build msg
			adMsg := types.NewADParam(userAddress, serverIP, time.Now().AddDate(0, 0, 1))
			passphrase, err := keys.GetPassphrase(cliCtx.GetFromName())
			if err != nil {
				return err
			}
			// sign
			sigBytes, pubkey, err := txBldr.Keybase().Sign(info.GetName(), passphrase, adMsg.GetSignBytes())
			if err != nil {
				return err
			}
			stdSig := auth.StdSignature{
				PubKey:    pubkey,
				Signature: sigBytes,
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgIPALClaim(cliCtx.FromAddress(), userAddress, serverIP, time.Now().AddDate(0, 0, 1), stdSig)

			//if err := msg.ValidateBasic(); err != nil {
			//	return err
			//}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagServerIP, "", "server ip")
	cmd.Flags().String(flagUser, "", "proxy account")
	cmd.MarkFlagRequired(flagServerIP)
	cmd.MarkFlagRequired(flagUser)

	cmd = client.PostCommands(cmd)[0]

	return cmd
}
