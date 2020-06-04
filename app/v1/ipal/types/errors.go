package types

import (
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

var (
	ErrEmptyInputs        = sdkerrors.New(ModuleName, 1, "input empty")
	ErrBadDenom           = sdkerrors.New(ModuleName, 2, "bad denom")
	ErrBondInsufficient   = sdkerrors.New(ModuleName, 3, "bond insufficient")
	ErrMonikerExist       = sdkerrors.New(ModuleName, 4, "moniker exists")
	ErrEndpointsEmpty     = sdkerrors.New(ModuleName, 6, "no endpoints")
	ErrEndpointsDuplicate = sdkerrors.New(ModuleName, 7, "endpoints duplicate")
	ErrEndpointsFormat    = sdkerrors.New(ModuleName, 8, "endpoints format error")
)

type EndpointDuplicateErrDetector struct {
	V map[int]int
}

func (d *EndpointDuplicateErrDetector) detecte(t int) error {
	d.V[t]++

	if d.V[t] > 1 {
		return sdkerrors.Wrapf(ErrEndpointsDuplicate, "endpoint type: [%d] is duplicate", t)
	}

	return nil
}

func EndpointsDupCheck(eps Endpoints) error {
	d := EndpointDuplicateErrDetector{
		V: make(map[int]int),
	}

	for _, v := range eps {
		if e := d.detecte(int(v.Type)); e != nil {
			return e
		}
	}

	return nil
}
