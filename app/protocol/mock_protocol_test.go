package protocol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMockProtocol(t *testing.T) {
	var version uint64 = 1
	mockProtocol := NewMockProtocol(version)

	require.Equal(t, version, mockProtocol.GetVersion())
	require.NotPanics(t, func() {
		mockProtocol.LoadContext()
	})
	require.NotPanics(t, func() {
		mockProtocol.Init()
	})

	require.NotNil(t, mockProtocol.GetRouter())
	require.NotNil(t, mockProtocol.GetQueryRouter())

	require.Nil(t, mockProtocol.GetAnteHandler())
	require.Nil(t, mockProtocol.GetInitChainer())
	require.Nil(t, mockProtocol.GetBeginBlocker())
	require.Nil(t, mockProtocol.GetEndBlocker())

	require.Nil(t, mockProtocol.GetSimulationManager())

}
