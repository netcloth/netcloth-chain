package types

import (
	"encoding/json"
	sdk "github.com/NetCloth/netcloth-chain/types"
)

// MsgSend defines a transfer message
type MsgSend struct {
	From   sdk.AccAddress `json:"from_address" yaml:"from_address"`
	To     sdk.AccAddress `json:"to_address" yaml:"to_address"`
	Amount sdk.Coins      `json:"amount" yaml:"amount"`
}

var _ sdk.Msg = MsgSend{}

// NewMsgSend is a constructor function for MsgSend
func NewMsgSend(from sdk.AccAddress, to sdk.AccAddress, amount sdk.Coins) MsgSend {
	return MsgSend{
		from,
		to,
		amount,
	}
}

// Route should return the name of the module
func (msg MsgSend) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSend) Type() string { return "send" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSend) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}

	if msg.To.Empty() {
		return sdk.ErrInvalidAddress("missing recipient address")
	}

	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}

	if !msg.Amount.IsAllPositive() {
		return sdk.ErrInsufficientCoins("send amount must be positive")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSend) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}
