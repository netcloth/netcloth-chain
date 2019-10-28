package types

import (
    "fmt"
    "gopkg.in/yaml.v2"
    "strings"

    "github.com/NetCloth/netcloth-chain/codec"
    sdk "github.com/NetCloth/netcloth-chain/types"
)

type ServiceType uint64

const (
    Chatting ServiceType = 1 << iota
    Storage
)

func ServiceTypeFromString(s string) ServiceType {
    var v ServiceType

    s = strings.Replace(s, " ", "", -1)
    types := strings.Split(s, "|")
    for _, t := range types {
        if t == "chatting" {
            v |= Chatting
        }

        if t == "storage" {
            v |= Storage
        }
    }
    return v
}

type ServiceNode struct {
    OperatorAddress sdk.AccAddress  `json:"operator_address" yaml:"operator_address"` // address of the ServiceNode's operator
    Moniker         string          `json:"moniker" yaml:"moniker"`                   // name
    Website         string          `json:"website" yaml:"website"`                   // optional website link
    ServiceType     ServiceType     `json:"service_type" yaml:"service_type"`
    ServerEndPoint  string          `json:"server_endpoint" yaml:"server_endpoint"`   // server endpoint for app client
    Details         string          `json:"details" yaml:"details"`                   // optional details
    Bond            sdk.Coin        `json:"bond" yaml:"bond"`
}

type ServiceNodes []ServiceNode

func (v ServiceNodes) String() (out string) {
    for _, val := range v {
        out += val.String() + "\n"
    }
    return strings.TrimSpace(out)
}

func NewServiceNode(operator sdk.AccAddress, moniker, website string, serviceType ServiceType, serverEndPoint, details string, amount sdk.Coin) ServiceNode {
    return ServiceNode {
        OperatorAddress:    operator,
        Moniker:            moniker,
        Website:            website,
        ServiceType:        serviceType,
        ServerEndPoint:     serverEndPoint,
        Details:            details,
        Bond:               amount,
    }
}

func (obj ServiceNode) MarshalYAML() (interface{}, error) {
    bs, err := yaml.Marshal(struct {
        OperatorAddress sdk.AccAddress
        Moniker         string
        Website         string
        ServerEndPoint  string
        Details         string
        Bond     sdk.Coin
    }{
        OperatorAddress:    obj.OperatorAddress,
        Moniker:            obj.Moniker,
        Website:            obj.Website,
        ServerEndPoint:     obj.ServerEndPoint,
        Details:            obj.Details,
        Bond:               obj.Bond,
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
Operator Address:		%s
Moniker:				%s
Website: 				%s
ServerEndPoint:			%s
Details:				%s
Bond: 			        %s`, obj.OperatorAddress, obj.Moniker, obj.Website, obj.ServerEndPoint, obj.Details, obj.Bond)
}