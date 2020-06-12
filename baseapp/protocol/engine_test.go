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
