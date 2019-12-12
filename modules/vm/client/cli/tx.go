package cli

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/netcloth/netcloth-chain/client"
	"github.com/netcloth/netcloth-chain/client/context"
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/auth"
	"github.com/netcloth/netcloth-chain/modules/auth/client/utils"
	"github.com/netcloth/netcloth-chain/modules/vm/types"
	sdk "github.com/netcloth/netcloth-chain/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
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
		Example: "nchcli vm create --from=<user key name> --amount=<amount> --code_file=<code file>",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			amount := viper.GetString(flagAmount)
			if "" == amount {
				amount = "0unch"
			}
			coin, err := sdk.ParseCoin(amount)
			if err != nil {
				return err
			}

			codeFile := viper.GetString(flagCodeFile)
			codeFile, err = filepath.Abs(codeFile)
			if 0 == len(codeFile) {
				return errors.New("code_file can not be empty")
			}

			hexcode, err := ioutil.ReadFile(codeFile)
			if err != nil {
				return err
			}

			hexcode = bytes.TrimSpace(hexcode)

			if 0 == len(hexcode) {
				return errors.New("code can not be empty")
			}

			if len(hexcode)%2 != 0 {
				return errors.New(fmt.Sprintf("Invalid input length for hex data (%d)\n", len(hexcode)))
			}

			code, err := hex.DecodeString(string(hexcode))
			if err != nil {
				return err
			}

			msg := types.NewMsgContractCreate(cliCtx.GetFromAddress(), coin, code)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagCodeFile, "", "contract code file")
	cmd.Flags().String(flagAmount, "", "send tokens to contract amount (e.g. 1000000unch)")

	cmd.MarkFlagRequired(flagCodeFile)

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

func ContractCallCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "call",
		Short:   "call a contract",
		Example: "nchcli vm call --from=<user key name> --contract_addr=<contract_addr> --amount=<amount> --abi_file=<abi_file> --method=<method> --args=<args>",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			amount := viper.GetString(flagAmount)
			if "" == amount {
				amount = "0unch"
			}
			coin, err := sdk.ParseCoin(amount)
			if err != nil {
				return err
			}

			abiFile := viper.GetString(flagAbiFile)
			abiFile, err = filepath.Abs(abiFile)
			if 0 == len(abiFile) {
				return errors.New("abi_file can not be empty")
			}

			abiData, err := ioutil.ReadFile(abiFile)
			abiObj, err := abi.JSON(strings.NewReader(string(abiData)))
			if err != nil {
				return err
			}

			method := viper.GetString(flagMethod)
			argsString := viper.GetString(flagArgs)
			argsBinary, err := hex.DecodeString(argsString)
			if err != nil {
				return err
			}

			m, exist := abiObj.Methods[method]
			var payload []byte
			if exist {
				if len(m.Inputs) != len(argsBinary)/32 {
					return errors.New(fmt.Sprint("args count dismatch"))
				}

				readyArgs, err := m.Inputs.UnpackValues(argsBinary)
				if err != nil {
					return err
				}

				payload, err = abiObj.Pack(method, readyArgs...)
				if err != nil {
					return err
				}
			} else {
				return errors.New(fmt.Sprintf("method %s not exist\n", method))
			}

			dump := make([]byte, len(payload)*2)
			hex.Encode(dump[:], payload)
			fmt.Fprintf(os.Stderr, fmt.Sprintf("paylaod = %s\n", string(dump)))

			contractAddr, err := sdk.AccAddressFromBech32(viper.GetString(flagContractAddr))
			if err != nil {
				return err
			}

			msg := types.NewMsgContractCall(cliCtx.GetFromAddress(), contractAddr, coin, payload)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagContractAddr, "", "contract bech32 addr")
	cmd.Flags().String(flagAmount, "", "send tokens to contract amount (e.g. 1000000unch)")
	cmd.Flags().String(flagMethod, "", "contract method")
	cmd.Flags().String(flagArgs, "", "contract method arg list, e.g. [f(a uint, b uint) a=1,b=1] --> 00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001")
	cmd.Flags().String(flagAbiFile, "", "contract abi file")

	cmd.MarkFlagRequired(flagContractAddr)
	cmd.MarkFlagRequired(flagMethod)
	cmd.MarkFlagRequired(flagArgs)

	cmd = client.PostCommands(cmd)[0]

	return cmd
}
