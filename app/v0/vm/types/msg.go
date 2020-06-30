package types

import (
	"github.com/netcloth/netcloth-chain/hexutil"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

const (
	TypeMsgContractCreate = "contract_create"
	TypeMsgContractCall   = "contract_call"
)

var (
	_ sdk.Msg = &MsgContract{}
)

type MsgContract struct {
	From    sdk.AccAddress `json:"from" yaml:"from"`
	To      sdk.AccAddress `json:"to" yaml:"to"`
	Payload hexutil.Bytes  `json:"payload" yaml:"payload"`
	Amount  sdk.Coin       `json:"amount" yaml:"amount"`
}

func (msg MsgContract) Route() string {
	return RouterKey
}

func (msg MsgContract) Type() string {
	if msg.To.Empty() {
		return TypeMsgContractCreate
	}
	return TypeMsgContractCall
}

func (msg MsgContract) ValidateBasic() error {
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "msg missing from address")
	}
	if !msg.Amount.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "msg amount is invalid: "+msg.Amount.String())
	}
	if msg.Amount.Denom != sdk.NativeTokenName {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "denom must be %s", sdk.NativeTokenName)
	}
	if len(msg.Payload) == 0 {
		return ErrNoPayload
	}

	return nil
}

func (msg MsgContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgContract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func NewMsgContract(from, to sdk.AccAddress, payload []byte, amount sdk.Coin) MsgContract {
	return MsgContract{
		From:    from,
		To:      to,
		Payload: payload,
		Amount:  amount,
	}
}

type MsgContractQuery MsgContract

func NewMsgContractQuery(from, to sdk.AccAddress, payload []byte, amount sdk.Coin) MsgContractQuery {
	return MsgContractQuery{
		From:    from,
		To:      to,
		Payload: payload,
		Amount:  amount,
	}
}
