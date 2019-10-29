package types

import (
    "fmt"
    sdk "github.com/NetCloth/netcloth-chain/types"
    "strings"
)

const (
    maxMonikerLength                = 64
    maxWebsiteLength                = 64
    maxServerEndPointLength         = 64
    maxDetailsLength                = 1024
)

var (
    _ sdk.Msg = MsgServiceNodeClaim{}
)

type MsgServiceNodeClaim struct {
    OperatorAddress sdk.AccAddress  `json:"operator_address" yaml:"operator_address"` // address of the ServiceNode's operator
    Moniker         string          `json:"moniker" yaml:"moniker"`                   // name
    Website         string          `json:"website" yaml:"website"`                   // optional website link
    ServiceType     uint64          `json:"service_type" yaml:"service_type"`
    ServerEndPoint  string          `json:"server_endpoint" yaml:"server_endpoint"`   // server endpoint for app client
    Details         string          `json:"details" yaml:"details"`                   // optional details
    Bond            sdk.Coin        `json:"bond" yaml:"bond"`
}

func NewMsgServiceNodeClaim(operator sdk.AccAddress, moniker, website string, serviceType uint64, serverEndPoint, details string, Bond sdk.Coin) MsgServiceNodeClaim {
    return MsgServiceNodeClaim {
        OperatorAddress:    operator,
        Moniker:            moniker,
        Website:            website,
        ServiceType:        serviceType,
        ServerEndPoint:     serverEndPoint,
        Details:            details,
        Bond:               Bond,
    }
}

func (msg MsgServiceNodeClaim) Route() string { return RouterKey }

func (msg MsgServiceNodeClaim) Type() string { return "serviceNodeClaim" }

func (msg *MsgServiceNodeClaim) TrimSpace() {
    msg.Moniker = strings.TrimSpace(msg.Moniker)
    msg.Website = strings.TrimSpace(msg.Website)
    msg.Details = strings.TrimSpace(msg.Details)
    msg.ServerEndPoint = strings.TrimSpace(msg.ServerEndPoint)
}

func (msg MsgServiceNodeClaim) ValidateBasic() sdk.Error {
    if msg.OperatorAddress.Empty() {
        return sdk.ErrInvalidAddress("missing operator address")
    }

    if msg.Moniker == "" {
        return ErrEmptyInputs("moniker empty")
    }

    if msg.ServerEndPoint == "" {
        return ErrEmptyInputs("server empty")
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

    if len(msg.ServerEndPoint) > maxServerEndPointLength {
        return ErrStringTooLong(fmt.Sprintf("server too long, max length: %v", maxServerEndPointLength))
    }

    if len(msg.Details) > maxDetailsLength {
        return ErrStringTooLong(fmt.Sprintf("details too long, max length: %v", maxDetailsLength))
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
