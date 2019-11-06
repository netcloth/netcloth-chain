package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NetCloth/netcloth-chain/modules/auth"
	sdk "github.com/NetCloth/netcloth-chain/types"
)

const (
	maxUserAddressLength   = 256
	maxServerAddressLength = 256
)

var (
	_ sdk.Msg = MsgIPALClaim{}
)

type ServiceInfo struct {
	Type    uint64 `json:"type" yaml:"type"`
	Address string `json:"address" yaml:"address"`
}

type ADParam struct {
	UserAddress string      `json:"user_address" yaml:"user_address"`
	ServiceInfo ServiceInfo `json:"service_info" yaml:"service_info"`
	Expiration  time.Time   `json:"expiration"`
}

type IPALUserRequest struct {
	Params ADParam           `json:"params" yaml:"params"`
	Sig    auth.StdSignature `json:"signature" yaml:"signature`
}

type MsgIPALClaim struct {
	From        sdk.AccAddress  `json:"from" yaml:"from`
	UserRequest IPALUserRequest `json:"user_request" yaml:"user_request"`
}

func (i ServiceInfo) Validate() sdk.Error {
	if i.Address == "" {
		return ErrEmptyInputs("server address empty")
	}

	if len(i.Address) > maxServerAddressLength {
		return ErrStringTooLong("server address too long")
	}

	return nil
}

func (i ServiceInfo) String() string {
	return fmt.Sprintf(`ServiceInfo{Type:%s,Address:%s`, i.Type, i.Address)
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

	if len(p.UserAddress) > maxUserAddressLength {
		return ErrStringTooLong("user address too long")
	}

	return p.ServiceInfo.Validate()
}

func NewADParam(userAddress string, serviceAddress string, serviceType uint64, expiration time.Time) ADParam {
	return ADParam{
		UserAddress: userAddress,
		ServiceInfo: ServiceInfo{Type: serviceType, Address: serviceAddress},
		Expiration:  expiration,
	}
}

func NewIPALUserRequest(userAddress string, serviceAddress string, serviceType uint64, expiration time.Time, sig auth.StdSignature) IPALUserRequest {
	return IPALUserRequest{
		Params: NewADParam(userAddress, serviceAddress, serviceType, expiration),
		Sig:    sig,
	}
}

func NewMsgIPALClaim(from sdk.AccAddress, userAddress string, serviceAddress string, serviceType uint64, expiration time.Time, sig auth.StdSignature) MsgIPALClaim {
	return MsgIPALClaim{
		from,
		NewIPALUserRequest(userAddress, serviceAddress, serviceType, expiration, sig),
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
