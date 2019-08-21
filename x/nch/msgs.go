package nch

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/NetCloth/netcloth-chain/types"
)

const (
	MinTransferFee = 1
)

// MsgTransfer defines a transfer message
type MsgTransfer struct {
	From  sdk.AccAddress
	To    sdk.AccAddress
	Value sdk.Coins
	Fee   sdk.Coins
}

// NewMsgTransfer is a constructor function for MsgTransfer
func NewMsgTransfer(from sdk.AccAddress, to sdk.AccAddress, value sdk.Coins) MsgTransfer {
	return MsgTransfer{
		from,
		to,
		value,
		sdk.Coins{sdk.NewInt64Coin(types.AppCoin, MinTransferFee)},
	}
}

// Route should return the name of the module
func (msg MsgTransfer) Route() string { return "nch" }

// Type should return the action
func (msg MsgTransfer) Type() string { return "transfer" }

// ValidateBasic runs stateless checks on the message
func (msg MsgTransfer) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("missing from address")
	}

	if msg.To.Empty() {
		return sdk.ErrInvalidAddress("missing recipient address")
	}

	if !msg.Value.IsValid() {
		return sdk.ErrInvalidCoins("transfer amount is invalid: " + msg.Value.String())
	}

	if !msg.Value.IsAllPositive() {
		return sdk.ErrInsufficientCoins("transfer amount must be positive")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgTransfer) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}