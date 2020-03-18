package types

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
)

type ServiceType uint64

const (
	Chatting ServiceType = 1
	Storage
)

func EndpointsFromString(endpointsStr, endpointDelimiter, endpointTypeDelimiter string) (endpoints Endpoints, e error) {
	endpointsFormatErr := ErrEndpointsFormat(endpointDelimiter, endpointTypeDelimiter)

	endpointStrings := strings.Split(endpointsStr, endpointDelimiter)
	for _, endpointString := range endpointStrings {
		endpointString = strings.Trim(endpointString, " \t")

		if len(endpointString) > 0 {
			typeAndEndpoint := strings.Split(endpointString, endpointTypeDelimiter)
			if len(typeAndEndpoint) != 2 {
				return nil, endpointsFormatErr
			}

			typeAndEndpoint[0] = strings.Trim(typeAndEndpoint[0], " \t")
			typeAndEndpoint[1] = strings.Trim(typeAndEndpoint[1], " \t")

			if len(typeAndEndpoint[0]) == 0 || len(typeAndEndpoint[1]) == 0 {
				return nil, endpointsFormatErr
			}

			endpointType, err := strconv.Atoi(typeAndEndpoint[0])
			if err != nil {
				return nil, endpointsFormatErr
			}

			endpoints = append(endpoints, NewEndpoint(uint64(endpointType), typeAndEndpoint[1]))
		} else {
			return nil, endpointsFormatErr
		}
	}

	if err := EndpointsDupCheck(endpoints); err != nil {
		return nil, err
	}

	return endpoints, nil
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

func MustMarshalServiceNode(cdc *codec.Codec, obj ServiceNode) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(obj)
}

func MustUnmarshalServiceNode(cdc *codec.Codec, value []byte) ServiceNode {
	obj, err := UnmarshalServiceNode(cdc, value)
	if err != nil {
		panic(err)
	}
	return obj
}

func UnmarshalServiceNode(cdc *codec.Codec, value []byte) (obj ServiceNode, err error) {
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
