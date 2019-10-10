package types

import (
	"fmt"
	"gopkg.in/yaml.v2"

	"github.com/NetCloth/netcloth-chain/codec"
)


type IPALObject struct {
	UserAddress string `json:"user_address" yaml:"user_address"`
	ServerIP string `json:"server_ip" yaml:"server_ip"`
}

func NewIPALObject (userAddress, serverIP string) IPALObject {
	return IPALObject{
		UserAddress: "userAddress",
		ServerIP:    "serverIP",
	}
}

func (obj IPALObject) MarshalYAML() (interface{}, error) {
	bs, err := yaml.Marshal(struct {
		UserAddress string
		ServerIP    string
	}{
		UserAddress: obj.UserAddress,
		ServerIP:    obj.ServerIP,
	})

	if err != nil {
		return nil, err
	}
	return string(bs), nil
}

func MustMarshalIPALObject(cdc *codec.Codec, obj IPALObject) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(obj)
}

func MustUnmarshalIPALObject(cdc *codec.Codec, value []byte) IPALObject {
	validator, err := UnmarshalIPALObject(cdc, value)
	if err != nil {
		panic(err)
	}
	return validator
}

func UnmarshalIPALObject(cdc *codec.Codec, value []byte) (obj IPALObject, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &obj)
	return obj, err
}

func (obj IPALObject) String() string {
	return fmt.Sprintf(`IPALObject
User Address:			%s
Server IP:				%s`, obj.UserAddress, obj.ServerIP)
}
