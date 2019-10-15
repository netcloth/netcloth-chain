package types

import (
	"encoding/json"
	"github.com/NetCloth/netcloth-chain/modules/auth"
	sdk "github.com/NetCloth/netcloth-chain/types"
	"time"
)

const (
	maxStringLength = 64
)

type ADParam struct {
	UserAddress string    `json:"user_address" yaml:"user_address"`
	ServerIP    string    `json:"server_ip" yaml:"server_ip"`
	Expiration  time.Time `json:"expiration"`
}

func (p ADParam) GetSignBytes() []byte {
	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (p ADParam) Validate() sdk.Error {
	if p.UserAddress == "" {
		return ErrEmptyInputs(DefaultCodespace)
	}

	if p.ServerIP == "" {
		return ErrEmptyInputs(DefaultCodespace)
	}

	if len(p.UserAddress) > maxStringLength {
		return ErrStringTooLong(DefaultCodespace)
	}

	if len(p.ServerIP) > maxStringLength {
		return ErrStringTooLong(DefaultCodespace)
	}

	return nil
}

type IPALUserRequest struct {
	Params ADParam           `json:"params" yaml:"params"`
	Sig    auth.StdSignature `json:"signature" yaml:"signature`
}

// MsgIPALClaim defines an ipal claim message
type MsgIPALClaim struct {
	From        sdk.AccAddress  `json:"from" yaml:"from`
	UserRequest IPALUserRequest `json:"user_request" yaml:"user_request"`
}

func NewADParam(userAddress string, serverIP string, expiration time.Time) ADParam {
	return ADParam{
		UserAddress: userAddress,
		ServerIP:    serverIP,
		Expiration:  expiration,
	}
}

func NewIPALUserRequest(userAddress string, serverIP string, expiration time.Time, sig auth.StdSignature) IPALUserRequest {
	return IPALUserRequest{
		Params: NewADParam(userAddress, serverIP, expiration),
		Sig:    sig,
	}
}

var _ sdk.Msg = MsgIPALClaim{}

// NewMsgIPALClaim is a constructor function for MsgIPALClaim
func NewMsgIPALClaim(from sdk.AccAddress, userAddress string, serverIP string, expiration time.Time, sig auth.StdSignature) MsgIPALClaim {
	return MsgIPALClaim{
		from,
		NewIPALUserRequest(userAddress, serverIP, expiration, sig),
	}
}

// Route should return the name of the module
func (msg MsgIPALClaim) Route() string { return RouterKey }

// Type should return the action
func (msg MsgIPALClaim) Type() string { return "ipal_claim" }

// ValidateBasic runs stateless checks on the message
func (msg MsgIPALClaim) ValidateBasic() sdk.Error {
	// check msg sender
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}

	// check userAddress and serverIP
	err := msg.UserRequest.Params.Validate()
	if err != nil {
		return err
	}

	// check user request signature
	pubKey := msg.UserRequest.Sig.PubKey
	signBytes := msg.UserRequest.Params.GetSignBytes()
	if !pubKey.VerifyBytes(signBytes, msg.UserRequest.Sig.Signature) {
		return ErrInvalidSignature(DefaultCodespace)
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
