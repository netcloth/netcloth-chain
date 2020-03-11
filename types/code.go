package types

import (
	"encoding/json"
	"fmt"
	"github.com/netcloth/netcloth-chain/hexutil"
)

type (
	Code    []byte
	Storage map[Hash]Hash
)

func (c Code) String() string {
	return hexutil.Encode(c)
}

func (c Code) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Code) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	d, err := hexutil.Decode(s)

	if err != nil {
		return err
	}

	*c = d

	return nil
}

func (s Storage) String() (str string) {
	for key, value := range s {
		str += fmt.Sprintf("%X : %X\n", key, value)
	}
	return
}

func (s Storage) Copy() Storage {
	cpy := make(Storage)
	for key, value := range s {
		cpy[key] = value
	}

	return cpy
}
