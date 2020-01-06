package types

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	TypeMsgContractCreate = "create_contract"
	TypeMsgContractCall   = "call_contract"
)

var (
	_ sdk.Msg = &MsgContractCreate{}
	_ sdk.Msg = &MsgContractCall{}
)

// MsgContractCreate - struct for contract create
type MsgContractCreate struct {
	From   sdk.AccAddress `json:"from" yaml:"from"`
	Amount sdk.Coin       `json:"amount" yaml:"amount"`
	Code   []byte         `json:"code" yaml:"code"`
}

func NewMsgContractCreate(From sdk.AccAddress, Amount sdk.Coin, Code []byte) MsgContractCreate {
	return MsgContractCreate{
		From:   From,
		Amount: Amount,
		Code:   Code,
	}
}

// MsgContractCall - struct for contract call
type MsgContractCall struct {
	From      sdk.AccAddress `json:"from" yaml:"from"`
	Recipient sdk.AccAddress `json:"recipient" yaml:"recipient"`
	Amount    sdk.Coin       `json:"amount" yaml:"amount"`
	Payload   []byte         `json:"payload" yaml:"payload"`
}

func NewMsgContractCall(From, to sdk.AccAddress, Amount sdk.Coin, args []byte) MsgContractCall {
	return MsgContractCall{
		From:      From,
		Recipient: to,
		Amount:    Amount,
		Payload:   args,
	}
}

// MsgContractCreate
func (msg MsgContractCreate) Route() string { return RouterKey }
func (msg MsgContractCreate) Type() string  { return TypeMsgContractCreate }

func (msg MsgContractCreate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func (msg MsgContractCreate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgContractCreate) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("msg missing sender address")
	}
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("msg amount is invalid: " + msg.Amount.String())
	}
	if !msg.Amount.IsPositive() && !msg.Amount.IsZero() {
		return sdk.ErrInsufficientCoins("msg amount must be positive")
	}

	if len(msg.Code) == 0 {
		return ErrNoCodeExist()
	}

	return nil
}

// MsgContractCall
func (msg MsgContractCall) Route() string { return RouterKey }
func (msg MsgContractCall) Type() string  { return TypeMsgContractCall }

func (msg MsgContractCall) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func (msg MsgContractCall) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgContractCall) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("msg missing sender address")
	}
	if msg.Recipient.Empty() {
		return sdk.ErrInvalidAddress("missing recipient address")
	}
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("msg amount is invalid: " + msg.Amount.String())
	}
	if !msg.Amount.IsPositive() && !msg.Amount.IsZero() {
		return sdk.ErrInsufficientCoins("msg amount must be positive")
	}
	return nil
}
