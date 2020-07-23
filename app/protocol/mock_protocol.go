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

// GetInitChainer
func (m *MockProtocol) GetInitChainer() sdk.InitChainer {
	return m.initChainer
}

// GetBeginBlocker
func (m MockProtocol) GetBeginBlocker() sdk.BeginBlocker {
	return m.beginBlocker
}

// GetEndBlocker
func (m MockProtocol) GetEndBlocker() sdk.EndBlocker {
	return m.endBlocker
}

// ExportAppStateAndValidators
func (m *MockProtocol) ExportAppStateAndValidators(ctx sdk.Context, forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, validators []types.GenesisValidator, err error) {
	return json.RawMessage{}, nil, nil
}

// LoadContext
func (m *MockProtocol) LoadContext() {
}

// Init
func (m *MockProtocol) Init() {
}

// GetCodec
func (m *MockProtocol) GetCodec() *codec.Codec {
	return codec.New()
}

func (m *MockProtocol) SetRouter(router sdk.Router) {
	m.router = router
}

func (m *MockProtocol) SetQuearyRouter(queryRouter sdk.QueryRouter) {
	m.queryRouter = queryRouter
}

func (m *MockProtocol) SetAnteHandler(anteHandler sdk.AnteHandler) {
	m.anteHandler = anteHandler
}

func (m *MockProtocol) SetInitChainer(initChainer sdk.InitChainer) {
	m.initChainer = initChainer
}

// for simulation
func (m *MockProtocol) GetSimulationManager() interface{} {
	return nil
}
