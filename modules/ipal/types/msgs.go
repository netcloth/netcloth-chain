package types

import (
	"encoding/json"
	sdk "github.com/NetCloth/netcloth-chain/types"
)

// MsgIPALClaim defines an ipal claim message
type MsgIPALClaim struct {
	From        sdk.AccAddress `json:"from" yaml:"from`
	UserAddress string         `json:"user_address" yaml:"user_address"`
	ServerIP    string         `json:"server_ip" yaml:"server_ip"`
}

var _ sdk.Msg = MsgIPALClaim{}

// NewMsgIPALClaim is a constructor function for MsgIPALClaim
func NewMsgIPALClaim(from sdk.AccAddress, userAddress string, serverIP string) MsgIPALClaim {
	return MsgIPALClaim{
		from,
		userAddress,
		serverIP,
	}
}

// Route should return the name of the module
func (msg MsgIPALClaim) Route() string { return RouterKey }

// Type should return the action
func (msg MsgIPALClaim) Type() string { return "ipal_claim" }

// ValidateBasic runs stateless checks on the message
func (msg MsgIPALClaim) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}

	if msg.UserAddress == "" {
		return ErrEmptyInputs("missing user address")
	}

	if msg.ServerIP == "" {
		return ErrEmptyInputs("missing server ip")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgIPALClaim) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgIPALClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}