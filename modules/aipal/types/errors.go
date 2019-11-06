package types

import (
	"fmt"

	sdk "github.com/NetCloth/netcloth-chain/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeEmptyInputs           sdk.CodeType = 100
	CodeEndpointsFormatErr    sdk.CodeType = 102
	CodeEndpointsEmptyErr     sdk.CodeType = 103
	CodeEndpointsDuplicateErr sdk.CodeType = 104
	CodeBadDenom              sdk.CodeType = 111
	CodeBondInsufficient      sdk.CodeType = 112
	CodeMonikerExist          sdk.CodeType = 113
)

func ErrEmptyInputs(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeEmptyInputs, msg)
}

func ErrBadDenom(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeBadDenom, msg)
}

func ErrBondInsufficient(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeBondInsufficient, msg)
}

func ErrMonikerExist(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeMonikerExist, msg)
}

func ErrEndpointsFormat() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeEndpointsFormatErr, "endpoints format err, should be in format: serviceType|endpoint,serviceType|endpoint, serviceType is a number, endpoint is a string")
}

func ErrEndpointsEmpty() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeEndpointsEmptyErr, "no endpoints")
}

func ErrEndpointsDuplicate(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeEndpointsDuplicateErr, msg)
}

type EndpointDuplicateErrDetector struct {
	V map[int]int
}

func (d *EndpointDuplicateErrDetector) detecte(t int) sdk.Error {
	d.V[t]++

	if d.V[t] > 1 {
		return ErrEndpointsDuplicate(fmt.Sprintf("endpoint type: [%d] is duplicate", t))
	}

	return nil
}

func EndpointsDupCheck(eps Endpoints) sdk.Error {
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
