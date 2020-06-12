package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/netcloth/netcloth-chain/baseapp"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/netcloth/netcloth-chain/app/protocol"
	v0 "github.com/netcloth/netcloth-chain/app/v0"
	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/version"
)

const (
	appName = "nch"
)

var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.nchcli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.nchd")
)

type NCHApp struct {
	*baseapp.BaseApp
}

func NewNCHApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, invCheckPeriod uint, baseAppOptions ...func(*baseapp.BaseApp)) *NCHApp {
	baseApp := baseapp.NewBaseApp(appName, logger, db, baseAppOptions...)

	baseApp.SetCommitMultiStoreTracer(traceStore)
	baseApp.SetAppVersion(version.Version)

	mainStoreKey := protocol.Keys[protocol.MainStoreKey]
	protocolKeeper := sdk.NewProtocolKeeper(mainStoreKey)
	engine := protocol.NewProtocolEngine(protocolKeeper)
	baseApp.SetProtocolEngine(&engine)
	baseApp.MountKVStores(protocol.Keys)
	baseApp.MountTransientStores(protocol.TKeys)

	var app = &NCHApp{baseApp}

	// set hook function postEndBlocker
	baseApp.PostEndBlocker = app.postEndBlocker

	if loadLatest {
		err := app.LoadLatestVersion(mainStoreKey)
		if err != nil {
			cmn.Exit(err.Error())
		}
	}

	engine.Add(v0.NewProtocolV0(0, logger, protocolKeeper, app.DeliverTx, invCheckPeriod, nil))

	loaded, current := engine.LoadCurrentProtocol(app.GetCms().GetKVStore(mainStoreKey))
	if !loaded {
		cmn.Exit(fmt.Sprintf("Your software doesn't support the required protocol (version %d)!, to upgrade nchd", current))
	}
	logger.Info(fmt.Sprintf("launch app with protocol version: %d", current))

	// set txDeocder
	app.SetTxDecoder(auth.DefaultTxDecoder(engine.GetCurrentProtocol().GetCodec()))

	return app
}

func MakeLatestCodec() *codec.Codec {
	return v0.MakeCodec()
}

func (app *NCHApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, protocol.Keys[protocol.MainStoreKey])
}

// hook function for BaseApp's EndBlock(upgrade)
func (app *NCHApp) postEndBlocker(res *abci.ResponseEndBlock) {
	appVersion := app.Engine.GetCurrentVersion()
	for _, event := range res.Events {
		if event.Type == sdk.AppVersionEvent {
			for _, attr := range event.Attributes {
				if string(attr.Key) == sdk.AppVersionEvent {
					appVersion, _ = strconv.ParseUint(string(attr.Value), 10, 64)
					break
				}
			}

			break
		}
	}

	if appVersion <= app.Engine.GetCurrentVersion() {
		return
	}

	success := app.Engine.Activate(appVersion)
	if success {
		app.SetTxDecoder(auth.DefaultTxDecoder(app.Engine.GetCurrentProtocol().GetCodec()))
		return
	}

	app.Log(fmt.Sprintf("activate version from %d to %d failed, please upgrade your app", app.Engine.GetCurrentVersion(), appVersion))
	return
}

func (app *NCHApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	return app.Engine.GetCurrentProtocol().ExportAppStateAndValidators(ctx, forZeroHeight, jailWhiteList)
}
