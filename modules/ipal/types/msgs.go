package types

import (
	"encoding/json"
	sdk "github.com/NetCloth/netcloth-chain/types"
	"go/types"
	"time"
)

const (
	maxStringLength = 64
)

type IPALUserRequest struct {
	UserAddress string          `json:"user_address" yaml:"user_address"`
	ServerIP    string          `json:"server_ip" yaml:"server_ip"`
	Expiration  time.Time       `json:"expiration"`
	Sig         types.Signature `json:"signature" yaml:"signature`
}

// MsgIPALClaim defines an ipal claim message
type MsgIPALClaim struct {
	From        sdk.AccAddress  `json:"from" yaml:"from`
	UserRequest IPALUserRequest `json: "user_request" yaml:"user_request"`
}

func NewIPALUserRequest(userAddress string, serverIP string, expiration time.Time) IPALUserRequest {
	return IPALUserRequest{
		UserAddress: userAddress,
		ServerIP:    serverIP,
		Expiration:  expiration,
		Sig:         types.Signature{},
	}
}

var _ sdk.Msg = MsgIPALClaim{}

// NewMsgIPALClaim is a constructor function for MsgIPALClaim
func NewMsgIPALClaim(from sdk.AccAddress, userAddress string, serverIP string, expiration time.Time) MsgIPALClaim {
	return MsgIPALClaim{
		from,
		NewIPALUserRequest(userAddress, serverIP, expiration),
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

	if msg.UserRequest.UserAddress == "" {
		return ErrEmptyInputs(DefaultCodespace)
	}

	if msg.UserRequest.ServerIP == "" {
		return ErrEmptyInputs(DefaultCodespace)
	}

	if len(msg.UserRequest.UserAddress) > maxStringLength {
		return ErrStringTooLong(DefaultCodespace)
	}

	if len(msg.UserRequest.ServerIP) > maxStringLength {
		return ErrStringTooLong(DefaultCodespace)
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
