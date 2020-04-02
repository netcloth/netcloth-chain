package types

import (
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

type ServiceType uint64

const (
	Chatting ServiceType = 1
	Storage
)

func EndpointsFromString(endpointsStr, endpointDelimiter, endpointTypeDelimiter string) (endpoints Endpoints, e error) {
	endpointStrings := strings.Split(endpointsStr, endpointDelimiter)
	for _, endpointString := range endpointStrings {
		endpointString = strings.Trim(endpointString, " \t")

		if len(endpointString) > 0 {
			typeAndEndpoint := strings.Split(endpointString, endpointTypeDelimiter)
			if len(typeAndEndpoint) != 2 {
				return nil, sdkerrors.Wrapf(ErrEndpointsFormat, "should be in format: serviceType%sendpoint%sserviceType%sendpoint", endpointTypeDelimiter, endpointDelimiter, endpointTypeDelimiter)
			}

			typeAndEndpoint[0] = strings.Trim(typeAndEndpoint[0], " \t")
			typeAndEndpoint[1] = strings.Trim(typeAndEndpoint[1], " \t")

			if len(typeAndEndpoint[0]) == 0 || len(typeAndEndpoint[1]) == 0 {
				return nil, sdkerrors.Wrapf(ErrEndpointsFormat, "should be in format: serviceType%sendpoint%sserviceType%sendpoint", endpointTypeDelimiter, endpointDelimiter, endpointTypeDelimiter)
			}

			endpointType, err := strconv.Atoi(typeAndEndpoint[0])
			if err != nil {
				return nil, sdkerrors.Wrapf(ErrEndpointsFormat, "should be in format: serviceType%sendpoint%sserviceType%sendpoint", endpointTypeDelimiter, endpointDelimiter, endpointTypeDelimiter)
			}

			endpoints = append(endpoints, NewEndpoint(uint64(endpointType), typeAndEndpoint[1]))
		} else {
			return nil, sdkerrors.Wrapf(ErrEndpointsFormat, "should be in format: serviceType%sendpoint%sserviceType%sendpoint", endpointTypeDelimiter, endpointDelimiter, endpointTypeDelimiter)
		}
	}

	if err := EndpointsDupCheck(endpoints); err != nil {
		return nil, err
	}

	return endpoints, nil
}

type IPALNode struct {
	OperatorAddress sdk.AccAddress `json:"operator_address" yaml:"operator_address"` // address of the IPALNode's operator
	Moniker         string         `json:"moniker" yaml:"moniker"`                   // name
	Website         string         `json:"website" yaml:"website"`                   // optional website link
	Details         string         `json:"details" yaml:"details"`                   // optional details
	Extension       string         `json:"extension" yaml:"extension"`
	Endpoints       Endpoints      `json:"endpoints" yaml:"endpoints"`
	Bond            sdk.Coin       `json:"bond" yaml:"bond"`
}

type IPALNodes []IPALNode

func (v IPALNodes) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

func NewIPALNode(operator sdk.AccAddress, moniker, website, details, extension string, endpoints Endpoints, amount sdk.Coin) IPALNode {
	return IPALNode{
		OperatorAddress: operator,
		Moniker:         moniker,
		Website:         website,
		Details:         details,
		Extension:       extension,
		Endpoints:       endpoints,
		Bond:            amount,
	}
}

func (obj IPALNode) MarshalYAML() (interface{}, error) {
	bs, err := yaml.Marshal(struct {
		OperatorAddress sdk.AccAddress
		Moniker         string
		Website         string
		Endpoints       Endpoints
		Details         string
		Extension       string
		Bond            sdk.Coin
	}{
		OperatorAddress: obj.OperatorAddress,
		Moniker:         obj.Moniker,
		Website:         obj.Website,
		Endpoints:       obj.Endpoints,
		Details:         obj.Details,
		Extension:       obj.Extension,
		Bond:            obj.Bond,
	})

	if err != nil {
		return nil, err
	}

	return string(bs), nil
}

func MustMarshalIPALNode(cdc *codec.Codec, obj IPALNode) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(obj)
}

func MustUnmarshalIPALNode(cdc *codec.Codec, value []byte) IPALNode {
	obj, err := UnmarshalIPALNode(cdc, value)
	if err != nil {
		panic(err)
	}
	return obj
}

func UnmarshalIPALNode(cdc *codec.Codec, value []byte) (obj IPALNode, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &obj)
	return obj, err
}

func (obj IPALNode) String() string {
	out, _ := yaml.Marshal(obj)
	return string(out)
}
