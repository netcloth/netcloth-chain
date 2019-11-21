package types

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

var (
	_ sdk.Msg = &MsgContractCreate{}
	_ sdk.Msg = &MsgContractCall{}
)

// MsgContractCreate - struct for contract create
type MsgContractCreate struct {
	From   sdk.AccAddress `json:"from" yaml:"from"`
	Amount sdk.Coins      `json:"amount" yaml:"amount"`
	Code   []byte         `json:"code" yaml:"code"`
}

// MsgContractCall - struct for contract call
type MsgContractCall struct {
	From      sdk.AccAddress `json:"from" yaml:"from"`
	Recipient sdk.AccAddress `json:"recipient" yaml:"recipient"`
	Amount    sdk.Coins      `json:"amount" yaml:"amount"`
	Payload   []byte         `json:"payload" yaml:"payload"`
}

// MsgContractCreate
func (msg MsgContractCreate) Route() string { return RouterKey }
func (msg MsgContractCreate) Type() string  { return "create_contract" }

func (msg MsgContractCreate) GetSigners() []sdk.AccAddress {
	return nil
}

func (msg MsgContractCreate) GetSignBytes() []byte {
	return nil
}

func (msg MsgContractCreate) ValidateBasic() sdk.Error {
	return nil
}

// MsgContractCall
func (msg MsgContractCall) Route() string { return RouterKey }
func (msg MsgContractCall) Type() string  { return "call_contract" }

func (msg MsgContractCall) GetSigners() []sdk.AccAddress {
	return nil
}

func (msg MsgContractCall) GetSignBytes() []byte {
	return nil
}

func (msg MsgContractCall) ValidateBasic() sdk.Error {
	return nil
}
