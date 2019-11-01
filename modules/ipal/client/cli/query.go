package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/NetCloth/netcloth-chain/client"
	"github.com/NetCloth/netcloth-chain/client/context"
	"github.com/NetCloth/netcloth-chain/codec"
	"github.com/NetCloth/netcloth-chain/modules/ipal/types"
	"github.com/NetCloth/netcloth-chain/version"
)

func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	ipalQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for ipal",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ipalQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryCIPAL(queryRoute, cdc),
	)...)

	return ipalQueryCmd
}

func GetCmdQueryCIPAL(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "ipal",
		Short: "Querying commands for ipal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about an individual ipal object.
	Example:
	$ %s query ipal <user-address>
	`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr := args[0]

			res, _, err := cliCtx.QueryStore(types.GetIPALObjectKey(addr), types.StoreKey)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return fmt.Errorf("No ipal object found with address %s", addr)
			}

			return cliCtx.PrintOutput(types.MustUnmarshalIPALObject(cdc, res))
		},
	}
}
