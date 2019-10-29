package cli

import (
    "github.com/NetCloth/netcloth-chain/client"
    "github.com/NetCloth/netcloth-chain/client/context"
    "github.com/NetCloth/netcloth-chain/codec"
    "github.com/NetCloth/netcloth-chain/modules/auth"
    "github.com/NetCloth/netcloth-chain/modules/auth/client/utils"
    "github.com/NetCloth/netcloth-chain/modules/aipal/types"
    sdk "github.com/NetCloth/netcloth-chain/types"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
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
        Example: "nchcli aipal claim --from=<user key name> --moniker=<name> --website=<website> --server=<server_endpoint> --details=<details>",
        RunE: func(cmd *cobra.Command, args []string) error {
            txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            moniker := viper.GetString(flagMoniker)
            website := viper.GetString(flagWebsite)
            serverEndPoint := viper.GetString(flagServerEndPoint)
            serviceType := viper.GetUint64(flagServiceType)
            details := viper.GetString(flagDetails)
            stakeAmount := viper.GetString(flagBond)

            coin, err := sdk.ParseCoin(stakeAmount)
            if err != nil {
                return err
            }

            msg := types.NewMsgServiceNodeClaim(cliCtx.GetFromAddress(), moniker, website, serviceType, serverEndPoint, details, coin)
            if err := msg.ValidateBasic(); err != nil {
                return err
            }
            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }

    cmd.Flags().String(flagMoniker, "", "server node moniker")
    cmd.Flags().String(flagWebsite, "", "server node website")
    cmd.Flags().String(flagServerEndPoint, "", "server node endpoint")
    cmd.Flags().String(flagServiceType, "1", "service type, 64 bits control 64 kind of service types, bit1:control chatting service, bit2 control storage service")
    cmd.Flags().String(flagDetails, "", "server node details")
    cmd.Flags().String(flagBond, "", "stake amount")

    cmd.MarkFlagRequired(flagMoniker)
    cmd.MarkFlagRequired(flagServerEndPoint)
    cmd.MarkFlagRequired(flagBond)

    cmd = client.PostCommands(cmd)[0]

    return cmd
}
