package types

import (
	"fmt"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

var (
	ErrEmptyInputs        = sdkerrors.Register(ModuleName, 1, "input empty")
	ErrBadDenom           = sdkerrors.Register(ModuleName, 2, "bad denom")
	ErrBondInsufficient   = sdkerrors.Register(ModuleName, 3, "bond insufficient")
	ErrMonikerExist       = sdkerrors.Register(ModuleName, 4, "moniker exists")
	ErrEndpointsEmpty     = sdkerrors.Register(ModuleName, 6, "no endpoints")
	ErrEndpointsDuplicate = sdkerrors.Register(ModuleName, 7, "endpoints duplicate")
)

func ErrEndpointsFormat(flagEndpointsDecollator, flagEndpointDecollator string) *sdkerrors.Error {
	return sdkerrors.Register(ModuleName, 5, fmt.Sprintf("endpoints format err, should be in format: serviceType%sendpoint%sserviceType%sendpoint, serviceType is a number, endpoint is a string", flagEndpointDecollator, flagEndpointsDecollator, flagEndpointDecollator))
}

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
