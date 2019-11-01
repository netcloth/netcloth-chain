package types

import (
	"encoding/json"
	"time"

	"github.com/NetCloth/netcloth-chain/modules/auth"
	sdk "github.com/NetCloth/netcloth-chain/types"
)

const (
	maxUserAddressLength = 64
	maxServerIPLength    = 64
)

var (
	_ sdk.Msg = MsgIPALClaim{}
)

type ADParam struct {
	UserAddress string    `json:"user_address" yaml:"user_address"`
	ServerIP    string    `json:"server_ip" yaml:"server_ip"`
	Expiration  time.Time `json:"expiration"`
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

func (p ADParam) GetSignBytes() []byte {
	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (p ADParam) Validate() sdk.Error {
	if p.UserAddress == "" {
		return ErrEmptyInputs("user address empty")
	}

	if p.ServerIP == "" {
		return ErrEmptyInputs("server ip empty")
	}

	if len(p.UserAddress) > maxUserAddressLength {
		return ErrStringTooLong("user address too long")
	}

	if len(p.ServerIP) > maxServerIPLength {
		return ErrStringTooLong("server ip too long")
	}

	return nil
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

func NewMsgIPALClaim(from sdk.AccAddress, userAddress string, serverIP string, expiration time.Time, sig auth.StdSignature) MsgIPALClaim {
	return MsgIPALClaim{
		from,
		NewIPALUserRequest(userAddress, serverIP, expiration, sig),
	}
}

func (msg MsgIPALClaim) Route() string { return RouterKey }

func (msg MsgIPALClaim) Type() string { return "ipal_claim" }

func (msg MsgIPALClaim) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}

	err := msg.UserRequest.Params.Validate()
	if err != nil {
		return err
	}

	pubKey := msg.UserRequest.Sig.PubKey
	signBytes := msg.UserRequest.Params.GetSignBytes()
	if !pubKey.VerifyBytes(signBytes, msg.UserRequest.Sig.Signature) {
		return ErrInvalidSignature("user request signature invalid")
	}

	return nil
}

func (msg MsgIPALClaim) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgIPALClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}
