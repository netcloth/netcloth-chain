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
	_ sdk.Msg = MsgServiceNodeClaim{}
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

type MsgServiceNodeClaim struct {
	OperatorAddress sdk.AccAddress `json:"operator_address" yaml:"operator_address"` // address of the ServiceNode's operator
	Moniker         string         `json:"moniker" yaml:"moniker"`                   // name
	Website         string         `json:"website" yaml:"website"`                   // optional website link
	ServerEndPoint  string         `json:"server_endpoint" yaml:"server_endpoint"`   // server endpoint for app client
	Details         string         `json:"details" yaml:"details"`                   // optional details
	StakeShares     sdk.Coin       `json:"stake_shares" yaml:"stake_shares"`         // total stake shares
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
		return ErrInvalidSignature("user request signature invalid")
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

func NewMsgServiceNodeClaim(operator sdk.AccAddress, moniker, website, serverEndPoint, details string, amount sdk.Coin) MsgServiceNodeClaim {
	return MsgServiceNodeClaim{
		OperatorAddress: operator,
		Moniker:         moniker,
		Website:         website,
		ServerEndPoint:  serverEndPoint,
		Details:         details,
		StakeShares:     amount,
	}
}

func (msg MsgServiceNodeClaim) Route() string { return RouterKey }

func (msg MsgServiceNodeClaim) Type() string { return "server_node_claim" }

// ValidateBasic runs stateless checks on the message
func (msg MsgServiceNodeClaim) ValidateBasic() sdk.Error {
	// check msg sender
	if msg.OperatorAddress.Empty() {
		return sdk.ErrInvalidAddress("missing operator address")
	}

	if msg.Moniker == "" {
		return ErrEmptyInputs("moniker empty")
	}

	if msg.Website == "" {
		return ErrEmptyInputs("website empty")
	}

	if msg.ServerEndPoint == "" {
		return ErrEmptyInputs("server empty")
	}

	if msg.StakeShares.IsNegative() {
		return ErrEmptyInputs("stake amount must > 0 ")
	}

	return nil
}

func (msg MsgServiceNodeClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.OperatorAddress)}
}

func (msg MsgServiceNodeClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
