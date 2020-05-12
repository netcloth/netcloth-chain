package guardian

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/netcloth/netcloth-chain/app/v0/genutil"
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/server"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func AddGenesisGuardianCmd(ctx *server.Context, cdc *codec.Codec, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-guardian [address] [description]",
		Short: "Add genesis guardian to genesis.json",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			description := args[1]

			genGuardian := NewGuardian(description, Genesis, addr, addr)
			if err := genGuardian.Validate(); err != nil {
				return err
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return err
			}

			var genesisGuardians GenesisState
			cdc.MustUnmarshalJSON(appState[ModuleName], &genesisGuardians)

			if genesisGuardians.Contains(addr) {
				return fmt.Errorf("cannot add guardian at existing address %v", addr)
			}

			genesisGuardians.Profilers = append(genesisGuardians.Profilers, genGuardian)

			genesisStateBz := cdc.MustMarshalJSON(genesisGuardians)
			appState[ModuleName] = genesisStateBz

			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return err
			}

			genDoc.AppState = appStateJSON

			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	return cmd
}
