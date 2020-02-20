package types

import (
	"fmt"
	"strings"

	sdk "github.com/netcloth/netcloth-chain/types"
)

var (
	_ sdk.Msg = MsgServiceNodeClaim{}
)

type Endpoint struct {
	Type     uint64 `json:"type" yaml:"type"`
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

func NewEndpoint(endpointType uint64, endpoint string) Endpoint {
	return Endpoint{
		Type:     endpointType,
		Endpoint: endpoint,
	}
}
func (e Endpoint) String() string {
	return fmt.Sprintf(`Endpoint
type: %d
endpoint: %s`, e.Type, e.Endpoint)
}

type Endpoints []Endpoint

func (e Endpoints) String() (out string) {
	for _, val := range e {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

type MsgServiceNodeClaim struct {
	OperatorAddress sdk.AccAddress `json:"operator_address" yaml:"operator_address"` // address of the ServiceNode's operator
	Moniker         string         `json:"moniker" yaml:"moniker"`                   // name
	Website         string         `json:"website" yaml:"website"`                   // optional website link
	Details         string         `json:"details" yaml:"details"`                   // optional details
	Endpoints       Endpoints      `json:"endpoints" yaml:"endpoints"`               // server endpoint for app client
	Bond            sdk.Coin       `json:"bond" yaml:"bond"`
}

func NewMsgServiceNodeClaim(operator sdk.AccAddress, moniker, website, details string, endpoints Endpoints, Bond sdk.Coin) MsgServiceNodeClaim {
	return MsgServiceNodeClaim{
		OperatorAddress: operator,
		Moniker:         moniker,
		Website:         website,
		Details:         details,
		Endpoints:       endpoints,
		Bond:            Bond,
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
		return ErrEmptyInputs("bond amount must > 0 ")
	}

	if len(msg.Endpoints) <= 0 {
		return ErrEndpointsEmpty()
	}

	if err := EndpointsDupCheck(msg.Endpoints); err != nil {
		return err
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
