package cli

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/netcloth/netcloth-chain/client/context"
	"github.com/netcloth/netcloth-chain/version"
	"github.com/spf13/cobra"

	"github.com/netcloth/netcloth-chain/client"
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/vm/types"

	sdk "github.com/netcloth/netcloth-chain/types"
)

func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	vmQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for ipal",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	vmQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryParams(cdc),
		GetCmdQueryCode(cdc),
	)...)
	return vmQueryCmd
}

func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current vm parameters",
		Long: strings.TrimSpace(fmt.Sprintf(`Query values set as vm parameters.
Example:
$ %s query vm params`, version.ClientName)),

		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			route := fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryParameters)
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

func GetCmdQueryCode(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "code",
		Short: "Querying commands for Contract Code",
		Long: strings.TrimSpace(fmt.Sprintf(`Query Contract Code by accAddr.
Example:
$ %s query vm code [address]`, version.ClientName)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData("custom/vm/code", addr)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return fmt.Errorf("No code found with address %s", args[0])
			}

			dst := make([]byte, 2*len(res))
			hex.Encode(dst, res)

			fmt.Println(string(dst))

			return nil
		},
	}
}
