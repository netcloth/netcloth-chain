package subspace

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestKeyTable(t *testing.T) {
	table := NewKeyTable()

	require.Panics(t, func() { table.RegisterType(ParamSetPair{[]byte(""), nil, nil}) })
	require.Panics(t, func() { table.RegisterType(ParamSetPair{[]byte("!@#$%"), nil, nil}) })
	require.Panics(t, func() { table.RegisterType(ParamSetPair{[]byte("hello,"), nil, nil}) })
	require.Panics(t, func() { table.RegisterType(ParamSetPair{[]byte("hello"), nil, nil}) })

	require.NotPanics(t, func() {
		table.RegisterType(ParamSetPair{keyBondDenom, string("stake"), validateBondDenom})
	})
	require.NotPanics(t, func() {
		table.RegisterType(ParamSetPair{keyMaxValidators, uint16(100), validateMaxValidators})
	})
	require.Panics(t, func() {
		table.RegisterType(ParamSetPair{keyUnbondingTime, time.Duration(1), nil})
	})
	require.NotPanics(t, func() {
		table.RegisterType(ParamSetPair{keyUnbondingTime, time.Duration(1), validateMaxValidators})
	})
	require.NotPanics(t, func() {
		newTable := NewKeyTable()
		newTable.RegisterParamSet(&params{})
	})

	require.Panics(t, func() { table.RegisterParamSet(&params{}) })
	require.Panics(t, func() { NewKeyTable(ParamSetPair{[]byte(""), nil, nil}) })

	require.NotPanics(t, func() {
		NewKeyTable(
			ParamSetPair{[]byte("test"), string("stake"), validateBondDenom},
			ParamSetPair{[]byte("test2"), uint16(100), validateMaxValidators},
		)
	})
}
