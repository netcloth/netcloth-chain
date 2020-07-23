package protocol

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	db "github.com/tendermint/tm-db"

	"github.com/netcloth/netcloth-chain/store"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestProtocolEngine(t *testing.T) {
	ctx, mainKey := createEngineTestInput(t)
	protocolKeeper := sdk.NewProtocolKeeper(mainKey)
	engine := NewProtocolEngine(protocolKeeper)
	require.NotEqual(t, nil, engine.GetProtocolKeeper())

	//check app upgrade config
	appUpgradeConfig := sdk.NewUpgradeConfig(0, sdk.NewProtocolDefinition(0, "NCH", 1024, sdk.NewDec(0)))
	protocolKeeper.SetUpgradeConfig(ctx, appUpgradeConfig)
	auc, ok := engine.GetUpgradeConfigByStore(ctx.KVStore(mainKey))
	require.Equal(t, true, ok)
	require.Equal(t, "NCH", auc.Protocol.Software)

	// add protocol randomly
	num := rand.Intn(3) + 1
	for i := 0; i < num; i++ {
		engine.Add(NewMockProtocol(uint64(i)))
		require.Equal(t, true, engine.Activate(uint64(i)))
		protocolKeeper.SetCurrentVersion(ctx, uint64(i))
	}

	currentProtocol := engine.GetCurrentProtocol()
	ok, currentVersionFromStore := engine.LoadCurrentProtocol(ctx.KVStore(mainKey))
	require.Equal(t, true, ok)
	require.Equal(t, currentVersionFromStore, currentProtocol.GetVersion(), engine.GetCurrentVersion())
}

func TestEnginePanics(t *testing.T) {
	_, mainKey := createEngineTestInput(t)
	protocolKeeper := sdk.NewProtocolKeeper(mainKey)
	engine := NewProtocolEngine(protocolKeeper)
	require.NotEqual(t, nil, engine.GetProtocolKeeper())

	// engine.next==0 && protocol.version==1 panics
	testProtocol := NewMockProtocol(1)
	require.Panics(t, func() {
		engine.Add(testProtocol)
	})

	//no protocol v1 in the engine, panics
	engine.current = uint64(1)
	require.Panics(t, func() {
		engine.GetCurrentProtocol()
	})
}

func createEngineTestInput(t *testing.T) (sdk.Context, *sdk.KVStoreKey) {
	keyMain := sdk.NewKVStoreKey("main")

	db := db.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyMain, sdk.StoreTypeIAVL, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)
	ctx := sdk.NewContext(ms, abci.Header{}, false, nil)

	return ctx, keyMain

}
