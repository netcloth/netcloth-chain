package protocol

import (
	"encoding/json"

	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/types/module"

	"github.com/tendermint/tendermint/types"
)

var _ Protocol = &MockProtocol{}

// MockProtocol is designed for engine test
type MockProtocol struct {
	version uint64

	moduleManager *module.Manager

	router      sdk.Router
	queryRouter sdk.QueryRouter

	anteHandler      sdk.AnteHandler
	feeRefundHandler sdk.FeeRefundHandler
	initChainer      sdk.InitChainer
	beginBlocker     sdk.BeginBlocker
	endBlocker       sdk.EndBlocker
}

// NewMockProtocol creates a new instance of MockProtocol
func NewMockProtocol(version uint64) *MockProtocol {
	return &MockProtocol{
		version:       version,
		router:        NewRouter(),
		queryRouter:   NewQueryRouter(),
		moduleManager: module.NewManager(),
	}
}

// GetVersion gets version
func (m *MockProtocol) GetVersion() uint64 {
	return m.version
}

// GetRouter gets router
func (m *MockProtocol) GetRouter() sdk.Router {
	return m.router
}

// GetQueryRouter gets query router
func (m *MockProtocol) GetQueryRouter() sdk.QueryRouter {
	return m.queryRouter
}

// GetAnteHandler gets ante handler
func (m MockProtocol) GetAnteHandler() sdk.AnteHandler {
	return m.anteHandler
}

// GetFeeRefundHandler gets fee refund handler
func (m *MockProtocol) GetFeeRefundHandler() sdk.FeeRefundHandler {
	return m.feeRefundHandler
}

// GetInitChainer get initChainer
func (m *MockProtocol) GetInitChainer() sdk.InitChainer {
	return m.initChainer
}

// GetBeginBlocker gets beginBlocker
func (m MockProtocol) GetBeginBlocker() sdk.BeginBlocker {
	return m.beginBlocker
}

// GetEndBlocker gets endBlocker
func (m MockProtocol) GetEndBlocker() sdk.EndBlocker {
	return m.endBlocker
}

// ExportAppStateAndValidators exports the application state for a genesis file
func (m *MockProtocol) ExportAppStateAndValidators(ctx sdk.Context, forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, validators []types.GenesisValidator, err error) {
	return json.RawMessage{}, nil, nil
}

// LoadContext - do nothing
func (m *MockProtocol) LoadContext() {
	// do nothing
}

// Init - initialize
func (m *MockProtocol) Init() {
	// do nothing
}

// GetCodec gets codec
func (m *MockProtocol) GetCodec() *codec.Codec {
	return codec.New()
}

// SetRouter allows us to customize the router
func (m *MockProtocol) SetRouter(router sdk.Router) {
	m.router = router
}

// SetQueryRouter allows us to customize the query router
func (m *MockProtocol) SetQueryRouter(queryRouter sdk.QueryRouter) {
	m.queryRouter = queryRouter
}

// SetAnteHandler set the anteHandler
func (m *MockProtocol) SetAnteHandler(anteHandler sdk.AnteHandler) {
	m.anteHandler = anteHandler
}

// SetInitChainer set the initChainer
func (m *MockProtocol) SetInitChainer(initChainer sdk.InitChainer) {
	m.initChainer = initChainer
}

// GetSimulationManager - for simulation
func (m *MockProtocol) GetSimulationManager() interface{} {
	return nil
}
