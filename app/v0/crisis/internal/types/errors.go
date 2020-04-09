package types

import (
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

var (
	ErrNoSender         = sdkerrors.New(ModuleName, 1, "sender address is empty")
	ErrUnknownInvariant = sdkerrors.New(ModuleName, 2, "unknown invariant")
)
