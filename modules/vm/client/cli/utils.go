package cli

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func CodeFromFile(codeFile string) ([]byte, error) {
	codeFile, err := filepath.Abs(codeFile)
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
