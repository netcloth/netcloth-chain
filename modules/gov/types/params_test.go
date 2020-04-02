package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParamsEqual(t *testing.T) {
	p1 := DefaultParams()
	p2 := DefaultParams()

	ok := p1.Equal(p2)
	require.True(t, ok)

	p2.VotingParams.VotingPeriod = time.Hour * 24 * 1

	ok = p1.Equal(p2)
	require.False(t, ok)
}
