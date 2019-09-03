package token

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is they name of the bank module
const TokenRouterKey = "token"

//----------------------------------------
// MsgIssue

// MsgIssue - high level transaction of the coin module
type MsgIssue struct {
	Banker  sdk.AccAddress `json:"banker"`
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`}

var _ sdk.Msg = MsgIssue{}

// NewMsgIssue - construct arbitrary multi-in, multi-out send msg.
func NewMsgIssue(banker sdk.AccAddress, addr sdk.AccAddress, coins sdk.Coins) MsgIssue {
	return MsgIssue{
		banker,
		addr,
		coins,
	}
}

// Implements Msg.
// nolint
func (msg MsgIssue) Route() string { return TokenRouterKey }
func (msg MsgIssue) Type() string  { return "issue" }

// Implements Msg.
func (msg MsgIssue) ValidateBasic() sdk.Error {
	if len(msg.Address) == 0 {
		return sdk.ErrInvalidAddress(msg.Address.String())
	}
	if !msg.Coins.IsValid() {
		return sdk.ErrInvalidCoins("issue coins invalid " +  msg.Coins.String())
	}
	if !msg.Coins.IsAllPositive() {
		return sdk.ErrInvalidCoins("issue amount must be positive " + msg.Coins.String())
	}
	return nil
}

// Implements Msg.
func (msg MsgIssue) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgIssue) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Banker}
}