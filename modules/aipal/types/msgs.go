package types

import (
    "fmt"
    sdk "github.com/NetCloth/netcloth-chain/types"
    "strings"
)

const (
    maxMonikerLength   = 64
    maxWebsiteLength   = 64
    maxEndpointsLength = 1024
    maxDetailsLength   = 1024
)

var (
    _ sdk.Msg = MsgServiceNodeClaim{}
)

type Endpoint struct {
    Type        uint64      `json:"type" yaml:"type"`
    Endpoint    string      `json:"endpoint" yaml:"endpoint"`
}

func NewEndpoint(endpointType uint64, endpoint string) Endpoint {
    return Endpoint {
        Type: endpointType,
        Endpoint: endpoint,
    }
}

type Endpoints []Endpoint

func (e Endpoints) String() string {
    return fmt.Sprintf("%v", e)
}

type MsgServiceNodeClaim struct {
    OperatorAddress sdk.AccAddress  `json:"operator_address" yaml:"operator_address"`   // address of the ServiceNode's operator
    Moniker         string          `json:"moniker" yaml:"moniker"`                     // name
    Website         string          `json:"website" yaml:"website"`                     // optional website link
    Details         string          `json:"details" yaml:"details"`                     // optional details
    Endpoints       Endpoints       `json:"endpoints" yaml:"endpoints"`                 // server endpoint for app client
    Bond            sdk.Coin        `json:"bond" yaml:"bond"`
}

func NewMsgServiceNodeClaim(operator sdk.AccAddress, moniker, website, details string, endpoints Endpoints, Bond sdk.Coin) MsgServiceNodeClaim {
    return MsgServiceNodeClaim {
        OperatorAddress:    operator,
        Moniker:            moniker,
        Website:            website,
        Details:            details,
        Endpoints:          endpoints,
        Bond:               Bond,
    }
}

func (msg MsgServiceNodeClaim) Route() string { return RouterKey }

func (msg MsgServiceNodeClaim) Type() string { return "serviceNodeClaim" }

func (msg *MsgServiceNodeClaim) TrimSpace() {
    msg.Moniker = strings.TrimSpace(msg.Moniker)
    msg.Website = strings.TrimSpace(msg.Website)
    msg.Details = strings.TrimSpace(msg.Details)
}

func (msg MsgServiceNodeClaim) ValidateBasic() sdk.Error {
    if msg.OperatorAddress.Empty() {
        return sdk.ErrInvalidAddress("missing operator address")
    }

    if msg.Moniker == "" {
        return ErrEmptyInputs("moniker empty")
    }

    if msg.Bond.Denom != sdk.NativeTokenName {
       return ErrBadDenom(fmt.Sprintf("bond denom must be %s", sdk.NativeTokenName))
    }

    if msg.Bond.IsNegative() {
        return ErrEmptyInputs("stake amount must > 0 ")
    }

    if len(msg.Moniker) > maxMonikerLength {
        return ErrStringTooLong(fmt.Sprintf("moniker too long, max length: %v", maxMonikerLength))
    }

    if len(msg.Website) > maxWebsiteLength {
        return ErrStringTooLong(fmt.Sprintf("website too long, max length: %v", maxWebsiteLength))
    }

    if len(msg.Details) > maxDetailsLength {
        return ErrStringTooLong(fmt.Sprintf("details too long, max length: %v", maxDetailsLength))
    }

    if len(msg.Endpoints) > maxEndpointsLength {
        return ErrStringTooLong(fmt.Sprintf("endpoints too long, max length: %v", maxEndpointsLength))
    }

    if len(msg.Endpoints) <= 0 {
        return ErrEndpointsEmpty()
    }

    return nil
}

func (msg MsgServiceNodeClaim) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.OperatorAddress}
}

func (msg MsgServiceNodeClaim) GetSignBytes() []byte {
    bz := ModuleCdc.MustMarshalJSON(msg)
    return sdk.MustSortJSON(bz)
}
