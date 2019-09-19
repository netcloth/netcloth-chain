// nolint
package cli

import (
	"github.com/spf13/cobra"

	"github.com/NetCloth/netcloth-chain/codec"
	"github.com/NetCloth/netcloth-chain/modules/auth"
	"github.com/NetCloth/netcloth-chain/modules/auth/client/utils"
	"github.com/NetCloth/netcloth-chain/modules/crisis/internal/types"
	"github.com/NetCloth/netcloth-chain/client"
	"github.com/NetCloth/netcloth-chain/client/context"
	sdk "github.com/NetCloth/netcloth-chain/types"
)

// command to replace a delegator's withdrawal address
func GetCmdInvariantBroken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invariant-broken [module-name] [invariant-route]",
		Short: "submit proof that an invariant broken to halt the chain",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			senderAddr := cliCtx.GetFromAddress()
			moduleName, route := args[0], args[1]
			msg := types.NewMsgVerifyInvariant(senderAddr, moduleName, route)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Crisis transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(client.PostCommands(
		GetCmdInvariantBroken(cdc),
	)...)
	return txCmd
}
