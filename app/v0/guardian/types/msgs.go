package types

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

var _, _ sdk.Msg = MsgAddProfiler{}, MsgDeleteProfiler{}

type MsgAddProfiler struct {
	AddGuardian
}

func NewMsgAddProfiler(description string, address, addedBy sdk.AccAddress) MsgAddProfiler {
	return MsgAddProfiler{
		AddGuardian: AddGuardian{
			Description: description,
			Address:     address,
			AddedBy:     addedBy,
		},
	}
}

func (m MsgAddProfiler) Route() string {
	return RouterKey
}

func (m MsgAddProfiler) Type() string {
	return "MsgAddProfiler"
}

func (m MsgAddProfiler) ValidateBasic() error {
	return m.AddGuardian.ValidateBasic()
}

func (m MsgAddProfiler) GetSignBytes() []byte {
	bz := msgCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgAddProfiler) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.AddedBy}
}

type MsgDeleteProfiler struct {
	DeleteGuardian
}

func NewMsgDeleteProfiler(address, deletedBy sdk.AccAddress) MsgDeleteProfiler {
	return MsgDeleteProfiler{
		DeleteGuardian: DeleteGuardian{
			Address:   address,
			DeletedBy: deletedBy,
		},
	}
}

func (m MsgDeleteProfiler) Route() string {
	return RouterKey
}

func (m MsgDeleteProfiler) Type() string {
	return "MsgDeleteProfiler"
}

func (m MsgDeleteProfiler) ValidateBasic() error {
	return m.DeleteGuardian.ValidateBasic()
}

func (m MsgDeleteProfiler) GetSignBytes() []byte {
	bz := msgCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgDeleteProfiler) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.DeletedBy}
}

type AddGuardian struct {
	Description string         `json:"description"`
	Address     sdk.AccAddress `json:"address"`
	AddedBy     sdk.AccAddress `json:"added_by"`
}

func (g AddGuardian) ValidateBasic() error {
	if len(g.Description) == 0 {
		return ErrInvalidDescription()
	}

	if len(g.Address) == 0 {
		return ErrAddressEmpty()
	}

	if len(g.AddedBy) == 0 {
		return ErrAddedByEmpty()
	}

	if err := g.EnsureLength(); err != nil {
		return err
	}

	return nil
}

func (g AddGuardian) EnsureLength() error {
	if len(g.Description) > 70 {
		return ErrInvalidDescription()
	}

	return nil
}

type DeleteGuardian struct {
	Address   sdk.AccAddress `json:"address"`
	DeletedBy sdk.AccAddress `json:"deleted_by"`
}

func (g DeleteGuardian) ValidateBasic() error {
	if len(g.Address) == 0 {
		return ErrAddressEmpty()
	}

	if len(g.DeletedBy) == 0 {
		return ErrDeletedByEmpty()
	}
	return nil
}
