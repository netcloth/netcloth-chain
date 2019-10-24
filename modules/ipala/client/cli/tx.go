package cli

import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"

    "github.com/NetCloth/netcloth-chain/client"
    "github.com/NetCloth/netcloth-chain/client/context"
    "github.com/NetCloth/netcloth-chain/codec"
    "github.com/NetCloth/netcloth-chain/modules/auth"
    "github.com/NetCloth/netcloth-chain/modules/auth/client/utils"
    sdk "github.com/NetCloth/netcloth-chain/types"
    "github.com/NetCloth/netcloth-chain/modules/ipala/types"
)

func IPALACmd(cdc *codec.Codec) *cobra.Command {
    txCmd := &cobra.Command{
        Use:   types.ModuleName,
        Short: "IPAL transaction subcommands",
    }
    txCmd.AddCommand(
        ServerNodeClaimCmd(cdc),
    )
    return txCmd
}

func ServerNodeClaimCmd(cdc *codec.Codec) *cobra.Command {
    cmd := &cobra.Command{
        Use:     "claim",
        Short:   "Create and sign a ServerNodeClaim tx",
        Example: "nchcli ipal claim  --from=<user key name> --moniker=<name> --website=<website> --server=<server_endpoint> --details=<details>",
        RunE: func(cmd *cobra.Command, args []string) error {
            txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            moniker := viper.GetString(flagMoniker)
            website := viper.GetString(flagWebsite)
            serverEndPoint := viper.GetString(flagServerEndPoint)
            details := viper.GetString(flagDetails)
            stakeAmount := viper.GetString(flagBond)

            coin, err := sdk.ParseCoin(stakeAmount)
            if err != nil {
                return err
            }

            // build and sign the transaction, then broadcast to Tendermint
            msg := types.NewMsgServiceNodeClaim(cliCtx.GetFromAddress(),moniker, website, serverEndPoint, details, coin)
            if err := msg.ValidateBasic(); err != nil {
                return err
            }
            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }

    cmd.Flags().String(flagMoniker, "", "server node moniker")
    cmd.Flags().String(flagWebsite, "", "server node website")
    cmd.Flags().String(flagServerEndPoint, "", "server node endpoint")
    cmd.Flags().String(flagDetails, "", "server node details")
    cmd.Flags().String(flagBond, "", "stake amount")

    cmd.MarkFlagRequired(flagMoniker)
    cmd.MarkFlagRequired(flagServerEndPoint)
    cmd.MarkFlagRequired(flagBond)

    cmd = client.PostCommands(cmd)[0]

    return cmd
}