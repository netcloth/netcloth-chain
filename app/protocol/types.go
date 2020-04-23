package protocol

import (
	"encoding/json"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
)

type Protocol interface {
	GetVersion() uint64
	GetRouter() sdk.Router
	GetQueryRouter() sdk.QueryRouter
	GetAnteHandler() sdk.AnteHandler
	GetFeeRefundHandler() sdk.FeeRefundHandler
	GetInitChainer() sdk.InitChainer
	GetBeginBlocker() sdk.BeginBlocker
	GetEndBlocker() sdk.EndBlocker

	ExportAppStateAndValidators(ctx sdk.Context, forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error)

	Load()
	Init(ctx sdk.Context)
	GetCodec() *codec.Codec

	//for test
	SetAnteHandler(anteHandler sdk.AnteHandler)
}
