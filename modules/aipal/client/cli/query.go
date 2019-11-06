package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/NetCloth/netcloth-chain/client"
	"github.com/NetCloth/netcloth-chain/client/context"
	"github.com/NetCloth/netcloth-chain/codec"
	"github.com/NetCloth/netcloth-chain/modules/aipal/types"
	sdk "github.com/NetCloth/netcloth-chain/types"
	"github.com/NetCloth/netcloth-chain/version"
)

func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	ipalQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for aipal",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ipalQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryParams(queryRoute, cdc),
		GetCmdQueryServiceNodeList(cdc),
		GetCmdQueryServiceNode(cdc),
	)...)

	return ipalQueryCmd

}

func GetCmdQueryParams(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current cipal parameters",
		Long: strings.TrimSpace(fmt.Sprintf(`Query values set as aipal parameters.
Example:
$ %s query aipal params`, version.ClientName)),

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

func GetCmdQueryServiceNodeList(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Querying commands for ServiceNodes",
		Long: strings.TrimSpace(fmt.Sprintf(`List all ServiceNodes.
Example:
$ %s query aipal servicenodes`, version.ClientName)),

		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			resKVs, _, err := cliCtx.QuerySubspace(types.ServiceNodeByBondKey, types.StoreKey)
			if err != nil {
				return err
			}

			var serverNodes types.ServiceNodes
			if len(resKVs) > 0 {
				for i := len(resKVs) - 1; i >= 0; i-- {
					serverNodes = append(serverNodes, types.MustUnmarshalServiceNode(cdc, resKVs[i].Value))
				}
			}

			return cliCtx.PrintOutput(serverNodes)
		},
	}
}

func GetCmdQueryServiceNode(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "node",
		Short: "Querying commands for ServiceNode",
		Long: strings.TrimSpace(fmt.Sprintf(`Query ServiceNode by accAddr.
Example:
$ %s query aipal node`, version.ClientName)),

		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryStore(types.GetServiceNodeKey(addr), types.StoreKey)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return fmt.Errorf("No ServiceNode found with address %s", addr)
			}

			return cliCtx.PrintOutput(types.MustUnmarshalServiceNode(cdc, res))
		},
	}
}
