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
	if len(codeFile) == 0 {
		return nil, errors.New("code_file can not be empty")
	}

	hexcode, err := ioutil.ReadFile(codeFile)
	if err != nil {
		return nil, err
	}

	hexcode = bytes.TrimSpace(hexcode)

	if len(hexcode) == 0 {
		return nil, errors.New("code can not be empty")
	}

	if len(hexcode)%2 != 0 {
		return nil, fmt.Errorf("invalid input length for hex data (%d)", len(hexcode))
	}

	code, err := hex.DecodeString(string(hexcode))
	if err != nil {
		return nil, err
	}

	return code, nil
}

func AbiFromFile(abiFile string) (abiObj abi.ABI, err error) {
	abiFile, err = filepath.Abs(abiFile)
	if err != nil {
		return
	}
	if len(abiFile) == 0 {
		err = errors.New("abi_file can not be empty")
		return
	}

	abiData, err := ioutil.ReadFile(abiFile)
	if err != nil {
		return
	}

	return abi.JSON(strings.NewReader(string(abiData)))
}

var emptyAddrHex = [20]byte{}

func stringAddr2hexAddr(stringAddr string) [20]byte {
	addrStr := stringAddr
	if len(addrStr) <= 2 {
		return emptyAddrHex
	}

	if addrStr[:3] == "nch" {
		addr, err := sdk.AccAddressFromBech32(addrStr)
		if err != nil {
			return emptyAddrHex
		}

		if len(addr) != 20 {
			return emptyAddrHex
		}

		var addrBin [20]byte
		copy(addrBin[:], addr)
		return addrBin
	} else {
		if addrStr[:2] == "0x" {
			addrStr = addrStr[2:]
		}

		if len(addrStr) != 40 {
			return emptyAddrHex
		}

		addrV, err := hexutil.Decode(addrStr)
		if err != nil {
			return emptyAddrHex
		}
		var v = [20]byte{}
		copy(v[:], addrV)
		return v
	}
}

func GenPayload(abiFile, method string, args []string) (payload []byte, m abi.Method, err error) {
	//fmt.Fprintf(os.Stderr, fmt.Sprintf("abiFile = %s, method = %s, args = %v, len=%d\n", abiFile, method, args, len(args)))

	emptyMethod := abi.Method{}
	abiObj, err := AbiFromFile(abiFile)
	if err != nil {
		return nil, emptyMethod, err
	}

	if len(method) == 0 { //constructor
		m = abiObj.Constructor
	} else if v, ok := abiObj.Methods[method]; ok {
		m = v
	} else {
		return nil, emptyMethod, fmt.Errorf("method %s not exist", method)
	}

	var readyArgs []interface{}

	if len(args) != len(m.Inputs) {
		return nil, emptyMethod, fmt.Errorf("args number dismatch,  expected %d args, actual %d args", len(m.Inputs), len(args))
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
				return nil, emptyMethod, fmt.Errorf("wrong address format, actual address[%s]", addrStr)
			}

			if addrStr[:3] == "nch" {
				addr, err := sdk.AccAddressFromBech32(addrStr)
				if err != nil {
					return nil, emptyMethod, err
				}

				if len(addr) != 20 {
					return nil, emptyMethod, fmt.Errorf("wrong bech32 address format, actual address[%s]", addrStr)
				}

				var addrBin [20]byte
				copy(addrBin[:], addr)
				readyArgs = append(readyArgs, addrBin)
			} else {
				if addrStr[:2] == "0x" {
					addrStr = addrStr[2:]
				}

				if len(addrStr) != 40 {
					return nil, emptyMethod, fmt.Errorf("address must have 40 chars except the prefix '0x', actual %d chars", len(addrStr))
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
				return nil, emptyMethod, fmt.Errorf("must have %d chars except the prefix '0x', actual %d chars", m.Inputs[i].Type.Size*2, len(bytes))
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
				return nil, emptyMethod, fmt.Errorf("parse int failed")
			}

			if v.Cmp(typeMinValue) == -1 || v.Cmp(typeMaxValue) == 1 {
				return nil, emptyMethod, fmt.Errorf("value of type[%s] must be in range [%d, %d]", m.Inputs[i].Type.String(), typeMinValue, typeMaxValue)
			}

			unsignedUint64 := v.Uint64()
			switch m.Inputs[i].Type.Size {
			case 8:
				if m.Inputs[i].Type.T == abi.IntTy {
					readyArgs = append(readyArgs, int8(unsignedUint64))
				} else {
					readyArgs = append(readyArgs, uint8(unsignedUint64))
				}
			case 16:
				if m.Inputs[i].Type.T == abi.IntTy {
					readyArgs = append(readyArgs, int16(unsignedUint64))
				} else {
					readyArgs = append(readyArgs, uint16(unsignedUint64))
				}
			case 32:
				if m.Inputs[i].Type.T == abi.IntTy {
					readyArgs = append(readyArgs, int32(unsignedUint64))
				} else {
					readyArgs = append(readyArgs, uint32(unsignedUint64))
				}
			case 64:
				if m.Inputs[i].Type.T == abi.IntTy {
					readyArgs = append(readyArgs, int64(unsignedUint64))
				} else {
					readyArgs = append(readyArgs, unsignedUint64)
				}
			default:
				readyArgs = append(readyArgs, v)
			}

		case abi.SliceTy:
			d := a[1 : len(a)-1]
			rawItems := strings.Split(d, ",")
			var hexAddrs [][20]byte
			for _, item := range rawItems {
				addrHex := stringAddr2hexAddr(item)
				if addrHex == emptyAddrHex {
					return nil, emptyMethod, fmt.Errorf("addr:[%s] invalid", a)
				}
				hexAddrs = append(hexAddrs, addrHex)
			}
			readyArgs = append(readyArgs, hexAddrs)

		default:
			return nil, emptyMethod, fmt.Errorf("no supported type [%s:%d]", m.Inputs[i].Type.String(), m.Inputs[i].Type.T)
		}
	}

	//fmt.Fprintf(os.Stderr, fmt.Sprintf("readyArgs = %v\n", readyArgs))

	payload, err = abiObj.Pack(method, readyArgs...)
	if err != nil {
		return nil, emptyMethod, err
	}

	return payload, m, err
}
