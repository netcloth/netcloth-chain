package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"

	"github.com/NetCloth/netcloth-chain/client/context"
	"github.com/NetCloth/netcloth-chain/codec"
	"github.com/NetCloth/netcloth-chain/modules/ipal/types"
	"github.com/NetCloth/netcloth-chain/version"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   types.ModuleName,
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

func GetServerNodeCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "server-nodes",
		Short: "Querying commands for ServerNodes",
		Long: strings.TrimSpace(
			fmt.Sprintf(`List all ServerNodes.

Example:
$ %s query ipal server-nodes
`,

				version.ClientName,
			),

		),

		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			resKVs, _, err := cliCtx.QuerySubspace(types.ServerNodeObjectKey, types.StoreKey)
			if err != nil {
				return err
			}

			var serverNodes types.ServerNodeObjects
			for _, kv := range resKVs {
				serverNodes = append(serverNodes, types.MustUnmarshalServerNodeObject(cdc, kv.Value))
			}

			return cliCtx.PrintOutput(serverNodes)
		},
	}
}