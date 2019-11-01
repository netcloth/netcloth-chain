package types

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/NetCloth/netcloth-chain/codec"
	sdk "github.com/NetCloth/netcloth-chain/types"
)

type ServiceType uint64

const (
	Chatting ServiceType = 1
	Storage
)

func EndpointsFromString(s string) (r Endpoints, e sdk.Error) {
	ss := strings.Split(s, ",")
	for _, v := range ss {
		v = strings.ReplaceAll(v, " ", "")
		if len(v) > 0 {
			es := strings.Split(v, "|")
			if len(es) != 2 {
				return nil, ErrEndpointsFormat()
			}

			if len(es[0]) == 0 || len(es[1]) == 0 {
				return nil, ErrEndpointsFormat()
			}

			Type, err := strconv.Atoi(es[0])
			if err != nil {
				return nil, ErrEndpointsFormat()
			}

			r = append(r, NewEndpoint(uint64(Type), es[1]))
		} else {
			return nil, ErrEndpointsFormat()
		}
	}

	return r, nil
}

type ServiceNode struct {
	OperatorAddress sdk.AccAddress `json:"operator_address" yaml:"operator_address"` // address of the ServiceNode's operator
	Moniker         string         `json:"moniker" yaml:"moniker"`                   // name
	Website         string         `json:"website" yaml:"website"`                   // optional website link
	Details         string         `json:"details" yaml:"details"`                   // optional details
	Endpoints       Endpoints      `json:"endpoints" yaml:"endpoints"`
	Bond            sdk.Coin       `json:"bond" yaml:"bond"`
}

type ServiceNodes []ServiceNode

func (v ServiceNodes) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

func NewServiceNode(operator sdk.AccAddress, moniker, website string, details string, endpoints Endpoints, amount sdk.Coin) ServiceNode {
	return ServiceNode{
		OperatorAddress: operator,
		Moniker:         moniker,
		Website:         website,
		Details:         details,
		Endpoints:       endpoints,
		Bond:            amount,
	}
}

func (obj ServiceNode) MarshalYAML() (interface{}, error) {
	bs, err := yaml.Marshal(struct {
		OperatorAddress sdk.AccAddress
		Moniker         string
		Website         string
		Endpoints       Endpoints
		Details         string
		Bond            sdk.Coin
	}{
		OperatorAddress: obj.OperatorAddress,
		Moniker:         obj.Moniker,
		Website:         obj.Website,
		Endpoints:       obj.Endpoints,
		Details:         obj.Details,
		Bond:            obj.Bond,
	})

	if err != nil {
		return nil, err
	}

	return string(bs), nil
}

func MustMarshalServerNodeObject(cdc *codec.Codec, obj ServiceNode) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(obj)
}

func MustUnmarshalServerNodeObject(cdc *codec.Codec, value []byte) ServiceNode {
	obj, err := UnmarshalServerNodeObject(cdc, value)
	if err != nil {
		panic(err)
	}
	return obj
}

func UnmarshalServerNodeObject(cdc *codec.Codec, value []byte) (obj ServiceNode, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &obj)
	return obj, err
}

func (obj ServiceNode) String() string {
	return fmt.Sprintf(`ServerNodeObject
Operator Address:       %s
Moniker:                %s
Website:                %s
Endpoints:              %s
Details:                %s
Bond:                   %s`, obj.OperatorAddress, obj.Moniker, obj.Website, obj.Endpoints.String(), obj.Details, obj.Bond)
}
