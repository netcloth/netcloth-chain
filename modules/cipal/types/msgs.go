package types

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/tendermint/tendermint/crypto"

	"github.com/NetCloth/netcloth-chain/modules/auth"
	sdk "github.com/NetCloth/netcloth-chain/types"
)

const (
	maxUserAddressLength   = 256
	maxServerAddressLength = 256
)

var (
	_ sdk.Msg = MsgCIPALClaim{}
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

type CIPALUserRequest struct {
	Params ADParam           `json:"params" yaml:"params"`
	Sig    auth.StdSignature `json:"signature" yaml:"signature`
}

type MsgCIPALClaim struct {
	From        sdk.AccAddress   `json:"from" yaml:"from`
	UserRequest CIPALUserRequest `json:"user_request" yaml:"user_request"`
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

func NewCIPALUserRequest(userAddress string, serviceAddress string, serviceType uint64, expiration time.Time, sig auth.StdSignature) CIPALUserRequest {
	return CIPALUserRequest{
		Params: NewADParam(userAddress, serviceAddress, serviceType, expiration),
		Sig:    sig,
	}
}

func NewMsgCIPALClaim(from sdk.AccAddress, userAddress string, serviceAddress string, serviceType uint64, expiration time.Time, sig auth.StdSignature) MsgCIPALClaim {
	return MsgCIPALClaim{
		from,
		NewCIPALUserRequest(userAddress, serviceAddress, serviceType, expiration, sig),
	}
}

func (msg MsgCIPALClaim) Route() string { return RouterKey }

func (msg MsgCIPALClaim) Type() string { return "cipal_claim" }

func (msg MsgCIPALClaim) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}

	err := msg.UserRequest.Params.Validate()
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "++++++++++++++++++++++++++++++++++\n")
	d, _ := msg.UserRequest.Sig.MarshalYAML()
	fmt.Fprintf(os.Stderr, fmt.Sprintf("params = %s\n", msg.UserRequest.Params.GetSignBytes()))
	fmt.Fprintf(os.Stderr, fmt.Sprintf("params hex = %x\n", msg.UserRequest.Params.GetSignBytes()))
	fmt.Fprintf(os.Stderr, fmt.Sprintf("params hex sha256 = %x\n", crypto.Sha256(msg.UserRequest.Params.GetSignBytes())))
	fmt.Fprintf(os.Stderr, fmt.Sprintf("user sig = %s\n", d))
	fmt.Fprintf(os.Stderr, fmt.Sprintf("user sig.pubkey = %v\n", msg.UserRequest.Sig.PubKey))
	fmt.Fprintf(os.Stderr, fmt.Sprintf("user addr = %x\n", msg.UserRequest.Sig.PubKey.Address()))
	fmt.Fprintf(os.Stderr, fmt.Sprintf("user sig hex = %x\n", msg.UserRequest.Sig.Signature))
	fmt.Fprintf(os.Stderr, "----------------------------------\n")
	pubKey := msg.UserRequest.Sig.PubKey
	signBytes := msg.UserRequest.Params.GetSignBytes()
	if !pubKey.VerifyBytes(signBytes, msg.UserRequest.Sig.Signature) {
		return ErrInvalidSignature("user request signature invalid")
	}

	return nil
}

func (msg MsgCIPALClaim) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgCIPALClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}
