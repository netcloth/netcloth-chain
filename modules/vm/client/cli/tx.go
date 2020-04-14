package cli

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/netcloth/netcloth-chain/client"
	"github.com/netcloth/netcloth-chain/client/context"
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/hexutil"
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
		ContractCallCmd2(cdc),
	)
	return txCmd
}

func ContractCreateCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a contract",
		Example: "nchcli vm create --from=<user key name> --amount=<amount> --code_file=<code file> --args=<args>",
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

			argsString := viper.GetString(flagArgs)
			if len(argsString) != 0 {
				argsBin, err := hex.DecodeString(argsString)
				if err != nil {
					return err
				}
				code = append(code, argsBin...)
			}

			msg := types.NewMsgContract(cliCtx.GetFromAddress(), nil, code, coin)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagCodeFile, "", "contract code file")
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
		Example: "nchcli vm call --from=<user key name> --contract_addr=<contract_addr> --amount=<amount> --abi_file=<abi_file> --method=<method> --args=<args>",
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
			abiObj, err := AbiFromFile(abiFile)
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
				//if len(m.Inputs) != len(argsBinary)/32 {
				//	return errors.New(fmt.Sprint("args count dismatch"))
				//}

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
	cmd.Flags().String(flagArgs, "", "contract method arg list, e.g. [f(a uint, b uint) a=1,b=1] --> 00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001")
	cmd.Flags().String(flagAbiFile, "", "contract abi file")

	cmd.MarkFlagRequired(flagContractAddr)
	cmd.MarkFlagRequired(flagMethod)
	cmd.MarkFlagRequired(flagArgs)

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

func ContractCallCmd2(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "call2",
		Short:   "call a contract",
		Example: `nchcli vm call2 --from=<user key name> --contract_addr=<contract_addr> --amount=<amount> --abi_file=<abi_file> --method=<method> --args=''`,
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
			abiObj, err := AbiFromFile(abiFile)
			if err != nil {
				return err
			}

			method := viper.GetString(flagMethod)
			m, exist := abiObj.Methods[method]
			if !exist {
				return errors.New(fmt.Sprintf("method %s not exist\n", method))
			}

			argList := viper.GetStringSlice(flagArgList)
			fmt.Fprintf(os.Stderr, fmt.Sprintf("argList = %v, len=%d\n", argList, len(argList)))
			var readyArgs []interface{}

			if len(argList) != len(m.Inputs) {
				return errors.New(fmt.Sprintf("args number dismatch expected %d args, actual %d args\n", len(m.Inputs), len(argList)))
			}

			for i, a := range argList {
				switch m.Inputs[i].Type.T {
				case abi.BoolTy:
					boolV, err := strconv.ParseBool(a)
					if err != nil {
						return err
					}
					readyArgs = append(readyArgs, boolV)

				case abi.AddressTy:
					addrStr := a
					if len(addrStr) <= 2 {
						return errors.New(fmt.Sprintf("wrong address format, actual address[%s]", addrStr))
					}

					if addrStr[:3] == "nch" {
						addr, err := sdk.AccAddressFromBech32(addrStr)
						if err != nil {
							return err
						}

						if len(addr) != 20 {
							return errors.New(fmt.Sprintf("wrong bech32 address format, actual address[%s]", addrStr))
						}

						var addrBin [20]byte
						copy(addrBin[:], addr)
						readyArgs = append(readyArgs, addrBin)
					} else {
						if addrStr[:2] == "0x" {
							addrStr = addrStr[2:]
						}

						if len(addrStr) != 40 {
							return errors.New(fmt.Sprintf("address must have 40 chars except the prefix '0x', actual %d chars\n", len(addrStr)))
						}

						addrV, err := hexutil.Decode(addrStr)
						if err != nil {
							return err
						}
						var v = [20]byte{}
						copy(v[:], addrV)
						readyArgs = append(readyArgs, v)
					}

				case abi.StringTy:
					readyArgs = append(readyArgs, a)

				case abi.BytesTy:
					v, err := hexutil.Decode(a)
					if err != nil {
						return err
					}
					readyArgs = append(readyArgs, v)

				case abi.FixedBytesTy:
					bytes := a
					if bytes[:2] == "0x" {
						bytes = bytes[2:]
					}
					if len(bytes) != m.Inputs[i].Type.Size*2 {
						return errors.New(fmt.Sprintf("must have %d chars except the prefix '0x', actual %d chars\n", m.Inputs[i].Type.Size*2, len(bytes)))
					}
					dv, err := hexutil.Decode(bytes)
					if err != nil {
						return err
					}

					var byteValue byte
					byteType := reflect.TypeOf(byteValue)
					byteArrayType := reflect.ArrayOf(m.Inputs[i].Type.Size, byteType)
					byteArrayValue := reflect.New(byteArrayType).Elem()
					for j := 0; j < m.Inputs[i].Type.Size; j++ {
						byteArrayValue.Index(j).Set(reflect.ValueOf(dv[j]))
					}

					readyArgs = append(readyArgs, byteArrayValue.Interface())

				case abi.IntTy, abi.UintTy: //TODO add overflow check
					v, success := big.NewInt(0).SetString(a, 10)
					if !success {
						return errors.New(fmt.Sprintf("parse int failed"))
					}

					unsignedUint64 := v.Uint64()
					if m.Inputs[i].Type.Size == 8 {
						if m.Inputs[i].Type.T == abi.IntTy {
							readyArgs = append(readyArgs, int8(unsignedUint64))
						} else {
							readyArgs = append(readyArgs, uint8(unsignedUint64))
						}
					} else if m.Inputs[i].Type.Size == 16 {
						if m.Inputs[i].Type.T == abi.IntTy {
							readyArgs = append(readyArgs, int16(unsignedUint64))
						} else {
							readyArgs = append(readyArgs, uint16(unsignedUint64))
						}
					} else if m.Inputs[i].Type.Size == 32 {
						if m.Inputs[i].Type.T == abi.IntTy {
							readyArgs = append(readyArgs, int32(unsignedUint64))
						} else {
							readyArgs = append(readyArgs, uint32(unsignedUint64))
						}
					} else if m.Inputs[i].Type.Size == 64 {
						if m.Inputs[i].Type.T == abi.IntTy {
							readyArgs = append(readyArgs, int64(unsignedUint64))
						} else {
							readyArgs = append(readyArgs, uint64(unsignedUint64))
						}
					} else {
						readyArgs = append(readyArgs, v)
					}

				default:
					return errors.New(fmt.Sprintf("no supported type [%s:%d]", m.Inputs[i].Type.String(), m.Inputs[i].Type.T))
				}
			}

			fmt.Fprintf(os.Stderr, fmt.Sprintf("readyArgs = %v\n", readyArgs))

			var payload []byte
			payload, err = abiObj.Pack(method, readyArgs...)
			if err != nil {
				return err
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
	cmd.Flags().String(flagArgList, "", "contract method arg list")
	cmd.Flags().String(flagAbiFile, "", "contract abi file")

	cmd.MarkFlagRequired(flagContractAddr)
	cmd.MarkFlagRequired(flagMethod)
	cmd.MarkFlagRequired(flagArgList)

	cmd = client.PostCommands(cmd)[0]

	return cmd
}
