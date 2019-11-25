package cli

import (
	"github.com/spf13/cobra"

	"github.com/netcloth/netcloth-chain/client"
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/vm/types"
)

func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	ipalQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for ipal",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	return ipalQueryCmd
}
