package token

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is they name of the bank module
const TokenRouterKey = "token"

//----------------------------------------
// MsgTokenIssue

// MsgTokenIssue - high level transaction of the coin module
type MsgTokenIssue struct {
	Banker  sdk.AccAddress `json:"banker"`
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`}

var _ sdk.Msg = MsgTokenIssue{}

// NewMsgTokenIssue - construct arbitrary multi-in, multi-out send msg.
func NewMsgTokenIssue(banker sdk.AccAddress, addr sdk.AccAddress, coins sdk.Coins) MsgTokenIssue {
	return MsgTokenIssue{
		banker,
		addr,
		coins,
	}
}

// Implements Msg.
// nolint
func (msg MsgTokenIssue) Route() string { return TokenRouterKey }
func (msg MsgTokenIssue) Type() string  { return "issue" }

// Implements Msg.
func (msg MsgTokenIssue) ValidateBasic() sdk.Error {
	if len(msg.Address) == 0 {
		return sdk.ErrInvalidAddress(msg.Address.String())
	}
	if !msg.Coins.IsValid() {
		return sdk.ErrInvalidCoins(msg.Coins.String())
	}
	if !msg.Coins.IsAnyNegative() {
		return sdk.ErrInvalidCoins(msg.Coins.String())
	}
	return nil
}

// Implements Msg.
func (msg MsgTokenIssue) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgTokenIssue) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Banker}
}