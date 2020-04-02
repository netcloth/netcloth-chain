package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParamsEqual(t *testing.T) {
	p1 := DefaultParams()
	p2 := DefaultParams()

	ok := p1.Equal(p2)
	require.True(t, ok)

	p2.BlocksPerYear = uint64(10 * 60 * 8766 / 5)

	ok = p1.Equal(p2)
	require.False(t, ok)
}
