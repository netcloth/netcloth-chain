package types

import (
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

var (
	ErrNoInputs            = sdkerrors.New(ModuleName, 1, "no inputs to send transaction")
	ErrNoOutputs           = sdkerrors.New(ModuleName, 2, "no outputs to send transaction")
	ErrInputOutputMismatch = sdkerrors.New(ModuleName, 3, "sum inputs != sum outputs")
	ErrSendDisabled        = sdkerrors.New(ModuleName, 4, "send transactions are disabled")
)
