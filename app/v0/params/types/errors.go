package types

import (
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

var (
	ErrUnknownSubspace  = sdkerrors.New(ModuleName, 1, "unknown subspace")
	ErrSettingParameter = sdkerrors.New(ModuleName, 2, "failed to set parameter")
	ErrEmptyChanges     = sdkerrors.New(ModuleName, 3, "submitted parameter changes are empty")
	ErrEmptySubspace    = sdkerrors.New(ModuleName, 4, "parameter subspace is empty")
	ErrEmptyKey         = sdkerrors.New(ModuleName, 5, "parameter key is empty")
	ErrEmptyValue       = sdkerrors.New(ModuleName, 6, "parameter value is empty")
)
