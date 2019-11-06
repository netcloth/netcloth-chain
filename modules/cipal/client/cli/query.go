package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/NetCloth/netcloth-chain/client"
	"github.com/NetCloth/netcloth-chain/client/context"
	"github.com/NetCloth/netcloth-chain/codec"
	"github.com/NetCloth/netcloth-chain/modules/cipal/types"
	"github.com/NetCloth/netcloth-chain/version"
)

func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cipalQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for cipal",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cipalQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryCIPAL(queryRoute, cdc),
	)...)

	return cipalQueryCmd
}

func GetCmdQueryCIPAL(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "cipal",
		Short: "Querying commands for cipal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about an individual cipal object.
	Example:
	$ %s query cipal <user-address>
	`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr := args[0]

			res, _, err := cliCtx.QueryStore(types.GetCIPALObjectKey(addr), types.StoreKey)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return fmt.Errorf("No cipal object found with address %s", addr)
			}

			return cliCtx.PrintOutput(types.MustUnmarshalCIPALObject(cdc, res))
		},
	}
}
