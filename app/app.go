package app

import (
	"encoding/json"
	"fmt"
	"github.com/netcloth/netcloth-chain/codec"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"io"
	"os"

	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/netcloth/netcloth-chain/app/protocol"
	v0 "github.com/netcloth/netcloth-chain/app/v0"
	"github.com/netcloth/netcloth-chain/app/v0/auth"
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
	*BaseApp
	invCheckPeriod uint
}

func NewNCHApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, invCheckPeriod uint, baseAppOptions ...func(*BaseApp)) *NCHApp {
	baseApp := NewBaseApp(appName, logger, db, baseAppOptions...)

	baseApp.SetCommitMultiStoreTracer(traceStore)
	baseApp.SetAppVersion(version.Version)

	protocolKeeper := sdk.NewProtocolKeeper(protocol.MainKVStoreKey)
	engine := protocol.NewProtocolEngine(protocolKeeper)
	baseApp.SetProtocolEngine(&engine)

	if !baseApp.fauxMerkleMode {
		baseApp.MountStore(protocol.MainKVStoreKey, sdk.StoreTypeIAVL)
	} else {
		baseApp.MountStore(protocol.MainKVStoreKey, sdk.StoreTypeDB)
	}

	baseApp.MountKVStores(protocol.Keys)
	baseApp.MountTransientStores(protocol.TKeys)

	if loadLatest {
		err := baseApp.LoadLatestVersion(protocol.MainKVStoreKey)
		if err != nil {
			cmn.Exit(err.Error())
		}
	}

	engine.Add(v0.NewProtocolV0(0, logger, protocolKeeper, baseApp.DeliverTx, invCheckPeriod, nil))

	loaded, current := engine.LoadCurrentProtocol(baseApp.cms.GetKVStore(protocol.MainKVStoreKey))
	if !loaded {
		cmn.Exit(fmt.Sprintf("Your software doesn't support the required protocol (version %d)!, to upgrade nchd", current))
	} else {
		fmt.Println(fmt.Sprintf("blockchain current protocol version :%d", current))
	}

	baseApp.txDecoder = auth.DefaultTxDecoder(engine.GetCurrentProtocol().GetCodec())

	var app = &NCHApp{
		BaseApp:        baseApp,
		invCheckPeriod: invCheckPeriod,
	}

	return app
}

func MakeLatestCodec() *codec.Codec {
	return v0.MakeCodec()
}

func (app *NCHApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, protocol.MainKVStoreKey)
}

func (app *NCHApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	return app.Engine.GetCurrentProtocol().ExportAppStateAndValidators(ctx, forZeroHeight, jailWhiteList)
}
