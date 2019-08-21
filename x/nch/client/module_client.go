package client

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	nchTxCmd := &cobra.Command{
		Use:   "nch",
		Short: "nch transactions subcommands",
	}

	return nchTxCmd
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	// Group nch queries under a subcommand
	nchTxCmd := &cobra.Command{
		Use:   "nch",
		Short: "Querying commands for the nch module",
	}

	return nchTxCmd
}
