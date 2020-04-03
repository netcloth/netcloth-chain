package cli

import (
	"fmt"
	"github.com/netcloth/netcloth-chain/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"github.com/netcloth/netcloth-chain/app/v0/auth/client/utils"
	"github.com/netcloth/netcloth-chain/app/v0/guardian/types"
	"github.com/netcloth/netcloth-chain/client/context"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func GuardianCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "guardian transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(client.PostCommands(
		GetCmdCreateProfiler(cdc),
		GetCmdDeleteProfiler(cdc),
	)...)

	return txCmd
}

func GetCmdCreateProfiler(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add-profiler",
		Short:   "Add a new profiler",
		Example: "nchcli guardian add-profiler --from=<key-name> --address=<added address> --description=<name>",

		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			profilerAddressStr := viper.GetString(FlagAddress)
			if len(profilerAddressStr) == 0 {
				return fmt.Errorf("must use --address flag")
			}

			profilerAddr, err := sdk.AccAddressFromBech32(profilerAddressStr)
			if err != nil {
				return err
			}

			description := viper.GetString(FlagDescription)
			if len(description) == 0 {
				return fmt.Errorf("must use --description flag")
			}

			msg := types.NewMsgAddProfiler(description, profilerAddr, fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagAddress, "", "bech32 encoded account address")
	cmd.Flags().String(FlagDescription, "", "bdescription of account")

	cmd.MarkFlagRequired(FlagAddress)
	cmd.MarkFlagRequired(FlagDescription)

	return cmd
}

func GetCmdDeleteProfiler(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete-profiler",
		Short:   "Delete a profiler",
		Example: "nchcli guardian delete-profiler --from=<key-name> --address=<deleted address>",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			profilerAddressStr := viper.GetString(FlagAddress)
			if len(profilerAddressStr) == 0 {
				return fmt.Errorf("must use --address flag")
			}

			profilerAddr, err := sdk.AccAddressFromBech32(profilerAddressStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteProfiler(profilerAddr, fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagAddress, "", "bech32 encoded account address")
	cmd.MarkFlagRequired(FlagAddress)

	return cmd
}
