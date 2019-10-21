package types

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/NetCloth/netcloth-chain/codec"
	sdk "github.com/NetCloth/netcloth-chain/types"
)

// IPALObject stores the binding of UserAddress and ServerNode ip
type IPALObject struct {
	UserAddress string `json:"user_address" yaml:"user_address"`
	ServerIP    string `json:"server_ip" yaml:"server_ip"`
}

// IPALObject store info of ServerNode
type ServerNodeObject struct {
	OperatorAddress sdk.AccAddress `json:"operator_address" yaml:"operator_address"` // address of the ServiceNode's operator
	Moniker         string         `json:"moniker" yaml:"moniker"`                   // name
	Website         string         `json:"website" yaml:"website"`                   // optional website link
	ServerEndPoint  string         `json:"server_endpoint" yaml:"server_endpoint"`   // server endpoint for app client
	Details         string         `json:"details" yaml:"details"`                   // optional details
	StakeShares     sdk.Coin       `json:"stake_shares" yaml:"stake_shares"`         // total stake shares
}

// ServerNodeObjects is a collection of ServerNodeObject
type ServerNodeObjects []ServerNodeObject

func (v ServerNodeObjects) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
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

func NewServerNodeObject(operator sdk.AccAddress, moniker, website, serverEndPoint, details string, amount sdk.Coin) ServerNodeObject {
	return ServerNodeObject{
		OperatorAddress: operator,
		Moniker:         moniker,
		Website:         website,
		ServerEndPoint:  serverEndPoint,
		Details:         details,
		StakeShares:     amount,
	}
}

func (obj ServerNodeObject) MarshalYAML() (interface{}, error) {
	bs, err := yaml.Marshal(struct {
		OperatorAddress sdk.AccAddress
		Moniker         string
		Website         string
		ServerEndPoint  string
		Details         string
		StakeShares     sdk.Coin
	}{
		OperatorAddress: obj.OperatorAddress,
		Moniker:         obj.Moniker,
		Website:         obj.Website,
		ServerEndPoint:  obj.ServerEndPoint,
		Details:         obj.Details,
		StakeShares:     obj.StakeShares,
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
Website: 				%s
ServerEndPoint:			%s
Details:				%s
StakeShares: 			%s`,
		obj.OperatorAddress, obj.Moniker, obj.Website, obj.ServerEndPoint, obj.Details, obj.StakeShares)
}
