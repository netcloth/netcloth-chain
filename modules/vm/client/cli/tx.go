package cli

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/netcloth/netcloth-chain/client"
	"github.com/netcloth/netcloth-chain/client/context"
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/auth"
	"github.com/netcloth/netcloth-chain/modules/auth/client/utils"
	"github.com/netcloth/netcloth-chain/modules/vm/types"
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
		Short:   "Create a contract",
		Example: "nchcli vm create --from=<user key name> --amount=<amount> --code_file=<code file> --args=<args> --abi_file=<abi_file>",
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
					return errors.New(fmt.Sprintf("must use --abi_file to appoint abi file when use constructor params\n"))
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

	cmd.Flags().String(flagCodeFile, "", "contract code file")
	cmd.Flags().String(flagAbiFile, "", "contract abi file")
	cmd.Flags().String(flagArgs, "", "contract construct function arg list, e.g. [constructor(a uint, b uint) a=1,b=1] --> 00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001")
	cmd.Flags().String(flagAmount, "", "send tokens to contract amount (e.g. 1000000pnch)")

	cmd.MarkFlagRequired(flagCodeFile)

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

func ContractCallCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "call",
		Short:   "call a contract",
		Example: `nchcli vm call --from=<user key name> --contract_addr=<contract_addr> --amount=<amount> --abi_file=<abi_file> --method=<method> --args=''`,
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
				return nil
			}

			dump := make([]byte, len(payload)*2)
			hex.Encode(dump[:], payload)
			fmt.Fprintf(os.Stderr, fmt.Sprintf("paylaod = %s\n", string(dump)))

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
	cmd.Flags().String(flagAmount, "", "send tokens to contract amount (e.g. 1000000pnch)")
	cmd.Flags().String(flagMethod, "", "contract method")
	cmd.Flags().String(flagArgs, "", "contract method arg list")
	cmd.Flags().String(flagAbiFile, "", "contract abi file")

	cmd.MarkFlagRequired(flagContractAddr)
	cmd.MarkFlagRequired(flagMethod)
	cmd.MarkFlagRequired(flagArgs)

	cmd = client.PostCommands(cmd)[0]

	return cmd
}
