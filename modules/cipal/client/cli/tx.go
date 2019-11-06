package cli

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/NetCloth/netcloth-chain/client"
	"github.com/NetCloth/netcloth-chain/client/context"
	"github.com/NetCloth/netcloth-chain/client/keys"
	"github.com/NetCloth/netcloth-chain/codec"
	"github.com/NetCloth/netcloth-chain/modules/auth"
	"github.com/NetCloth/netcloth-chain/modules/auth/client/utils"
	"github.com/NetCloth/netcloth-chain/modules/cipal/types"
	sdk "github.com/NetCloth/netcloth-chain/types"
)

func CIPALCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "CIPAL transaction subcommands",
	}
	txCmd.AddCommand(
		CIPALClaimCmd(cdc),
	)
	return txCmd
}

func CIPALClaimCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "claim",
		Short:   "Create and sign a CIPALClaim tx",
		Example: "nchcli cipal claim --user=<user key name> --proxy=<proxy key name> --service_address=<service address> --service_type=<service type>",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtxUser := context.NewCLIContextWithFrom(viper.GetString(flagUser)).WithCodec(cdc)

			info, err := txBldr.Keybase().Get(cliCtxUser.GetFromName())
			if err != nil {
				return err
			}
			userAddress := info.GetAddress().String()

			serviceAddress := viper.GetString(flagServiceAddress)
			serviceType := viper.GetUint64(flagServiceType)
			expiration := time.Now().UTC().AddDate(0, 0, 1)
			adMsg := types.NewADParam(userAddress, serviceAddress, serviceType, expiration)

			// build msg
			passphrase, err := keys.GetPassphrase(cliCtxUser.GetFromName())
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
			cliCtxProxy := context.NewCLIContextWithFrom(viper.GetString(flagProxy)).WithCodec(cdc)
			msg := types.NewMsgCIPALClaim(cliCtxProxy.GetFromAddress(), userAddress, serviceAddress, serviceType, expiration, stdSig)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtxProxy, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagUser, "", "user account")
	cmd.Flags().String(flagProxy, "", "proxy account")
	cmd.Flags().String(flagServiceAddress, "", "service address")
	cmd.Flags().String(flagServiceType, "", "service type. 1:chatting, 2:storage...")

	cmd.MarkFlagRequired(flagUser)
	cmd.MarkFlagRequired(flagProxy)
	cmd.MarkFlagRequired(flagServiceAddress)
	cmd.MarkFlagRequired(flagServiceType)

	cmd = client.PostCommands(cmd)[0]

	return cmd
}
