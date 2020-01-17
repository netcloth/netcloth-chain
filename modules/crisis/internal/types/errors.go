package types

import (
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

var (
	ErrNoSender         = sdkerrors.Register(ModuleName, 1, "sender address is empty")
	ErrUnknownInvariant = sdkerrors.Register(ModuleName, 2, "unknown invariant")
)
