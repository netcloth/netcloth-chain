package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/netcloth/netcloth-chain/client"
	"github.com/netcloth/netcloth-chain/client/context"
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/ipal/types"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/version"
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
		GetCmdQueryParams(queryRoute, cdc),
		GetCmdQueryIPALNodeList(cdc),
		GetCmdQueryIPALNode(cdc),
	)...)

	return ipalQueryCmd

}

func GetCmdQueryParams(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current ipal parameters",
		Long: strings.TrimSpace(fmt.Sprintf(`Query values set as ipal parameters.
Example:
$ %s query ipal params`, version.ClientName)),

		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryParameters)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Println("fail")
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(bz, &params)
			return cliCtx.PrintOutput(params)
		},
	}
}

func GetCmdQueryIPALNodeList(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Querying commands for IPALNodes",
		Long: strings.TrimSpace(fmt.Sprintf(`List all IPALNodes.
Example:
$ %s query ipal list`, version.ClientName)),

		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			resKVs, _, err := cliCtx.QuerySubspace(types.IPALNodeByBondKey, types.StoreKey)
			if err != nil {
				return err
			}

			var ipalNodes types.IPALNodes
			if len(resKVs) > 0 {
				for i := len(resKVs) - 1; i >= 0; i-- {
					ipalNodes = append(ipalNodes, types.MustUnmarshalIPALNode(cdc, resKVs[i].Value))
				}
			}

			return cliCtx.PrintOutput(ipalNodes)
		},
	}
}

func GetCmdQueryIPALNode(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "node",
		Short: "Querying commands for IPALNode",
		Long: strings.TrimSpace(fmt.Sprintf(`Query IPALNode by accAddr.
Example:
$ %s query ipal node [address]`, version.ClientName)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryStore(types.GetIPALNodeKey(addr), types.StoreKey)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return fmt.Errorf("No IPALNode found with address %s", addr)
			}

			return cliCtx.PrintOutput(types.MustUnmarshalIPALNode(cdc, res))
		},
	}
}
