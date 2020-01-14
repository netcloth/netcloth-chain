package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatorTestEquivalent(t *testing.T) {
	val1 := NewValidator(valAddr1, pk1, Description{})
	val2 := NewValidator(valAddr1, pk1, Description{})

	ok := val1.TestEquivalent(val2)
	require.True(t, ok)

	val2 = NewValidator(valAddr2, pk2, Description{})

	ok = val1.TestEquivalent(val2)
	require.False(t, ok)
}

func TestUpdateDescription(t *testing.T) {
	d1 := Description{
		Website: "https://validator.cosmos",
		Details: "Test validator",
	}

	d2 := Description{
		Moniker:  DoNotModifyDesc,
		Identity: DoNotModifyDesc,
		Website:  DoNotModifyDesc,
		Details:  DoNotModifyDesc,
	}

	d3 := Description{
		Moniker:  "",
		Identity: "",
		Website:  "",
		Details:  "",
	}

	d, err := d1.UpdateDescription(d2)
	require.Nil(t, err)
	require.Equal(t, d, d1)

	d, err = d1.UpdateDescription(d3)
	require.Nil(t, err)
	require.Equal(t, d, d3)

	d, err = d2.UpdateDescription(d3)
	require.Nil(t, err)
	require.Equal(t, d, d3)

}
