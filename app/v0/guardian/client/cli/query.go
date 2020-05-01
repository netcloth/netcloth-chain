package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/netcloth/netcloth-chain/app/v0/guardian/types"
	"github.com/netcloth/netcloth-chain/client"
	"github.com/netcloth/netcloth-chain/client/context"
	"github.com/netcloth/netcloth-chain/codec"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	guardianQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for ipal",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	guardianQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryProfilers(cdc),
	)...)

	return guardianQueryCmd

}

func GetCmdQueryProfilers(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "profilers",
		Short:   "Query for all profilers",
		Example: "nchcli query guardian profilers",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", "guardian", types.QueryProfilers), nil)

			if err != nil {
				return err
			}

			var profilers types.Profilers
			err = cdc.UnmarshalJSON(res, &profilers)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(profilers)
		},
	}
	return cmd
}
