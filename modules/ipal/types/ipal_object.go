package types

import (
	"fmt"
	"github.com/NetCloth/netcloth-chain/codec"
	sdk "github.com/NetCloth/netcloth-chain/types"
	"gopkg.in/yaml.v2"
)

type IPALObject struct {
	UserAddress string `json:"user_address" yaml:"user_address"`
	ServerIP    string `json:"server_ip" yaml:"server_ip"`
}

type ServerNodeObject struct {
	OperatorAddress sdk.AccAddress `json:"operator_address" yaml:"operator_address"` // address of the ServiceNode's operator
	Moniker         string         `json:"moniker" yaml:"moniker"`                   // name
	Identity        string         `json:"identity" yaml:"identity"`                 // optional identity signature (ex. UPort or Keybase)
	Website         string         `json:"website" yaml:"website"`                   // optional website link
	ServerEndPoint  string         `json:"server_endpoint" yaml:"server_endpoint"`   // server endpoint for app client
	Details         string         `json:"details" yaml:"details"`                   // optional details	DelegatorShares sdk.Dec        `json:"delegator_shares" yaml:"delegator_shares"` // total shares issued to a validator's delegators
}

func NewIPALObject(userAddress string, serverIP string) IPALObject {
	return IPALObject{
		userAddress,
		serverIP,
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

func NewServerNodeObject(operator sdk.AccAddress, moniker, identity, website, serverEndPoint, details string) ServerNodeObject {
	return ServerNodeObject{
		OperatorAddress: operator,
		Moniker:         moniker,
		Identity:        identity,
		Website:         website,
		ServerEndPoint:  serverEndPoint,
		Details:         details,
	}
}

func (obj ServerNodeObject) MarshalYAML() (interface{}, error) {
	bs, err := yaml.Marshal(struct {
		OperatorAddress sdk.AccAddress
		Moniker         string
		Identity        string
		Website         string
		ServerEndPoint  string
		Details         string
	}{
		OperatorAddress: obj.OperatorAddress,
		Moniker:         obj.Moniker,
		Identity:        obj.Identity,
		Website:         obj.Website,
		ServerEndPoint:  obj.ServerEndPoint,
		Details:         obj.Details,
	})

	if err != nil {
		return nil, err
	}
	return string(bs), nil
}

func MustMarshalServerNodeObject(cdc *codec.Codec, obj ServerNodeObject) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(obj)
}

func MustUnmarshalServerNodeObject(cdc *codec.Codec, value []byte) ServerNodeObject {
	obj, err := UnmarshalServerNodeObject(cdc, value)
	if err != nil {
		panic(err)
	}
	return obj
}

func UnmarshalServerNodeObject(cdc *codec.Codec, value []byte) (obj ServerNodeObject, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &obj)
	return obj, err
}

func (obj ServerNodeObject) String() string {
	return fmt.Sprintf(`ServerNodeObject
Operator Address:		%s
Moniker:				%s
Identity: 				%s
Website: 				%s
ServerEndPoint:			%s
Details:				%s`,
		obj.OperatorAddress, obj.Moniker, obj.Identity, obj.Website, obj.ServerEndPoint, obj.Details)
}
