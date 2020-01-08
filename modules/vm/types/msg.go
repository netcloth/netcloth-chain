package types

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	TypeMsgContract = "contract"
)

var (
	_ sdk.Msg = &MsgContract{}
)

type MsgContract struct {
	From    sdk.AccAddress `json:"from" yaml:"from"`
	To      sdk.AccAddress `json:"to" yaml:"to"`
	Payload []byte         `json:"payload" yaml:"payload"`
	Amount  sdk.Coin       `json:"amount" yaml:"amout"`
}

func (m MsgContract) Route() string {
	return RouterKey
}

func (m MsgContract) Type() string {
	return TypeMsgContract
}

func (m MsgContract) ValidateBasic() sdk.Error {
	if m.From.Empty() {
		return sdk.ErrInvalidAddress("msg missing from address")
	}
	if !m.Amount.IsValid() {
		return sdk.ErrInvalidCoins("msg amount is invalid: " + m.Amount.String())
	}
	if !m.Amount.IsPositive() && !m.Amount.IsZero() {
		return sdk.ErrInsufficientCoins("msg amount must be positive")
	}
	return nil
}

func (m MsgContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgContract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.From}
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
