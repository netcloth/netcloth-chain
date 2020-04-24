package mock

import (
	"encoding/json"
	"github.com/netcloth/netcloth-chain/app/protocol"
	"github.com/netcloth/netcloth-chain/app/v0/genutil"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/types/module"
	"github.com/tendermint/tendermint/types"
)

type ProtocolV0 struct {
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

func newMockProtocolV0() *ProtocolV0 {
	return &ProtocolV0{
		router:        protocol.NewRouter(),
		queryRouter:   protocol.NewQueryRouter(),
		moduleManager: module.NewManager(),
	}
}

var _ protocol.Protocol = &ProtocolV0{}

func (m *ProtocolV0) GetVersion() uint64 {
	return 0
}

func (m *ProtocolV0) GetRouter() sdk.Router {
	return m.router
}

func (m *ProtocolV0) GetQueryRouter() sdk.QueryRouter {
	return m.queryRouter
}

func (m ProtocolV0) GetAnteHandler() sdk.AnteHandler {
	return m.anteHandler
}

func (m *ProtocolV0) GetFeeRefundHandler() sdk.FeeRefundHandler {
	return m.feeRefundHandler
}

func (m *ProtocolV0) GetInitChainer() sdk.InitChainer {
	return m.initChainer
}

func (m ProtocolV0) GetBeginBlocker() sdk.BeginBlocker {
	return m.beginBlocker
}

func (m ProtocolV0) GetEndBlocker() sdk.EndBlocker {
	return m.endBlocker
}

func (m *ProtocolV0) ExportAppStateAndValidators(ctx sdk.Context, forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, validators []types.GenesisValidator, err error) {
	return json.RawMessage{}, nil, nil
}

func (m *ProtocolV0) Load() {
}

func (m *ProtocolV0) Init(ctx sdk.Context) {
}

func (m *ProtocolV0) GetCodec() *codec.Codec {
	return codec.New()
}

func (m *ProtocolV0) SetRouter(router sdk.Router) {
	m.router = router
}

func (m *ProtocolV0) SetQuearyRouter(queryRouter sdk.QueryRouter) {
	m.queryRouter = queryRouter
}

func (m *ProtocolV0) SetAnteHandler(anteHandler sdk.AnteHandler) {
	m.anteHandler = anteHandler
}

func (m *ProtocolV0) SetInitChainer(initChainer sdk.InitChainer) {
	m.initChainer = initChainer
}
