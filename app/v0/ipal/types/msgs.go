package types

import (
	"strings"

	"gopkg.in/yaml.v2"

	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

var (
	_ sdk.Msg = MsgIPALNodeClaim{}
)

const TypeMsgIPALNodeClaim = "ipalNodeClaim"

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
	out, _ := yaml.Marshal(e)
	return string(out)
}

type Endpoints []Endpoint

func (e Endpoints) String() (out string) {
	for _, val := range e {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

type MsgIPALNodeClaim struct {
	OperatorAddress sdk.AccAddress `json:"operator_address" yaml:"operator_address"` // address of the IPALNode's operator
	Moniker         string         `json:"moniker" yaml:"moniker"`                   // name
	Website         string         `json:"website" yaml:"website"`                   // optional website link
	Details         string         `json:"details" yaml:"details"`                   // optional details
	Extension       string         `json:"extension" yaml:"extension"`               // for future extension
	Endpoints       Endpoints      `json:"endpoints" yaml:"endpoints"`               // server endpoint for app client
	Bond            sdk.Coin       `json:"bond" yaml:"bond"`                         // bond coin for ranking
}

func NewMsgIPALNodeClaim(operator sdk.AccAddress, moniker, website, details, extension string, endpoints Endpoints, bond sdk.Coin) MsgIPALNodeClaim {
	return MsgIPALNodeClaim{
		OperatorAddress: operator,
		Moniker:         moniker,
		Website:         website,
		Details:         details,
		Extension:       extension,
		Endpoints:       endpoints,
		Bond:            bond,
	}
}

func (msg MsgIPALNodeClaim) Route() string { return RouterKey }

func (msg MsgIPALNodeClaim) Type() string { return TypeMsgIPALNodeClaim }

func (msg *MsgIPALNodeClaim) TrimSpace() {
	msg.Moniker = strings.TrimSpace(msg.Moniker)
	msg.Website = strings.TrimSpace(msg.Website)
	msg.Details = strings.TrimSpace(msg.Details)
	msg.Extension = strings.TrimSpace(msg.Extension)
}

func (msg MsgIPALNodeClaim) ValidateBasic() error {
	if msg.OperatorAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing operator address")
	}

	if msg.Moniker == "" {
		return sdkerrors.Wrap(ErrEmptyInputs, "moniker empty")
	}

	if msg.Bond.Denom != sdk.NativeTokenName {
		return sdkerrors.Wrapf(ErrBadDenom, "bond denom must be %s", sdk.NativeTokenName)
	}

	if msg.Bond.IsNegative() {
		return sdkerrors.Wrap(ErrEmptyInputs, "bond amount must > 0 ")
	}

	if len(msg.Endpoints) == 0 {
		return ErrEndpointsEmpty
	}

	if err := EndpointsDupCheck(msg.Endpoints); err != nil {
		return err
	}

	return nil
}

func (msg MsgIPALNodeClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OperatorAddress}
}

func (msg MsgIPALNodeClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
