package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/netcloth/netcloth-chain/app/v0/auth"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

const (
	maxUserAddressLength = 256
)

var (
	_ sdk.Msg = MsgCIPALClaim{}
)

// ServiceInfo defines the struct of service type and service address
type ServiceInfo struct {
	Type    uint64 `json:"type" yaml:"type"`
	Address string `json:"address" yaml:"address"`
}

// ADParam defines the struct of user address and service info
type ADParam struct {
	UserAddress string      `json:"user_address" yaml:"user_address"`
	ServiceInfo ServiceInfo `json:"service_info" yaml:"service_info"`
	Expiration  time.Time   `json:"expiration"`
}

// CIPALUserRequest defines the struct of user request for CIPAL claim
type CIPALUserRequest struct {
	Params ADParam           `json:"params" yaml:"params"`
	Sig    auth.StdSignature `json:"signature" yaml:"signature"`
}

// MsgCIPALClaim defines the transaction struct of CIPAL claim
type MsgCIPALClaim struct {
	From        sdk.AccAddress   `json:"from" yaml:"from"`
	UserRequest CIPALUserRequest `json:"user_request" yaml:"user_request"`
}

func (i ServiceInfo) String() string {
	return fmt.Sprintf(`ServiceInfo{Type:%d, Address:%s`, i.Type, i.Address)
}

// GetSignBytes - get the bytes for the message signer to sign on
func (p ADParam) GetSignBytes() []byte {
	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Validate - quick validity check
func (p ADParam) Validate() error {
	if p.UserAddress == "" {
		return sdkerrors.Wrap(ErrEmptyInputs, "user address empty")
	}

	if len(p.UserAddress) > maxUserAddressLength {
		return sdkerrors.Wrap(ErrStringTooLong, "user address too long")
	}

	return nil
}

// NewADParam - create a new instance of ADParam
func NewADParam(userAddress string, serviceAddress string, serviceType uint64, expiration time.Time) ADParam {
	return ADParam{
		UserAddress: userAddress,
		ServiceInfo: ServiceInfo{Type: serviceType, Address: serviceAddress},
		Expiration:  expiration,
	}
}

// NewCIPALUserRequest - create a new instance of CIPALUserRequest
func NewCIPALUserRequest(userAddress string, serviceAddress string, serviceType uint64, expiration time.Time, sig auth.StdSignature) CIPALUserRequest {
	return CIPALUserRequest{
		Params: NewADParam(userAddress, serviceAddress, serviceType, expiration),
		Sig:    sig,
	}
}

// NewMsgCIPALClaim - create a new instance of MsgCIPALClaim
func NewMsgCIPALClaim(from sdk.AccAddress, userAddress string, serviceAddress string, serviceType uint64, expiration time.Time, sig auth.StdSignature) MsgCIPALClaim {
	return MsgCIPALClaim{
		from,
		NewCIPALUserRequest(userAddress, serviceAddress, serviceType, expiration, sig),
	}
}

// Route Implements Msg
func (msg MsgCIPALClaim) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgCIPALClaim) Type() string { return "cipal_claim" }

// ValidateBasic Implements Msg
func (msg MsgCIPALClaim) ValidateBasic() error {
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}

	err := msg.UserRequest.Params.Validate()
	if err != nil {
		return err
	}

	pubKey := msg.UserRequest.Sig.PubKey
	signBytes := msg.UserRequest.Params.GetSignBytes()
	if !pubKey.VerifyBytes(signBytes, msg.UserRequest.Sig.Signature) {
		return sdkerrors.Wrap(ErrInvalidSignature, "user request signature invalid")
	}

	return nil
}

// GetSignBytes Implements Msg
func (msg MsgCIPALClaim) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgCIPALClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}
