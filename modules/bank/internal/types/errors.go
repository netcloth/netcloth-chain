package types

import (
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

var (
	ErrNoInputs            = sdkerrors.Register(ModuleName, 1, "no inputs to send transaction")
	ErrNoOutputs           = sdkerrors.Register(ModuleName, 2, "no outputs to send transaction")
	ErrInputOutputMismatch = sdkerrors.Register(ModuleName, 3, "sum inputs != sum outputs")
	ErrSendDisabled        = sdkerrors.Register(ModuleName, 4, "send transactions are disabled")
)
