package cli

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/netcloth/netcloth-chain/app/v0/vm/common/math"
	"github.com/netcloth/netcloth-chain/hexutil"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func CodeFromFile(codeFile string) ([]byte, error) {
	codeFile, _ = filepath.Abs(codeFile)
	if 0 == len(codeFile) {
		return nil, errors.New("code_file can not be empty")
	}

	hexcode, err := ioutil.ReadFile(codeFile)
	if err != nil {
		return nil, err
	}

	hexcode = bytes.TrimSpace(hexcode)

	if 0 == len(hexcode) {
		return nil, errors.New("code can not be empty")
	}

	if len(hexcode)%2 != 0 {
		return nil, errors.New(fmt.Sprintf("Invalid input length for hex data (%d)\n", len(hexcode)))
	}

	code, err := hex.DecodeString(string(hexcode))
	if err != nil {
		return nil, err
	}

	return code, nil
}

func AbiFromFile(abiFile string) (abiObj abi.ABI, err error) {
	abiFile, err = filepath.Abs(abiFile)
	if 0 == len(abiFile) {
		err = errors.New("abi_file can not be empty")
		return
	}

	abiData, err := ioutil.ReadFile(abiFile)
	if err != nil {
		return
	}

	return abi.JSON(strings.NewReader(string(abiData)))
}

func GenPayload(abiFile, method string, args []string) (payload []byte, m abi.Method, err error) {
	//fmt.Fprintf(os.Stderr, fmt.Sprintf("abiFile = %s, method = %s, args = %v, len=%d\n", abiFile, method, args, len(args)))

	emptyMethod := abi.Method{}
	abiObj, err := AbiFromFile(abiFile)
	if err != nil {
		return nil, emptyMethod, err
	}

	exist := false
	if len(method) == 0 { //constructor
		m = abiObj.Constructor
	} else {
		m, exist = abiObj.Methods[method]
		if !exist {
			return nil, emptyMethod, errors.New(fmt.Sprintf("method %s not exist\n", method))
		}
	}

	var readyArgs []interface{}

	if len(args) != len(m.Inputs) {
		return nil, emptyMethod, errors.New(fmt.Sprintf("args number dismatch,  expected %d args, actual %d args\n", len(m.Inputs), len(args)))
	}

	for i, a := range args {
		switch m.Inputs[i].Type.T {
		case abi.BoolTy:
			boolV, err := strconv.ParseBool(a)
			if err != nil {
				return nil, emptyMethod, err
			}
			readyArgs = append(readyArgs, boolV)

		case abi.AddressTy:
			addrStr := a
			if len(addrStr) <= 2 {
				return nil, emptyMethod, errors.New(fmt.Sprintf("wrong address format, actual address[%s]", addrStr))
			}

			if addrStr[:3] == "nch" {
				addr, err := sdk.AccAddressFromBech32(addrStr)
				if err != nil {
					return nil, emptyMethod, err
				}

				if len(addr) != 20 {
					return nil, emptyMethod, errors.New(fmt.Sprintf("wrong bech32 address format, actual address[%s]", addrStr))
				}

				var addrBin [20]byte
				copy(addrBin[:], addr)
				readyArgs = append(readyArgs, addrBin)
			} else {
				if addrStr[:2] == "0x" {
					addrStr = addrStr[2:]
				}

				if len(addrStr) != 40 {
					return nil, emptyMethod, errors.New(fmt.Sprintf("address must have 40 chars except the prefix '0x', actual %d chars\n", len(addrStr)))
				}

				addrV, err := hexutil.Decode(addrStr)
				if err != nil {
					return nil, emptyMethod, err
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
				return nil, emptyMethod, err
			}
			readyArgs = append(readyArgs, v)

		case abi.FixedBytesTy:
			bytes := a
			if bytes[:2] == "0x" {
				bytes = bytes[2:]
			}
			if len(bytes) != m.Inputs[i].Type.Size*2 {
				return nil, emptyMethod, errors.New(fmt.Sprintf("must have %d chars except the prefix '0x', actual %d chars\n", m.Inputs[i].Type.Size*2, len(bytes)))
			}
			dv, err := hexutil.Decode(bytes)
			if err != nil {
				return nil, emptyMethod, err
			}

			var byteValue byte
			byteType := reflect.TypeOf(byteValue)
			byteArrayType := reflect.ArrayOf(m.Inputs[i].Type.Size, byteType)
			byteArrayValue := reflect.New(byteArrayType).Elem()
			for j := 0; j < m.Inputs[i].Type.Size; j++ {
				byteArrayValue.Index(j).Set(reflect.ValueOf(dv[j]))
			}

			readyArgs = append(readyArgs, byteArrayValue.Interface())

		case abi.IntTy, abi.UintTy:
			var typeMinValue *big.Int
			var typeMaxValue *big.Int
			if abi.IntTy == m.Inputs[i].Type.T {
				typeMinValue = math.BigPow(2, int64(m.Inputs[i].Type.Size)-1)
				typeMinValue.Neg(typeMinValue)
				typeMaxValue = math.BigPow(2, int64(m.Inputs[i].Type.Size)-1)
				typeMaxValue = typeMaxValue.Sub(typeMaxValue, big.NewInt(1))
			} else {
				typeMinValue = big.NewInt(0)
				typeMaxValue = math.BigPow(2, int64(m.Inputs[i].Type.Size))
			}
			//fmt.Println(fmt.Sprintf("type:%s, bit size:%d, [%d, %d]", m.Inputs[i].Type.String(), m.Inputs[i].Type.Size, typeMinValue, typeMaxValue))

			v, success := big.NewInt(0).SetString(a, 10)
			if !success {
				return nil, emptyMethod, errors.New(fmt.Sprintf("parse int failed"))
			}

			if v.Cmp(typeMinValue) == -1 || v.Cmp(typeMaxValue) == 1 {
				return nil, emptyMethod, errors.New(fmt.Sprintf("value of type[%s] must be in range [%d, %d]", m.Inputs[i].Type.String(), typeMinValue, typeMaxValue))
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
			return nil, emptyMethod, errors.New(fmt.Sprintf("no supported type [%s:%d]", m.Inputs[i].Type.String(), m.Inputs[i].Type.T))
		}
	}

	//fmt.Fprintf(os.Stderr, fmt.Sprintf("readyArgs = %v\n", readyArgs))

	payload, err = abiObj.Pack(method, readyArgs...)
	if err != nil {
		return nil, emptyMethod, err
	}

	return
}
