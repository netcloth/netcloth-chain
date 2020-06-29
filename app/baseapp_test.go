package app

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/netcloth/netcloth-chain/app/protocol"
	"github.com/netcloth/netcloth-chain/app/v0/genutil"
	"github.com/netcloth/netcloth-chain/codec"
	store "github.com/netcloth/netcloth-chain/store/types"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
	"github.com/netcloth/netcloth-chain/types/module"
)

var (
	capKey1 = sdk.NewKVStoreKey("key1")
	capKey2 = sdk.NewKVStoreKey("key2")
)

func newEngine(options ...func(*MockProtocolV0)) protocol.ProtocolEngine {
	pk := sdk.NewProtocolKeeper(protocol.Keys[protocol.MainStoreKey])
	engine := protocol.NewProtocolEngine(pk)
	mockProtocolV0 := newMockProtocolV0()
	for _, option := range options {
		option(mockProtocolV0)
	}
	engine.Add(mockProtocolV0)
	engine.LoadProtocol(0)

	return engine
}

func defaultLogger() log.Logger {
	return log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
}

func newBaseApp(name string, options ...func(*BaseApp)) *BaseApp {
	logger := defaultLogger()
	db := dbm.NewMemDB()
	codec := codec.New()
	registerTestCodec(codec)

	bApp := NewBaseApp(name, logger, db, options...)
	bApp.txDecoder = testTxDecoder(codec)

	return bApp
}

func registerTestCodec(cdc *codec.Codec) {
	sdk.RegisterCodec(cdc)

	cdc.RegisterConcrete(&txTest{}, "nch/codec/baseapp/txTest", nil)
	cdc.RegisterConcrete(&msgCounter{}, "nch/baseapp/msgCounter", nil)
	cdc.RegisterConcrete(&msgCounter2{}, "nch/baseapp/msgCounter2", nil)
	cdc.RegisterConcrete(&msgNoRoute{}, "nch/baseapp/msgNoRoute", nil)
}

func setupBaseApp(t *testing.T, engine *protocol.ProtocolEngine, options ...func(*BaseApp)) *BaseApp {
	app := newBaseApp(t.Name(), options...)
	require.Equal(t, t.Name(), app.Name())

	if engine != nil {
		app.SetProtocolEngine(engine)
		app.Engine.LoadProtocol(app.Engine.GetCurrentVersion())
	}

	require.Panics(t, func() {
		app.LoadLatestVersion(capKey1)
	})

	app.MountStores(capKey1, capKey2)

	// stores are mounted
	err := app.LoadLatestVersion(capKey1)
	require.Nil(t, err)

	return app
}

func TestMountStores(t *testing.T) {
	engine := newEngine()
	app := setupBaseApp(t, &engine)

	store1 := app.cms.GetCommitKVStore(capKey1)
	require.NotNil(t, store1)
	store2 := app.cms.GetCommitKVStore(capKey2)
	require.NotNil(t, store2)
}

func setupProtocol(app *BaseApp, capKey sdk.StoreKey) {
	pk := sdk.NewProtocolKeeper(capKey)
	engine := protocol.NewProtocolEngine(pk)
	app.SetProtocolEngine(&engine)
	mockProtocolV0 := newMockProtocolV0()
	engine.Add(mockProtocolV0)
	engine.LoadProtocol(0)
}

func TestLoadVersion(t *testing.T) {
	logger := defaultLogger()
	pruningOpt := SetPruning(store.PruneSyncable)
	db := dbm.NewMemDB()
	name := t.Name()

	app := NewBaseApp(name, logger, db, pruningOpt)

	capKey := sdk.NewKVStoreKey("main")
	app.MountStores(capKey)
	setupProtocol(app, capKey)
	err := app.LoadLatestVersion(capKey) // needed to make stores non-nil
	require.Nil(t, err)

	emptyCommitID := sdk.CommitID{}

	// fresh store has zero/empty last commit
	lastHeight := app.LastBlockHeight()
	lastID := app.LastCommitID()
	require.Equal(t, int64(0), lastHeight)
	require.Equal(t, emptyCommitID, lastID)

	// execute a block, collect commit ID
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	res := app.Commit()
	commitID1 := sdk.CommitID{1, res.Data}

	// execute a block, collect commit ID
	header = abci.Header{Height: 2}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	res = app.Commit()
	commitID2 := sdk.CommitID{2, res.Data}

	// reload with LoadLatestVersion
	app = NewBaseApp(name, logger, db, pruningOpt)
	app.MountStores(capKey)
	setupProtocol(app, capKey)
	err = app.LoadLatestVersion(capKey)
	require.Nil(t, err)
	testLoadVersionHelper(t, app, int64(2), commitID2)

	// reload with LoadVersion, see if you can commit the same block and get
	// the same result
	app = NewBaseApp(name, logger, db, pruningOpt)
	app.MountStores(capKey)
	setupProtocol(app, capKey)
	err = app.LoadVersion(1, capKey)
	require.Nil(t, err)
	testLoadVersionHelper(t, app, int64(1), commitID1)
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	app.Commit()
	testLoadVersionHelper(t, app, int64(2), commitID2)
}

func TestAppVersionSetterGetter(t *testing.T) {
	logger := defaultLogger()
	pruningOpt := SetPruning(store.PruneSyncable)
	db := dbm.NewMemDB()
	name := t.Name()

	app := NewBaseApp(name, logger, db, pruningOpt)
	setupProtocol(app, capKey1)

	require.Equal(t, "", app.AppVersion())
	res := app.Query(abci.RequestQuery{Path: "app/version"})
	require.True(t, res.IsOK())
	require.Equal(t, "", string(res.Value))

	versionString := "1.0.0"
	app.SetAppVersion(versionString)
	require.Equal(t, versionString, app.AppVersion())
	res = app.Query(abci.RequestQuery{Path: "app/version"})
	require.True(t, res.IsOK())
	require.Equal(t, versionString, string(res.Value))
}

func TestLoadVersionInvalid(t *testing.T) {
	logger := log.NewNopLogger()
	pruningOpt := SetPruning(store.PruneSyncable)
	db := dbm.NewMemDB()
	name := t.Name()

	app := NewBaseApp(name, logger, db, pruningOpt)
	setupProtocol(app, capKey1)
	capKey := sdk.NewKVStoreKey("main")
	app.MountStores(capKey)
	setupProtocol(app, capKey1)
	err := app.LoadLatestVersion(capKey)
	require.Nil(t, err)

	// require error when loading an invalid version
	err = app.LoadVersion(-1, capKey)
	require.Error(t, err)

	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	res := app.Commit()
	commitID1 := sdk.CommitID{1, res.Data}

	// create a new app with the stores mounted under the same cap key
	app = NewBaseApp(name, logger, db, pruningOpt)
	app.MountStores(capKey)
	setupProtocol(app, capKey1)

	// require we can load the latest version
	err = app.LoadVersion(1, capKey)
	require.Nil(t, err)
	testLoadVersionHelper(t, app, int64(1), commitID1)

	// require error when loading an invalid version
	err = app.LoadVersion(2, capKey)
	require.Error(t, err)
}

func testLoadVersionHelper(t *testing.T, app *BaseApp, expectedHeight int64, expectedID sdk.CommitID) {
	lastHeight := app.LastBlockHeight()
	lastID := app.LastCommitID()
	require.Equal(t, expectedHeight, lastHeight)
	require.Equal(t, expectedID, lastID)
}

func TestOptionFunction(t *testing.T) {
	logger := defaultLogger()
	db := dbm.NewMemDB()
	bap := NewBaseApp("starting name", logger, db, testChangeNameHelper("new name"))
	require.Equal(t, bap.name, "new name", "BaseApp should have had name changed via option function")
}

func testChangeNameHelper(name string) func(*BaseApp) {
	return func(bap *BaseApp) {
		bap.name = name
	}
}

// Test that txs can be unmarshalled and read and that
// correct error codes are returned when not
func TestTxDecoder(t *testing.T) {
	codec := codec.New()
	registerTestCodec(codec)

	app := newBaseApp(t.Name())
	tx := newTxCounter(1, 0)
	txBytes := codec.MustMarshalBinaryLengthPrefixed(tx)

	dTx, err := app.txDecoder(txBytes)
	require.NoError(t, err)

	cTx := dTx.(txTest)
	require.Equal(t, tx.Counter, cTx.Counter)
}

// Test that Info returns the latest committed state.
func TestInfo(t *testing.T) {
	app := newBaseApp(t.Name())

	reqInfo := abci.RequestInfo{}
	res := app.Info(reqInfo)

	assert.Equal(t, "", res.Version)
	assert.Equal(t, t.Name(), res.GetData())
	assert.Equal(t, int64(0), res.LastBlockHeight)
	require.Equal(t, []uint8(nil), res.LastBlockAppHash)
}

func TestBaseAppOptionSeal(t *testing.T) {
	app := setupBaseApp(t, nil)

	require.Panics(t, func() {
		app.SetName("")
	})
	require.Panics(t, func() {
		app.SetAppVersion("")
	})
	require.Panics(t, func() {
		app.SetDB(nil)
	})
	require.Panics(t, func() {
		app.SetCMS(nil)
	})
	require.Panics(t, func() {
		app.SetAddrPeerFilter(nil)
	})
	require.Panics(t, func() {
		app.SetIDPeerFilter(nil)
	})
	require.Panics(t, func() {
		app.SetFauxMerkleMode()
	})
}

func TestSetMinGasPrices(t *testing.T) {
	minGasPrices := sdk.DecCoins{sdk.NewInt64DecCoin("stake", 5000)}
	app := newBaseApp(t.Name(), SetMinGasPrices(minGasPrices.String()))
	require.Equal(t, minGasPrices, app.minGasPrices)
}

func TestInitChainer(t *testing.T) {
	name := t.Name()
	db := dbm.NewMemDB()
	logger := defaultLogger()

	app := NewBaseApp(name, logger, db)
	capKey := sdk.NewKVStoreKey("main")
	capKey2 := sdk.NewKVStoreKey("key2")
	app.MountStores(capKey, capKey2)
	setupProtocol(app, capKey1)

	key, value := []byte("hello"), []byte("goodbye")
	var initChainer sdk.InitChainer = func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		store := ctx.KVStore(capKey)
		store.Set(key, value)
		return abci.ResponseInitChain{}
	}

	query := abci.RequestQuery{
		Path: "/store/main/key",
		Data: key,
	}

	app.InitChain(abci.RequestInitChain{})
	res := app.Query(query)
	require.Equal(t, 0, len(res.Value))

	app.Engine.GetCurrentProtocol().SetInitChainer(initChainer)

	err := app.LoadLatestVersion(capKey)
	require.Nil(t, err)
	require.Equal(t, int64(0), app.LastBlockHeight())

	app.InitChain(abci.RequestInitChain{AppStateBytes: []byte("{}"), ChainId: "test-chain-id"})

	chainID := app.deliverState.ctx.ChainID()
	require.Equal(t, "test-chain-id", chainID, "ChainID in deliverState not set correctly in InitChain")

	chainID = app.checkState.ctx.ChainID()
	require.Equal(t, "test-chain-id", chainID, "ChainID in checkState not set correctly in InitChain")

	app.Commit()
	res = app.Query(query)
	require.Equal(t, int64(1), app.LastBlockHeight())
	require.Equal(t, value, res.Value)

	// reload app
	app = NewBaseApp(name, logger, db)
	app.MountStores(capKey, capKey2)
	setupProtocol(app, capKey1)
	err = app.LoadLatestVersion(capKey) // needed to make stores non-nil
	require.Nil(t, err)
	require.Equal(t, int64(1), app.LastBlockHeight())

	// ensure we can still query after reloading
	res = app.Query(query)
	require.Equal(t, value, res.Value)

	// commit and ensure we can still query
	header := abci.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	app.Commit()

	res = app.Query(query)
	require.Equal(t, value, res.Value)
}

// Simple tx with a list of Msgs.
type txTest struct {
	Msgs       []sdk.Msg
	Counter    int64
	FailOnAnte bool
}

func (tx *txTest) setFailOnAnte(fail bool) {
	tx.FailOnAnte = fail
}

func (tx *txTest) setFailOnHandler(fail bool) {
	for i, msg := range tx.Msgs {
		tx.Msgs[i] = msgCounter{msg.(msgCounter).Counter, fail}
	}
}

// Implements Tx
func (tx txTest) GetMsgs() []sdk.Msg   { return tx.Msgs }
func (tx txTest) ValidateBasic() error { return nil }

const (
	routeMsgCounter  = "msgCounter"
	routeMsgCounter2 = "msgCounter2"
)

// ValidateBasic() fails on negative counters.
// Otherwise it's up to the handlers
type msgCounter struct {
	Counter       int64
	FailOnHandler bool
}

// Implements Msg
func (msg msgCounter) Route() string                { return routeMsgCounter }
func (msg msgCounter) Type() string                 { return "counter1" }
func (msg msgCounter) GetSignBytes() []byte         { return nil }
func (msg msgCounter) GetSigners() []sdk.AccAddress { return nil }
func (msg msgCounter) ValidateBasic() error {
	if msg.Counter >= 0 {
		return nil
	}
	return sdkerrors.Wrap(sdkerrors.ErrInvalidSequence, "counter should be a non-negative integer")
}

func newTxCounter(txInt int64, msgInts ...int64) *txTest {
	var msgs []sdk.Msg
	for _, msgInt := range msgInts {
		msgs = append(msgs, msgCounter{msgInt, false})
	}
	return &txTest{msgs, txInt, false}
}

// a msg we dont know how to route
type msgNoRoute struct {
	msgCounter
}

func (tx msgNoRoute) Route() string { return "noroute" }

// a msg we dont know how to decode
type msgNoDecode struct {
	msgCounter
}

func (tx msgNoDecode) Route() string { return routeMsgCounter }

// Another counter msg. Duplicate of msgCounter
type msgCounter2 struct {
	Counter int64
}

// Implements Msg
func (msg msgCounter2) Route() string                { return routeMsgCounter2 }
func (msg msgCounter2) Type() string                 { return "counter2" }
func (msg msgCounter2) GetSignBytes() []byte         { return nil }
func (msg msgCounter2) GetSigners() []sdk.AccAddress { return nil }
func (msg msgCounter2) ValidateBasic() error {
	if msg.Counter >= 0 {
		return nil
	}
	return sdkerrors.Wrap(sdkerrors.ErrInvalidSequence, "counter should be a non-negative integer")
}

// amino decode
func testTxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		var tx txTest
		if len(txBytes) == 0 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "tx bytes are empty")
		}
		err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
		if err != nil {
			return nil, sdkerrors.ErrTxDecode
		}
		return tx, nil
	}
}

func anteHandlerTxTest(t *testing.T, capKey *sdk.KVStoreKey, storeKey []byte) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
		store := ctx.KVStore(capKey)
		txTest := tx.(txTest)

		if txTest.FailOnAnte {
			return newCtx, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "ante handler failure")
		}

		_, err = incrementingCounter(t, store, storeKey, txTest.Counter)
		if err != nil {
			return newCtx, err
		}

		return newCtx, nil
	}
}

func handlerMsgCounter(t *testing.T, capKey *sdk.KVStoreKey, deliverKey []byte) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		store := ctx.KVStore(capKey)
		var msgCount int64

		switch m := msg.(type) {
		case *msgCounter:
			if m.FailOnHandler {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "message handler failure")
			}

			msgCount = m.Counter
		case *msgCounter2:
			msgCount = m.Counter
		}

		return incrementingCounter(t, store, deliverKey, msgCount)
	}
}

func i2b(i int64) []byte {
	return []byte{byte(i)}
}

func getIntFromStore(store sdk.KVStore, key []byte) int64 {
	bz := store.Get(key)
	if len(bz) == 0 {
		return 0
	}
	i, err := binary.ReadVarint(bytes.NewBuffer(bz))
	if err != nil {
		panic(err)
	}
	return i
}

func setIntOnStore(store sdk.KVStore, key []byte, i int64) {
	bz := make([]byte, 8)
	n := binary.PutVarint(bz, i)
	store.Set(key, bz[:n])
}

// check counter matches what's in store.
// increment and store
func incrementingCounter(t *testing.T, store sdk.KVStore, counterKey []byte, counter int64) (*sdk.Result, error) {
	storedCounter := getIntFromStore(store, counterKey)
	require.Equal(t, storedCounter, counter)
	setIntOnStore(store, counterKey, counter+1)
	return &sdk.Result{}, nil
}

//---------------------------------------------------------------------
// Tx processing - CheckTx, DeliverTx, SimulateTx.
// These tests use the serialized tx as input, while most others will use the
// Check(), Deliver(), Simulate() methods directly.
// Ensure that Check/Deliver/Simulate work as expected with the store.

// Test that successive CheckTx can see each others' effects
// on the store within a block, and that the CheckTx state
// gets reset to the latest committed state during Commit
func TestCheckTx(t *testing.T) {
	// This ante handler reads the key and checks that the value matches the current counter.
	// This ensures changes to the kvstore persist across successive CheckTx.
	counterKey := []byte("counter-key")

	anteOpt := func(p *MockProtocolV0) {
		p.SetAnteHandler(anteHandlerTxTest(t, capKey1, counterKey))
	}

	routerOpt := func(p *MockProtocolV0) {
		// TODO: can remove this once CheckTx doesnt process msgs.
		p.GetRouter().AddRoute(routeMsgCounter, func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
			return &sdk.Result{}, nil
		})
	}
	engine := newEngine(anteOpt, routerOpt)
	app := setupBaseApp(t, &engine)

	nTxs := int64(5)
	app.InitChain(abci.RequestInitChain{})

	// Create same codec used in txDecoder
	codec := codec.New()
	registerTestCodec(codec)

	for i := int64(0); i < nTxs; i++ {
		tx := newTxCounter(i, 0)
		txBytes, err := codec.MarshalBinaryLengthPrefixed(tx)
		require.NoError(t, err)
		r := app.CheckTx(abci.RequestCheckTx{Tx: txBytes})
		assert.True(t, r.IsOK(), fmt.Sprintf("%v", r))
	}

	checkStateStore := app.checkState.ctx.KVStore(capKey1)
	storedCounter := getIntFromStore(checkStateStore, counterKey)

	// Ensure AnteHandler ran
	require.Equal(t, nTxs, storedCounter)

	// If a block is committed, CheckTx state should be reset.
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()

	checkStateStore = app.checkState.ctx.KVStore(capKey1)
	storedBytes := checkStateStore.Get(counterKey)
	require.Nil(t, storedBytes)
}

// Test that successive DeliverTx can see each others' effects
// on the store, both within and across blocks.
func TestDeliverTx(t *testing.T) {
	// test increments in the ante
	anteKey := []byte("ante-key")
	anteOpt := func(p *MockProtocolV0) {
		p.SetAnteHandler(anteHandlerTxTest(t, capKey1, anteKey))
	}

	// test increments in the handler
	deliverKey := []byte("deliver-key")
	routerOpt := func(p *MockProtocolV0) {
		p.GetRouter().AddRoute(routeMsgCounter, handlerMsgCounter(t, capKey1, deliverKey))
	}

	engine := newEngine(anteOpt, routerOpt)
	app := setupBaseApp(t, &engine)
	app.InitChain(abci.RequestInitChain{})

	// Create same codec used in txDecoder
	codec := codec.New()
	registerTestCodec(codec)

	nBlocks := 3
	txPerHeight := 5

	for blockN := 0; blockN < nBlocks; blockN++ {
		header := abci.Header{Height: int64(blockN) + 1}
		app.BeginBlock(abci.RequestBeginBlock{Header: header})

		for i := 0; i < txPerHeight; i++ {
			counter := int64(blockN*txPerHeight + i)
			tx := newTxCounter(counter, counter)

			txBytes, err := codec.MarshalBinaryLengthPrefixed(tx)
			require.NoError(t, err)

			res := app.DeliverTx(abci.RequestDeliverTx{Tx: txBytes})
			require.True(t, res.IsOK(), fmt.Sprintf("%v", res))
		}

		app.EndBlock(abci.RequestEndBlock{})
		app.Commit()
	}
}

// Number of messages doesn't matter to CheckTx.
func TestMultiMsgCheckTx(t *testing.T) {
	// TODO: ensure we get the same results
	// with one message or many
}

// One call to DeliverTx should process all the messages, in order.
func TestMultiMsgDeliverTx(t *testing.T) {
	// increment the tx counter
	anteKey := []byte("ante-key")
	anteOpt := func(p *MockProtocolV0) {
		p.SetAnteHandler(anteHandlerTxTest(t, capKey1, anteKey))
	}

	// increment the msg counter
	deliverKey := []byte("deliver-key")
	deliverKey2 := []byte("deliver-key2")
	routerOpt := func(p *MockProtocolV0) {
		p.GetRouter().AddRoute(routeMsgCounter, handlerMsgCounter(t, capKey1, deliverKey))
		p.GetRouter().AddRoute(routeMsgCounter2, handlerMsgCounter(t, capKey1, deliverKey2))
	}

	engine := newEngine(anteOpt, routerOpt)
	app := setupBaseApp(t, &engine)

	// Create same codec used in txDecoder
	codec := codec.New()
	registerTestCodec(codec)

	// run a multi-msg tx
	// with all msgs the same route

	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	tx := newTxCounter(0, 0, 1, 2)
	txBytes, err := codec.MarshalBinaryLengthPrefixed(tx)
	require.NoError(t, err)
	res := app.DeliverTx(abci.RequestDeliverTx{Tx: txBytes})
	require.True(t, res.IsOK(), fmt.Sprintf("%v", res))

	store := app.deliverState.ctx.KVStore(capKey1)

	// tx counter only incremented once
	txCounter := getIntFromStore(store, anteKey)
	require.Equal(t, int64(1), txCounter)

	// msg counter incremented three times
	msgCounter := getIntFromStore(store, deliverKey)
	require.Equal(t, int64(3), msgCounter)

	// replace the second message with a msgCounter2

	tx = newTxCounter(1, 3)
	tx.Msgs = append(tx.Msgs, msgCounter2{0})
	tx.Msgs = append(tx.Msgs, msgCounter2{1})
	txBytes, err = codec.MarshalBinaryLengthPrefixed(tx)
	require.NoError(t, err)
	res = app.DeliverTx(abci.RequestDeliverTx{Tx: txBytes})
	require.True(t, res.IsOK(), fmt.Sprintf("%v", res))

	store = app.deliverState.ctx.KVStore(capKey1)

	// tx counter only incremented once
	txCounter = getIntFromStore(store, anteKey)
	require.Equal(t, int64(2), txCounter)

	// original counter increments by one
	// new counter increments by two
	msgCounter = getIntFromStore(store, deliverKey)
	require.Equal(t, int64(4), msgCounter)
	msgCounter2 := getIntFromStore(store, deliverKey2)
	require.Equal(t, int64(2), msgCounter2)
}

// Interleave calls to Check and Deliver and ensure
// that there is no cross-talk. Check sees results of the previous Check calls
// and Deliver sees that of the previous Deliver calls, but they don't see eachother.
func TestConcurrentCheckDeliver(t *testing.T) {
	// TODO
}

// Simulate a transaction that uses gas to compute the gas.
// Simulate() and Query("/app/simulate", txBytes) should give
// the same results.
func TestSimulateTx(t *testing.T) {
	gasConsumed := uint64(5)

	anteOpt := func(p *MockProtocolV0) {
		p.SetAnteHandler(func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
			newCtx = ctx.WithGasMeter(sdk.NewGasMeter(gasConsumed))
			return
		})
	}

	routerOpt := func(p *MockProtocolV0) {
		p.GetRouter().AddRoute(routeMsgCounter, func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
			ctx.GasMeter().ConsumeGas(gasConsumed, "test")
			return &sdk.Result{}, nil
		})
	}

	engine := newEngine(anteOpt, routerOpt)
	app := setupBaseApp(t, &engine)

	app.InitChain(abci.RequestInitChain{})

	// Create same codec used in txDecoder
	cdc := codec.New()
	registerTestCodec(cdc)

	nBlocks := 3
	for blockN := 0; blockN < nBlocks; blockN++ {
		count := int64(blockN + 1)
		header := abci.Header{Height: count}
		app.BeginBlock(abci.RequestBeginBlock{Header: header})

		tx := newTxCounter(count, count)
		txBytes, err := cdc.MarshalBinaryLengthPrefixed(tx)
		require.Nil(t, err)

		// simulate a message, check gas reported
		gInfo, result, err := app.Simulate(txBytes, tx)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, gasConsumed, gInfo.GasUsed)

		// simulate again, same result
		gInfo, result, err = app.Simulate(txBytes, tx)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, gasConsumed, gInfo.GasUsed)

		// simulate by calling Query with encoded tx
		query := abci.RequestQuery{
			Path: "/app/simulate",
			Data: txBytes,
		}
		queryResult := app.Query(query)
		require.True(t, queryResult.IsOK(), queryResult.Log)

		var res uint64
		err = codec.Cdc.UnmarshalBinaryLengthPrefixed(queryResult.Value, &res)
		require.NoError(t, err)
		require.Equal(t, gasConsumed, res)
		app.EndBlock(abci.RequestEndBlock{})
		app.Commit()
	}
}

func TestRunInvalidTransaction(t *testing.T) {
	anteOpt := func(p *MockProtocolV0) {
		p.SetAnteHandler(func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
			return
		})
	}
	routerOpt := func(p *MockProtocolV0) {
		p.GetRouter().AddRoute(routeMsgCounter, func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
			return &sdk.Result{}, nil
		})
	}

	engine := newEngine(anteOpt, routerOpt)
	app := setupBaseApp(t, &engine)

	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// Transaction with no messages
	{
		emptyTx := &txTest{}
		_, result, err := app.Deliver(emptyTx)
		require.Error(t, err)
		require.Nil(t, result)

		space, code, _ := sdkerrors.ABCIInfo(err, false)
		require.EqualValues(t, sdkerrors.ErrInvalidRequest.Codespace(), space, err)
		require.EqualValues(t, sdkerrors.ErrInvalidRequest.ABCICode(), code, err)
	}

	// Transaction where ValidateBasic fails
	{
		testCases := []struct {
			tx   *txTest
			fail bool
		}{
			{newTxCounter(0, 0), false},
			{newTxCounter(-1, 0), false},
			{newTxCounter(100, 100), false},
			{newTxCounter(100, 5, 4, 3, 2, 1), false},

			{newTxCounter(0, -1), true},
			{newTxCounter(0, 1, -2), true},
			{newTxCounter(0, 1, 2, -10, 5), true},
		}

		for _, testCase := range testCases {
			tx := testCase.tx
			_, result, err := app.Deliver(tx)

			if testCase.fail {
				require.Error(t, err)

				space, code, _ := sdkerrors.ABCIInfo(err, false)
				require.EqualValues(t, sdkerrors.ErrInvalidSequence.Codespace(), space, err)
				require.EqualValues(t, sdkerrors.ErrInvalidSequence.ABCICode(), code, err)
			} else {
				require.NotNil(t, result)
			}
		}
	}

	// Transaction with no known route
	{
		unknownRouteTx := txTest{[]sdk.Msg{msgNoRoute{}}, 0, false}
		_, result, err := app.Deliver(unknownRouteTx)
		require.Error(t, err)
		require.Nil(t, result)

		space, code, _ := sdkerrors.ABCIInfo(err, false)
		require.EqualValues(t, sdkerrors.ErrUnknownRequest.Codespace(), space, err)
		require.EqualValues(t, sdkerrors.ErrUnknownRequest.ABCICode(), code, err)

		unknownRouteTx = txTest{[]sdk.Msg{msgCounter{}, msgNoRoute{}}, 0, false}
		_, result, err = app.Deliver(unknownRouteTx)
		require.Error(t, err)
		require.Nil(t, result)

		space, code, _ = sdkerrors.ABCIInfo(err, false)
		require.EqualValues(t, sdkerrors.ErrUnknownRequest.Codespace(), space, err)
		require.EqualValues(t, sdkerrors.ErrUnknownRequest.ABCICode(), code, err)
	}

	// Transaction with an unregistered message
	{
		tx := newTxCounter(0, 0)
		tx.Msgs = append(tx.Msgs, msgNoDecode{})

		// new codec so we can encode the tx, but we shouldn't be able to decode
		newCdc := codec.New()
		registerTestCodec(newCdc)
		newCdc.RegisterConcrete(&msgNoDecode{}, "nch/baseapp/msgNoDecode", nil)

		txBytes, err := newCdc.MarshalBinaryLengthPrefixed(tx)
		require.NoError(t, err)
		res := app.DeliverTx(abci.RequestDeliverTx{Tx: txBytes})
		require.EqualValues(t, sdkerrors.ErrTxDecode.ABCICode(), res.Code)
		require.EqualValues(t, sdkerrors.ErrTxDecode.Codespace(), res.Codespace)
	}
}

// Test that transactions exceeding gas limits fail
func TestTxGasLimits(t *testing.T) {
	gasGranted := uint64(10)
	anteOpt := func(p *MockProtocolV0) {
		p.SetAnteHandler(func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
			newCtx = ctx.WithGasMeter(sdk.NewGasMeter(gasGranted))

			defer func() {
				if r := recover(); r != nil {
					switch rType := r.(type) {
					case sdk.ErrorOutOfGas:
						err = sdkerrors.Wrapf(sdkerrors.ErrOutOfGas, "out of gas in location: %v", rType.Descriptor)
					default:
						panic(r)
					}
				}
			}()

			count := tx.(*txTest).Counter
			newCtx.GasMeter().ConsumeGas(uint64(count), "counter-ante")

			return newCtx, nil
		})

	}

	routerOpt := func(p *MockProtocolV0) {
		p.GetRouter().AddRoute(routeMsgCounter, func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
			count := msg.(msgCounter).Counter
			ctx.GasMeter().ConsumeGas(uint64(count), "counter-handler")
			return &sdk.Result{}, nil
		})
	}

	engine := newEngine(anteOpt, routerOpt)
	app := setupBaseApp(t, &engine)

	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	testCases := []struct {
		tx      *txTest
		gasUsed uint64
		fail    bool
	}{
		{newTxCounter(0, 0), 0, false},
		{newTxCounter(1, 1), 2, false},
		{newTxCounter(9, 1), 10, false},
		{newTxCounter(1, 9), 10, false},
		{newTxCounter(10, 0), 10, false},
		{newTxCounter(0, 10), 10, false},
		{newTxCounter(0, 8, 2), 10, false},
		{newTxCounter(0, 5, 1, 1, 1, 1, 1), 10, false},
		{newTxCounter(0, 5, 1, 1, 1, 1), 9, false},

		{newTxCounter(9, 2), 11, true},
		{newTxCounter(2, 9), 11, true},
		{newTxCounter(9, 1, 1), 11, true},
		{newTxCounter(1, 8, 1, 1), 11, true},
		{newTxCounter(11, 0), 11, true},
		{newTxCounter(0, 11), 11, true},
		{newTxCounter(0, 5, 11), 16, true},
	}

	for i, tc := range testCases {
		tx := tc.tx
		gInfo, result, err := app.Deliver(tx)

		// check gas used and wanted
		require.Equal(t, tc.gasUsed, gInfo.GasUsed, fmt.Sprintf("tc #%d; gas: %v, result: %v, err: %s", i, gInfo, result, err))

		// check for out of gas
		if !tc.fail {
			require.NotNil(t, result, fmt.Sprintf("%d: %v, %v", i, tc, err))
		} else {
			require.Error(t, err)
			require.Nil(t, result)

			space, code, _ := sdkerrors.ABCIInfo(err, false)
			require.EqualValues(t, sdkerrors.ErrOutOfGas.Codespace(), space, err)
			require.EqualValues(t, sdkerrors.ErrOutOfGas.ABCICode(), code, err)
		}
	}
}

// Test that transactions exceeding gas limits fail
func TestMaxBlockGasLimits(t *testing.T) {
	gasGranted := uint64(10)
	anteOpt := func(p *MockProtocolV0) {
		p.SetAnteHandler(func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
			newCtx = ctx.WithGasMeter(sdk.NewGasMeter(gasGranted))

			defer func() {
				if r := recover(); r != nil {
					switch rType := r.(type) {
					case sdk.ErrorOutOfGas:
						err = sdkerrors.Wrapf(sdkerrors.ErrOutOfGas, "out of gas in location: %v", rType.Descriptor)
					default:
						panic(r)
					}
				}
			}()

			count := tx.(*txTest).Counter
			newCtx.GasMeter().ConsumeGas(uint64(count), "counter-ante")

			return
		})

	}

	routerOpt := func(p *MockProtocolV0) {
		p.GetRouter().AddRoute(routeMsgCounter, func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
			count := msg.(msgCounter).Counter
			ctx.GasMeter().ConsumeGas(uint64(count), "counter-handler")
			return &sdk.Result{}, nil
		})
	}

	engine := newEngine(anteOpt, routerOpt)
	app := setupBaseApp(t, &engine)

	app.InitChain(abci.RequestInitChain{
		ConsensusParams: &abci.ConsensusParams{
			Block: &abci.BlockParams{
				MaxGas: 100,
			},
		},
	})

	testCases := []struct {
		tx                *txTest
		numDelivers       int
		gasUsedPerDeliver uint64
		fail              bool
		failAfterDeliver  int
	}{
		{newTxCounter(0, 0), 0, 0, false, 0},
		{newTxCounter(9, 1), 2, 10, false, 0},
		{newTxCounter(10, 0), 3, 10, false, 0},
		{newTxCounter(10, 0), 10, 10, false, 0},
		{newTxCounter(2, 7), 11, 9, false, 0},
		{newTxCounter(10, 0), 10, 10, false, 0}, // hit the limit but pass

		{newTxCounter(10, 0), 11, 10, true, 10},
		{newTxCounter(10, 0), 15, 10, true, 10},
		{newTxCounter(9, 0), 12, 9, true, 11}, // fly past the limit
	}

	for i, tc := range testCases {
		fmt.Printf("debug i: %v\n", i)
		tx := tc.tx

		// reset the block gas
		header := abci.Header{Height: app.LastBlockHeight() + 1}
		app.BeginBlock(abci.RequestBeginBlock{Header: header})

		// execute the transaction multiple times
		for j := 0; j < tc.numDelivers; j++ {
			_, result, err := app.Deliver(tx)

			ctx := app.getState(runTxModeDeliver).ctx

			// check for failed transactions
			if tc.fail && (j+1) > tc.failAfterDeliver {
				require.Error(t, err, fmt.Sprintf("tc #%d; result: %v, err: %s", i, result, err))
				require.Nil(t, result, fmt.Sprintf("tc #%d; result: %v, err: %s", i, result, err))

				space, code, _ := sdkerrors.ABCIInfo(err, false)
				require.EqualValues(t, sdkerrors.ErrOutOfGas.Codespace(), space, err)
				require.EqualValues(t, sdkerrors.ErrOutOfGas.ABCICode(), code, err)
				require.True(t, ctx.BlockGasMeter().IsOutOfGas())
			} else {
				// check gas used and wanted
				blockGasUsed := ctx.BlockGasMeter().GasConsumed()
				expBlockGasUsed := tc.gasUsedPerDeliver * uint64(j+1)
				require.Equal(
					t, expBlockGasUsed, blockGasUsed,
					fmt.Sprintf("%d,%d: %v, %v, %v, %v", i, j, tc, expBlockGasUsed, blockGasUsed, result),
				)

				require.NotNil(t, result, fmt.Sprintf("tc #%d; currDeliver: %d, result: %v, err: %s", i, j, result, err))
				require.False(t, ctx.BlockGasMeter().IsPastLimit())
			}
		}
	}
}

func TestBaseAppAnteHandler(t *testing.T) {
	anteKey := []byte("ante-key")
	anteOpt := func(p *MockProtocolV0) {
		p.SetAnteHandler(anteHandlerTxTest(t, capKey1, anteKey))
	}

	deliverKey := []byte("deliver-key")
	routerOpt := func(p *MockProtocolV0) {
		p.GetRouter().AddRoute(routeMsgCounter, handlerMsgCounter(t, capKey1, deliverKey))
	}

	cdc := codec.New()
	engine := newEngine(anteOpt, routerOpt)
	app := setupBaseApp(t, &engine)

	app.InitChain(abci.RequestInitChain{})
	registerTestCodec(cdc)

	header := abci.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// execute a tx that will fail ante handler execution
	//
	// NOTE: State should not be mutated here. This will be implicitly checked by
	// the next txs ante handler execution (anteHandlerTxTest).
	tx := newTxCounter(0, 0)
	tx.setFailOnAnte(true)
	txBytes, err := cdc.MarshalBinaryLengthPrefixed(tx)
	require.NoError(t, err)
	res := app.DeliverTx(abci.RequestDeliverTx{Tx: txBytes})
	require.False(t, res.IsOK(), fmt.Sprintf("%v", res))

	ctx := app.getState(runTxModeDeliver).ctx
	store := ctx.KVStore(capKey1)
	require.Equal(t, int64(0), getIntFromStore(store, anteKey))

	// execute at tx that will pass the ante handler (the checkTx state should
	// mutate) but will fail the message handler
	tx = newTxCounter(0, 0)
	tx.setFailOnHandler(true)

	txBytes, err = cdc.MarshalBinaryLengthPrefixed(tx)
	require.NoError(t, err)

	res = app.DeliverTx(abci.RequestDeliverTx{Tx: txBytes})
	require.False(t, res.IsOK(), fmt.Sprintf("%v", res))

	ctx = app.getState(runTxModeDeliver).ctx
	store = ctx.KVStore(capKey1)
	require.Equal(t, int64(1), getIntFromStore(store, anteKey))
	require.Equal(t, int64(0), getIntFromStore(store, deliverKey))

	// execute a successful ante handler and message execution where state is
	// implicitly checked by previous tx executions
	tx = newTxCounter(1, 0)

	txBytes, err = cdc.MarshalBinaryLengthPrefixed(tx)
	require.NoError(t, err)

	res = app.DeliverTx(abci.RequestDeliverTx{Tx: txBytes})
	require.True(t, res.IsOK(), fmt.Sprintf("%v", res))

	ctx = app.getState(runTxModeDeliver).ctx
	store = ctx.KVStore(capKey1)
	require.Equal(t, int64(2), getIntFromStore(store, anteKey))
	require.Equal(t, int64(1), getIntFromStore(store, deliverKey))

	// commit
	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()
}

func TestGasConsumptionBadTx(t *testing.T) {
	gasWanted := uint64(5)
	anteOpt := func(p *MockProtocolV0) {
		p.SetAnteHandler(func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
			newCtx = ctx.WithGasMeter(sdk.NewGasMeter(gasWanted))

			defer func() {
				if r := recover(); r != nil {
					switch rType := r.(type) {
					case sdk.ErrorOutOfGas:
						log := fmt.Sprintf("out of gas in location: %v", rType.Descriptor)
						err = sdkerrors.Wrap(sdkerrors.ErrOutOfGas, log)
					default:
						panic(r)
					}
				}
			}()

			txTest := tx.(txTest)
			newCtx.GasMeter().ConsumeGas(uint64(txTest.Counter), "counter-ante")
			if txTest.FailOnAnte {
				return newCtx, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "ante handler failure")
			}

			return
		})
	}

	routerOpt := func(p *MockProtocolV0) {
		p.GetRouter().AddRoute(routeMsgCounter, func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
			count := msg.(msgCounter).Counter
			ctx.GasMeter().ConsumeGas(uint64(count), "counter-handler")
			return &sdk.Result{}, nil
		})
	}

	cdc := codec.New()
	registerTestCodec(cdc)

	engine := newEngine(anteOpt, routerOpt)
	app := setupBaseApp(t, &engine)

	app.InitChain(abci.RequestInitChain{
		ConsensusParams: &abci.ConsensusParams{
			Block: &abci.BlockParams{
				MaxGas: 9,
			},
		},
	})

	app.InitChain(abci.RequestInitChain{})

	header := abci.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	tx := newTxCounter(5, 0)
	tx.setFailOnAnte(true)
	txBytes, err := cdc.MarshalBinaryLengthPrefixed(tx)
	require.NoError(t, err)

	res := app.DeliverTx(abci.RequestDeliverTx{Tx: txBytes})
	require.False(t, res.IsOK(), fmt.Sprintf("%v", res))

	// require next tx to fail due to black gas limit
	tx = newTxCounter(5, 0)
	txBytes, err = cdc.MarshalBinaryLengthPrefixed(tx)
	require.NoError(t, err)

	res = app.DeliverTx(abci.RequestDeliverTx{Tx: txBytes})
	require.False(t, res.IsOK(), fmt.Sprintf("%v", res))
}

// Test that we can only query from the latest committed state.
func TestQuery(t *testing.T) {
	key, value := []byte("hello"), []byte("goodbye")
	anteOpt := func(p *MockProtocolV0) {
		p.SetAnteHandler(func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
			store := ctx.KVStore(capKey1)
			store.Set(key, value)
			return
		})
	}

	routerOpt := func(p *MockProtocolV0) {
		p.GetRouter().AddRoute(routeMsgCounter, func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
			store := ctx.KVStore(capKey1)
			store.Set(key, value)
			return &sdk.Result{}, nil
		})
	}

	engine := newEngine(anteOpt, routerOpt)
	app := setupBaseApp(t, &engine)

	app.InitChain(abci.RequestInitChain{})

	// NOTE: "/store/key1" tells us KVStore
	// and the final "/key" says to use the data as the
	// key in the given KVStore ...
	query := abci.RequestQuery{
		Path: "/store/key1/key",
		Data: key,
	}
	tx := newTxCounter(0, 0)

	// query is empty before we do anything
	res := app.Query(query)
	require.Equal(t, 0, len(res.Value))

	// query is still empty after a CheckTx
	_, resTx, err := app.Check(tx)
	require.NoError(t, err)
	require.NotNil(t, resTx)
	res = app.Query(query)
	require.Equal(t, 0, len(res.Value))

	// query is still empty after a DeliverTx before we commit
	header := abci.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	_, resTx, err = app.Deliver(tx)
	require.NoError(t, err)
	require.NotNil(t, resTx)
	res = app.Query(query)
	require.Equal(t, 0, len(res.Value))

	// query returns correct value after Commit
	app.Commit()
	res = app.Query(query)
	require.Equal(t, value, res.Value)
}

// Test p2p filter queries
func TestP2PQuery(t *testing.T) {
	addrPeerFilterOpt := func(bapp *BaseApp) {
		bapp.SetAddrPeerFilter(func(addrport string) abci.ResponseQuery {
			require.Equal(t, "1.1.1.1:8000", addrport)
			return abci.ResponseQuery{Code: uint32(3)}
		})
	}

	idPeerFilterOpt := func(bapp *BaseApp) {
		bapp.SetIDPeerFilter(func(id string) abci.ResponseQuery {
			require.Equal(t, "testid", id)
			return abci.ResponseQuery{Code: uint32(4)}
		})
	}

	app := setupBaseApp(t, nil, addrPeerFilterOpt, idPeerFilterOpt)

	addrQuery := abci.RequestQuery{
		Path: "/p2p/filter/addr/1.1.1.1:8000",
	}
	res := app.Query(addrQuery)
	require.Equal(t, uint32(3), res.Code)

	idQuery := abci.RequestQuery{
		Path: "/p2p/filter/id/testid",
	}
	res = app.Query(idQuery)
	require.Equal(t, uint32(4), res.Code)
}

func TestGetMaximumBlockGas(t *testing.T) {
	app := setupBaseApp(t, nil)

	app.setConsensusParams(&abci.ConsensusParams{Block: &abci.BlockParams{MaxGas: 0}})
	require.Equal(t, uint64(0), app.getMaximumBlockGas())

	app.setConsensusParams(&abci.ConsensusParams{Block: &abci.BlockParams{MaxGas: -1}})
	require.Equal(t, uint64(0), app.getMaximumBlockGas())

	app.setConsensusParams(&abci.ConsensusParams{Block: &abci.BlockParams{MaxGas: 5000000}})
	require.Equal(t, uint64(5000000), app.getMaximumBlockGas())

	app.setConsensusParams(&abci.ConsensusParams{Block: &abci.BlockParams{MaxGas: -5000000}})
	require.Panics(t, func() { app.getMaximumBlockGas() })
}

// NOTE: represents a new custom router for testing purposes of WithRouter()
type testCustomRouter struct {
	routes sync.Map
}

func (rtr *testCustomRouter) AddRoute(path string, h sdk.Handler) sdk.Router {
	rtr.routes.Store(path, h)
	return rtr
}

func (rtr *testCustomRouter) Route(ctx sdk.Context, path string) sdk.Handler {
	if v, ok := rtr.routes.Load(path); ok {
		if h, ok := v.(sdk.Handler); ok {
			return h
		}
	}
	return nil
}

func TestWithRouter(t *testing.T) {
	anteKey := []byte("ante-key")
	anteOpt := func(p *MockProtocolV0) {
		p.SetAnteHandler(anteHandlerTxTest(t, capKey1, anteKey))
	}

	deliverKey := []byte("deliver-key")
	routerOpt := func(p *MockProtocolV0) {
		p.SetRouter(&testCustomRouter{routes: sync.Map{}})
		p.GetRouter().AddRoute(routeMsgCounter, handlerMsgCounter(t, capKey1, deliverKey))
	}

	engine := newEngine(anteOpt, routerOpt)
	app := setupBaseApp(t, &engine)

	app.InitChain(abci.RequestInitChain{})

	codec := codec.New()
	registerTestCodec(codec)

	nBlocks := 3
	txPerHeight := 5

	for blockN := 0; blockN < nBlocks; blockN++ {
		header := abci.Header{Height: int64(blockN) + 1}
		app.BeginBlock(abci.RequestBeginBlock{Header: header})

		for i := 0; i < txPerHeight; i++ {
			counter := int64(blockN*txPerHeight + i)
			tx := newTxCounter(counter, counter)

			txBytes, err := codec.MarshalBinaryLengthPrefixed(tx)
			require.NoError(t, err)

			res := app.DeliverTx(abci.RequestDeliverTx{Tx: txBytes})
			require.True(t, res.IsOK(), fmt.Sprintf("%v", res))
		}

		app.EndBlock(abci.RequestEndBlock{})
		app.Commit()
	}
}

type MockProtocolV0 struct {
	moduleManager *module.Manager

	router      sdk.Router
	queryRouter sdk.QueryRouter

	anteHandler      sdk.AnteHandler
	feeRefundHandler sdk.FeeRefundHandler
	initChainer      sdk.InitChainer
	beginBlocker     sdk.BeginBlocker
	endBlocker       sdk.EndBlocker
	deliverTx        genutil.DeliverTxfn
}

func newMockProtocolV0() *MockProtocolV0 {
	return &MockProtocolV0{
		router:        protocol.NewRouter(),
		queryRouter:   protocol.NewQueryRouter(),
		moduleManager: module.NewManager(),
	}
}

var _ protocol.Protocol = &MockProtocolV0{}

func (m *MockProtocolV0) GetVersion() uint64 {
	return 0
}

func (m *MockProtocolV0) GetRouter() sdk.Router {
	return m.router
}

func (m *MockProtocolV0) GetQueryRouter() sdk.QueryRouter {
	return m.queryRouter
}

func (m MockProtocolV0) GetAnteHandler() sdk.AnteHandler {
	return m.anteHandler
}

func (m *MockProtocolV0) GetFeeRefundHandler() sdk.FeeRefundHandler {
	return m.feeRefundHandler
}

func (m *MockProtocolV0) GetInitChainer() sdk.InitChainer {
	return m.initChainer
}

func (m MockProtocolV0) GetBeginBlocker() sdk.BeginBlocker {
	return m.beginBlocker
}

func (m MockProtocolV0) GetEndBlocker() sdk.EndBlocker {
	return m.endBlocker
}

func (m *MockProtocolV0) ExportAppStateAndValidators(ctx sdk.Context, forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, validators []types.GenesisValidator, err error) {
	return json.RawMessage{}, nil, nil
}

func (m *MockProtocolV0) Load() {
}

func (m *MockProtocolV0) Init(ctx sdk.Context) {
}

func (m *MockProtocolV0) GetCodec() *codec.Codec {
	return codec.New()
}

func (m *MockProtocolV0) SetRouter(router sdk.Router) {
	m.router = router
}

func (m *MockProtocolV0) SetQuearyRouter(queryRouter sdk.QueryRouter) {
	m.queryRouter = queryRouter
}

func (m *MockProtocolV0) SetAnteHandler(anteHandler sdk.AnteHandler) {
	m.anteHandler = anteHandler
}

func (m *MockProtocolV0) SetInitChainer(initChainer sdk.InitChainer) {
	m.initChainer = initChainer
}
