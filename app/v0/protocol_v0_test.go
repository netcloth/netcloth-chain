package v0

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/netcloth/netcloth-chain/app/v0/mock"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func TestProtocolV0(t *testing.T) {
	mainKey := sdk.NewKVStoreKey("main")
	protocolKeeper := sdk.NewProtocolKeeper(mainKey)
	mockApp := mock.NewNCHApp()

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("test", "protocol_v0_test")
	protocolV0 := NewProtocolV0(1, logger, protocolKeeper, mockApp.DeliverTx, 0, nil)

	require.Panics(t, func() {
		protocolV0.GetCodec()
	})

	protocolV0.LoadContext()
	require.NotEqual(t, 0, len(protocolV0.moduleManager.Modules))
	require.NotEqual(t, nil, protocolV0.GetCodec())
}
