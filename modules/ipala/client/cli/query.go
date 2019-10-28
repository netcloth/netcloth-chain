package cli

import (
    "fmt"
    "github.com/NetCloth/netcloth-chain/client"
    "github.com/spf13/cobra"
    "strings"

    "github.com/NetCloth/netcloth-chain/client/context"
    "github.com/NetCloth/netcloth-chain/codec"
    "github.com/NetCloth/netcloth-chain/modules/ipala/types"
    "github.com/NetCloth/netcloth-chain/version"
)

func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
    ipalQueryCmd := &cobra.Command {
        Use:                        types.ModuleName,
        Short:                      "Querying commands for ipala",
        DisableFlagParsing:         true,
        SuggestionsMinimumDistance: 2,
        RunE:                       client.ValidateCmd,
    }

    ipalQueryCmd.AddCommand(client.GetCommands(
        GetCmdQueryParams(queryRoute, cdc),
        GetCmdQueryServerNode(cdc),
    )...)

    return ipalQueryCmd

}

func GetCmdQueryParams(storeName string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command {
        Use: "params",
        Args: cobra.NoArgs,
        Short: "Query the current ipal parameters",
        Long: strings.TrimSpace(fmt.Sprintf(`Query values set as ipala parameters.
Example:
$ %s query ipala params`, version.ClientName,)),

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

func GetCmdQueryServerNode(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command {
        Use: "servicenodes",
        Short: "Querying commands for ServiceNodes",
        Long: strings.TrimSpace(fmt.Sprintf(`List all ServiceNodes.
Example:
$ %s query ipala servicenode`, version.ClientName,)),

        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            resKVs, _, err := cliCtx.QuerySubspace(types.ServiceNodeKey, types.StoreKey)
            if err != nil {
                return err
            }

            var serverNodes types.ServiceNodes
            for _, kv := range resKVs {
                serverNodes = append(serverNodes, types.MustUnmarshalServerNodeObject(cdc, kv.Value))
            }

            return cliCtx.PrintOutput(serverNodes)
        },
    }
}