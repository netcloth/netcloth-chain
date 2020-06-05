package baseapp

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/gogo/protobuf/proto"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/netcloth/netcloth-chain/baseapp/protocol"
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/store"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

type runTxMode uint8

const (
	runTxModeCheck    runTxMode = iota // Check a transaction
	runTxModeReCheck                   // Recheck a (pending) transaction after a commit
	runTxModeSimulate                  // Simulate a transaction
	runTxModeDeliver                   // Deliver a transaction
)

var (
	_ abci.Application = (*BaseApp)(nil)

	// mainConsensusParamsKey defines a key to store the consensus params in the
	// main store.
	mainConsensusParamsKey = []byte("consensus_params")
)

// BaseApp reflects the ABCI application implementation.
type BaseApp struct {
	// initialized on creation
	logger    log.Logger
	name      string               // application name from abci.Info
	db        dbm.DB               // common DB backend
	cms       sdk.CommitMultiStore // Main (uncached) state
	TxDecoder sdk.TxDecoder        // unmarshal []byte into sdk.Tx

	Engine *protocol.ProtocolEngine

	// set upon LoadVersion or LoadLatestVersion.
	baseKey *sdk.KVStoreKey // Main KVStore in cms

	addrPeerFilter sdk.PeerFilter // filter peers by address and port
	idPeerFilter   sdk.PeerFilter // filter peers by node ID
	fauxMerkleMode bool           // if true, IAVL MountStores uses MountStoresDB for simulation speed.

	// --------------------
	// Volatile state
	// checkState is set on initialization and reset on Commit.
	// deliverState is set in InitChain and BeginBlock and cleared on Commit.
	// See methods setCheckState and setDeliverState.
	checkState   *state // for CheckTx
	deliverState *state // for DeliverTx

	// absent validators from begin block
	voteInfos []abci.VoteInfo // absent validators from begin block

	// consensus params
	// TODO: Move this in the future to baseapp param store on main store.
	consensusParams *abci.ConsensusParams

	// The minimum gas prices a validator is willing to accept for processing a
	// transaction. This is mainly used for DoS and spam prevention.
	minGasPrices sdk.DecCoins

	// flag for sealing options and parameters to a BaseApp
	sealed bool

	// height at which to halt the chain and gracefully shutdown
	haltHeight uint64

	// application's version string
	appVersion string
}

func NewBaseApp(name string, logger log.Logger, db dbm.DB, options ...func(*BaseApp)) *BaseApp { //TODO fixme crash if options use protocol instance(nil)
	app := &BaseApp{
		logger:         logger,
		name:           name,
		db:             db,
		cms:            store.NewCommitMultiStore(db),
		fauxMerkleMode: false,
	}

	for _, option := range options {
		option(app)
	}

	return app
}

func (app *BaseApp) WithEngine(engine *protocol.ProtocolEngine) *BaseApp {
	app.Engine = engine
	return app
}

func (app *BaseApp) Name() string {
	return app.name
}

func (app *BaseApp) AppVersion() string {
	return app.appVersion
}

func (app *BaseApp) Logger() log.Logger {
	return app.logger
}

// MountStores mounts all IAVL or DB stores to the provided keys in the BaseApp
// multistore.
func (app *BaseApp) MountStores(keys ...sdk.StoreKey) {
	for _, key := range keys {
		switch key.(type) {
		case *sdk.KVStoreKey:
			if !app.fauxMerkleMode {
				app.MountStore(key, sdk.StoreTypeIAVL)
			} else {
				// StoreTypeDB doesn't do anything upon commit, and it doesn't
				// retain history, but it's useful for faster simulation.
				app.MountStore(key, sdk.StoreTypeDB)
			}
		case *sdk.TransientStoreKey:
			app.MountStore(key, sdk.StoreTypeTransient)
		default:
			panic("Unrecognized store key type " + reflect.TypeOf(key).Name())
		}
	}
}

// MountStores mounts all IAVL or DB stores to the provided keys in the BaseApp
// multistore.
func (app *BaseApp) MountKVStores(keys map[string]*sdk.KVStoreKey) {
	for _, key := range keys {
		if !app.fauxMerkleMode {
			app.MountStore(key, sdk.StoreTypeIAVL)
		} else {
			// StoreTypeDB doesn't do anything upon commit, and it doesn't
			// retain history, but it's useful for faster simulation.
			app.MountStore(key, sdk.StoreTypeDB)
		}
	}
}

// MountStores mounts all IAVL or DB stores to the provided keys in the BaseApp
// multistore.
func (app *BaseApp) MountTransientStores(keys map[string]*sdk.TransientStoreKey) {
	for _, key := range keys {
		app.MountStore(key, sdk.StoreTypeTransient)
	}
}

// MountStoreWithDB mounts a store to the provided key in the BaseApp
// multistore, using a specified DB.
func (app *BaseApp) MountStoreWithDB(key sdk.StoreKey, typ sdk.StoreType, db dbm.DB) {
	app.cms.MountStoreWithDB(key, typ, db)
}

// MountStore mounts a store to the provided key in the BaseApp multistore,
// using the default DB.
func (app *BaseApp) MountStore(key sdk.StoreKey, typ sdk.StoreType) {
	app.cms.MountStoreWithDB(key, typ, nil)
}

// LoadLatestVersion loads the latest application version. It will panic if
// called more than once on a running BaseApp.
func (app *BaseApp) LoadLatestVersion(baseKey *sdk.KVStoreKey) error {
	err := app.cms.LoadLatestVersion()
	if err != nil {
		return err
	}
	return app.initFromMainStore(baseKey)
}

// LoadVersion loads the BaseApp application version. It will panic if called
// more than once on a running baseapp.
func (app *BaseApp) LoadVersion(version int64, baseKey *sdk.KVStoreKey) error {
	err := app.cms.LoadVersion(version)
	if err != nil {
		return err
	}
	return app.initFromMainStore(baseKey)
}

// LastCommitID returns the last CommitID of the multistore.
func (app *BaseApp) LastCommitID() sdk.CommitID {
	return app.cms.LastCommitID()
}

// LastBlockHeight returns the last committed block height.
func (app *BaseApp) LastBlockHeight() int64 {
	return app.cms.LastCommitID().Version
}

// initializes the remaining logic from app.cms
func (app *BaseApp) initFromMainStore(baseKey *sdk.KVStoreKey) error {
	mainStore := app.cms.GetKVStore(baseKey)
	if mainStore == nil {
		return errors.New("baseapp expects MultiStore with 'main' KVStore")
	}

	// memoize baseKey
	if app.baseKey != nil {
		panic("app.baseKey expected to be nil; duplicate init?")
	}
	app.baseKey = baseKey

	// Load the consensus params from the main store. If the consensus params are
	// nil, it will be saved later during InitChain.
	//
	// TODO: assert that InitChain hasn't yet been called.
	consensusParamsBz := mainStore.Get(mainConsensusParamsKey)
	if consensusParamsBz != nil {
		var consensusParams = &abci.ConsensusParams{}

		err := proto.Unmarshal(consensusParamsBz, consensusParams)
		if err != nil {
			panic(err)
		}

		app.setConsensusParams(consensusParams)
	}

	// needed for the export command which inits from store but never calls initchain
	app.setCheckState(abci.Header{})
	app.Seal()

	return nil
}

func (app *BaseApp) SetProtocolEngine(pe *protocol.ProtocolEngine) {
	if app.sealed {
		panic("SetProtocolEngine() on sealed BaseApp")
	}
	app.Engine = pe
}

func (app *BaseApp) setMinGasPrices(gasPrices sdk.DecCoins) {
	fmt.Fprintf(os.Stdout, "minimus gas prices: %s\n", gasPrices.String())
	app.minGasPrices = gasPrices
}

func (app *BaseApp) setHaltHeight(height uint64) {
	app.haltHeight = height
}

// Seal seals a BaseApp. It prohibits any further modifications to a BaseApp.
func (app *BaseApp) Seal() { app.sealed = true }

// IsSealed returns true if the BaseApp is sealed and false otherwise.
func (app *BaseApp) IsSealed() bool { return app.sealed }

// setCheckState sets checkState with the cached multistore and
// the context wrapping it.
// It is called by InitChain() and Commit()
func (app *BaseApp) setCheckState(header abci.Header) {
	ms := app.cms.CacheMultiStore()
	app.checkState = &state{
		ms:  ms,
		ctx: sdk.NewContext(ms, header, true, app.logger).WithMinGasPrices(app.minGasPrices),
	}
}

// setCheckState sets checkState with the cached multistore and
// the context wrapping it.
// It is called by InitChain() and BeginBlock(),
// and deliverState is set nil on Commit().
func (app *BaseApp) setDeliverState(header abci.Header) {
	ms := app.cms.CacheMultiStore()
	app.deliverState = &state{
		ms:  ms,
		ctx: sdk.NewContext(ms, header, false, app.logger),
	}
}

// setConsensusParams memoizes the consensus params.
func (app *BaseApp) setConsensusParams(consensusParams *abci.ConsensusParams) {
	app.consensusParams = consensusParams
}

// setConsensusParams stores the consensus params to the main store.
func (app *BaseApp) storeConsensusParams(consensusParams *abci.ConsensusParams) {
	consensusParamsBz, err := proto.Marshal(consensusParams)
	if err != nil {
		panic(err)
	}
	mainStore := app.cms.GetKVStore(app.baseKey)
	mainStore.Set(mainConsensusParamsKey, consensusParamsBz)
}

// getMaximumBlockGas gets the maximum gas from the consensus params. It panics
// if maximum block gas is less than negative one and returns zero if negative
// one.
func (app *BaseApp) getMaximumBlockGas() uint64 {
	if app.consensusParams == nil || app.consensusParams.Block == nil {
		return 0
	}

	maxGas := app.consensusParams.Block.MaxGas
	switch {
	case maxGas < -1:
		panic(fmt.Sprintf("invalid maximum block gas: %d", maxGas))

	case maxGas == -1:
		return 0

	default:
		return uint64(maxGas)
	}
}

func (app *BaseApp) validateHeight(req abci.RequestBeginBlock) error {
	if req.Header.Height < 1 {
		return fmt.Errorf("invalid height: %d", req.Header.Height)
	}

	prevHeight := app.LastBlockHeight()
	if req.Header.Height != prevHeight+1 {
		return fmt.Errorf("invalid height: %d; expected: %d", req.Header.Height, prevHeight+1)
	}

	return nil
}

// validateBasicTxMsgs executes basic validator calls for messages.
func validateBasicTxMsgs(msgs []sdk.Msg) error {
	if msgs == nil || len(msgs) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "must contain at least one message")
	}

	for _, msg := range msgs {
		// Validate the Msg.
		err := msg.ValidateBasic()
		if err != nil {
			return err
		}
	}

	return nil
}

// Returns the applications's deliverState if app is in runTxModeDeliver,
// otherwise it returns the application's checkstate.
func (app *BaseApp) getState(mode runTxMode) *state {
	if mode == runTxModeDeliver {
		return app.deliverState
	}

	return app.checkState
}

func (app *BaseApp) GetCms() sdk.CommitMultiStore {
	return app.cms
}

// retrieve the context for the tx w/ txBytes and other memoized values.
func (app *BaseApp) getContextForTx(mode runTxMode, txBytes []byte) (ctx sdk.Context) {
	ctx = app.getState(mode).ctx.
		WithTxBytes(txBytes).
		WithVoteInfos(app.voteInfos).
		WithConsensusParams(app.consensusParams)

	if mode == runTxModeReCheck {
		ctx = ctx.WithIsReCheckTx(true)
	}
	if mode == runTxModeSimulate {
		ctx, _ = ctx.CacheContext()
	}

	return
}

// cacheTxContext returns a new context based off of the provided context with
// a cache wrapped multi-store.
func (app *BaseApp) cacheTxContext(ctx sdk.Context, txBytes []byte) (sdk.Context, sdk.CacheMultiStore) {
	ms := ctx.MultiStore()
	msCache := ms.CacheMultiStore()
	if msCache.TracingEnabled() {
		msCache = msCache.SetTracingContext(
			sdk.TraceContext(
				map[string]interface{}{
					"txHash": fmt.Sprintf("%X", tmhash.Sum(txBytes)),
				},
			),
		).(sdk.CacheMultiStore)
	}

	return ctx.WithMultiStore(msCache), msCache
}

// ----------------------------------------------------------------------------
// ABCI

// runMsgs iterates through all the messages and executes them.
// nolint: gocyclo
func (app *BaseApp) runMsgs(ctx sdk.Context, msgs []sdk.Msg, mode runTxMode) (*sdk.Result, error) {
	idxLogs := make([]sdk.ABCIMessageLog, 0, len(msgs)) // a list of JSON-encoded logs with msg index

	var data []byte
	var gasSum uint64
	events := sdk.EmptyEvents()

	// NOTE: GasWanted is determined by the AnteHandler and GasUsed by the GasMeter.
	for i, msg := range msgs {
		// skip actual execution for (Re)CheckTx mode
		if mode == runTxModeCheck || mode == runTxModeReCheck {
			break
		}

		msgRoute := msg.Route()
		//handler := app.router.Route(ctx, msgRoute)
		handler := app.Engine.GetCurrentProtocol().GetRouter().Route(ctx, msgRoute) //TODO ctx is not needed
		if handler == nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized message route: %s; message index: %d", msgRoute, i)
		}

		idxLog := sdk.ABCIMessageLog{MsgIndex: uint16(i)}

		ctx.Simulate = false
		if mode == runTxModeSimulate {
			ctx.Simulate = true
		}

		msgResult, err := handler(ctx, msg)

		if err != nil {
			idxLog.Success = false
			idxLog.Log = fmt.Sprintf("failed to execute message; message index: %d. error: %s", i, err.Error())
			idxLogs = append(idxLogs, idxLog)

			logJSON := codec.Cdc.MustMarshalJSON(idxLogs)

			return &sdk.Result{
				Data:   data,
				Log:    strings.TrimSpace(string(logJSON)),
				Events: events,
			}, sdkerrors.Wrapf(err, "failed to execute message; message index: %d", i)
		}

		data = append(data, msgResult.Data...)
		gasSum += msgResult.GasUsed
		// append events from the message's execution and a message action event
		events = events.AppendEvent(sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type())))
		events = events.AppendEvents(msgResult.Events)

		idxLog.Success = true
		idxLogs = append(idxLogs, idxLog)
	}

	logJSON := codec.Cdc.MustMarshalJSON(idxLogs)
	return &sdk.Result{
		Data:    data, //TODO: how to get data of single msg
		Log:     strings.TrimSpace(string(logJSON)),
		Events:  events,
		GasUsed: gasSum,
	}, nil
}

// runTx processes a transaction. The transactions is processed via an
// anteHandler. The provided txBytes may be nil in some cases, eg. in tests. For
// further details on transaction execution, reference the BaseApp SDK
// documentation.
func (app *BaseApp) runTx(mode runTxMode, txBytes []byte, tx sdk.Tx) (gInfo sdk.GasInfo, result *sdk.Result, err error) {
	// NOTE: GasWanted should be returned by the AnteHandler. GasUsed is
	// determined by the GasMeter. We need access to the context to get the gas
	// meter so we initialize upfront.
	var gasWanted uint64

	ctx := app.getContextForTx(mode, txBytes)
	ms := ctx.MultiStore()

	// only run the tx if there is block gas remaining
	if mode == runTxModeDeliver && ctx.BlockGasMeter().IsOutOfGas() {
		gInfo = sdk.GasInfo{GasUsed: ctx.BlockGasMeter().GasConsumed()}
		return gInfo, nil, sdkerrors.Wrap(sdkerrors.ErrOutOfGas, "no block gas left to run tx")
	}

	var startingGas uint64
	if mode == runTxModeDeliver {
		startingGas = ctx.BlockGasMeter().GasConsumed()
	}

	// Add cache in fee refund. If an error is returned or panic happens during refund,
	// no value will be written into blockchain state
	defer func() {
		feeRefundHandler := app.Engine.GetCurrentProtocol().GetFeeRefundHandler()
		if mode == runTxModeDeliver && result != nil && feeRefundHandler != nil {
			result.GasUsed = ctx.GasMeter().GasConsumed()
			result.GasWanted = gasWanted

			var refundCtx sdk.Context
			var refundCache sdk.CacheMultiStore

			refundCtx, refundCache = app.cacheTxContext(ctx, txBytes)

			// refund unspent fee
			_, err := feeRefundHandler(refundCtx, tx, *result)
			if err != nil {
				panic(sdkerrors.Wrap(sdkerrors.ErrPanic, err.Error()))
			}
			refundCache.Write()
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			switch rType := r.(type) {
			// TODO: Use ErrOutOfGas instead of ErrorOutOfGas which would allow us
			// to keep the stracktrace.
			case sdk.ErrorOutOfGas:
				err = sdkerrors.Wrap(
					sdkerrors.ErrOutOfGas, fmt.Sprintf(
						"out of gas in location: %v; gasWanted: %d, gasUsed: %d",
						rType.Descriptor, gasWanted, ctx.GasMeter().GasConsumed(),
					),
				)

			default:
				err = sdkerrors.Wrap(
					sdkerrors.ErrPanic, fmt.Sprintf(
						"recovered: %v\nstack:\n%v", r, string(debug.Stack()),
					),
				)
			}

			result = nil
		}

		gInfo = sdk.GasInfo{GasWanted: gasWanted, GasUsed: ctx.GasMeter().GasConsumed()}
	}()

	// If BlockGasMeter() panics it will be caught by the above recover and will
	// return an error - in any case BlockGasMeter will consume gas past the limit.
	//
	// NOTE: This must exist in a separate defer function for the above recovery
	// to recover from this one.
	defer func() {
		if mode == runTxModeDeliver {
			ctx.BlockGasMeter().ConsumeGas(
				ctx.GasMeter().GasConsumedToLimit(), "block gas meter",
			)

			if ctx.BlockGasMeter().GasConsumed() < startingGas {
				panic(sdk.ErrorGasOverflow{Descriptor: "tx gas summation"})
			}
		}
	}()

	msgs := tx.GetMsgs()
	if err := validateBasicTxMsgs(msgs); err != nil {
		gInfo = sdk.GasInfo{GasUsed: ctx.BlockGasMeter().GasConsumed()}
		return gInfo, nil, err
	}

	anteHandler := app.Engine.GetCurrentProtocol().GetAnteHandler()
	if anteHandler != nil {
		var anteCtx sdk.Context
		var msCache sdk.CacheMultiStore

		// Cache wrap context before AnteHandler call in case it aborts.
		// This is required for both CheckTx and DeliverTx.
		// Ref: https://github.com/cosmos/cosmos-sdk/issues/2772
		//
		// NOTE: Alternatively, we could require that AnteHandler ensures that
		// writes do not happen if aborted/failed.  This may have some
		// performance benefits, but it'll be more difficult to get right.
		anteCtx, msCache = app.cacheTxContext(ctx, txBytes)

		newCtx, err := anteHandler(anteCtx, tx, mode == runTxModeSimulate)
		if !newCtx.IsZero() {
			// At this point, newCtx.MultiStore() is cache-wrapped, or something else
			// replaced by the AnteHandler. We want the original multistore, not one
			// which was cache-wrapped for the AnteHandler.
			//
			// Also, in the case of the tx aborting, we need to track gas consumed via
			// the instantiated gas meter in the AnteHandler, so we update the context
			// prior to returning.
			ctx = newCtx.WithMultiStore(ms)
		}

		// GasMeter expected to be set in AnteHandler
		gasWanted = ctx.GasMeter().Limit()

		if err != nil {
			return gInfo, nil, err
		}

		msCache.Write()
	}

	// Create a new Context based off of the existing Context with a cache-wrapped
	// MultiStore in case message processing fails. At this point, the MultiStore
	// is doubly cached-wrapped.
	runMsgCtx, msCache := app.cacheTxContext(ctx, txBytes)

	// Attempt to execute all messages and only update state if all messages pass
	// and we're in DeliverTx. Note, runMsgs will never return a reference to a
	// Result if any single message fails or does not have a registered Handler.
	r, err := app.runMsgs(runMsgCtx, msgs, mode)
	if err == nil && mode == runTxModeDeliver {
		msCache.Write()
	}

	return gInfo, r, err
}

// ----------------------------------------------------------------------------
// State

type state struct {
	ms  sdk.CacheMultiStore
	ctx sdk.Context
}

func (st *state) CacheMultiStore() sdk.CacheMultiStore {
	return st.ms.CacheMultiStore()
}

func (st *state) Context() sdk.Context {
	return st.ctx
}
