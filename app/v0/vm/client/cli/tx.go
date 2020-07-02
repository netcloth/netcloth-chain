package cli

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"github.com/netcloth/netcloth-chain/app/v0/auth/client/utils"
	"github.com/netcloth/netcloth-chain/app/v0/vm/types"
	"github.com/netcloth/netcloth-chain/client"
	"github.com/netcloth/netcloth-chain/client/context"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func VMCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "VM transaction subcommands",
	}
	txCmd.AddCommand(
		ContractCreateCmd(cdc),
		ContractCallCmd(cdc),
	)
	return txCmd
}

func ContractCreateCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create and sign a create contract tx",
		Example: "nchcli vm create --from=<user key name> --code_file=<code file> --amount=<amount> --args='arg1 arg2 arg3' --abi_file=<abi_file>",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			coin := sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(0))
			amount := viper.GetString(flagAmount)
			if len(amount) > 0 {
				coinInput, err := sdk.ParseCoin(amount)
				if err != nil {
					return err
				}
				coin = coinInput
			}

			codeFile := viper.GetString(flagCodeFile)
			code, err := CodeFromFile(codeFile)
			if err != nil {
				return err
			}

			var payload []byte
			argList := viper.GetStringSlice(flagArgs)
			if len(argList) != 0 {
				abiFile := viper.GetString(flagAbiFile)
				if len(abiFile) == 0 {
					return fmt.Errorf("must use --abi_file to appoint abi file when use constructor params\n")
				}
				payload, _, err = GenPayload(abiFile, "", argList)
				if err != nil {
					return err
				}
				code = append(code, payload...)
			}

			msg := types.NewMsgContract(cliCtx.GetFromAddress(), nil, code, coin)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagCodeFile, "", "contract code file path")
	cmd.Flags().String(flagAbiFile, "", "contract abi file path")
	cmd.Flags().String(flagArgs, "", "contract method arg list (e.g. --args='arg1 arg2 arg3')")
	cmd.Flags().String(flagAmount, "0pnch", "amount of coins to send (e.g. 100pnch)")

	cmd.MarkFlagRequired(flagCodeFile)

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

func ContractCallCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "call",
		Short:   "Create and sign a call contract tx",
		Example: `nchcli vm call --from=<user key name> --contract_addr=<contract_addr> --method=<method> --abi_file=<abi_file>  --args='arg1 arg2 arg3' --amount=<amount> `,
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			coin := sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(0))
			amount := viper.GetString(flagAmount)
			if len(amount) > 0 {
				coinInput, err := sdk.ParseCoin(amount)
				if err != nil {
					return err
				}
				coin = coinInput
			}

			abiFile := viper.GetString(flagAbiFile)
			method := viper.GetString(flagMethod)
			argList := viper.GetStringSlice(flagArgs)
			payload, _, err := GenPayload(abiFile, method, argList)
			if err != nil {
				return err
			}

			dump := make([]byte, len(payload)*2)
			hex.Encode(dump, payload)

			contractAddr, err := sdk.AccAddressFromBech32(viper.GetString(flagContractAddr))
			if err != nil {
				return err
			}

			msg := types.NewMsgContract(cliCtx.GetFromAddress(), contractAddr, payload, coin)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagContractAddr, "", "contract bech32 addr")
	cmd.Flags().String(flagAmount, "0pnch", "amount of coins to send (e.g. 1000000pnch)")
	cmd.Flags().String(flagMethod, "", "contract method")
	cmd.Flags().String(flagArgs, "", "contract method arg list")
	cmd.Flags().String(flagAbiFile, "", "contract abi file path")

	cmd.MarkFlagRequired(flagContractAddr)
	cmd.MarkFlagRequired(flagMethod)
	cmd.MarkFlagRequired(flagAbiFile)

	cmd = client.PostCommands(cmd)[0]

	return cmd
}
