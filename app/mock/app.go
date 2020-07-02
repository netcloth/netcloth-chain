package mock

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/netcloth/netcloth-chain/app/protocol"
	v0 "github.com/netcloth/netcloth-chain/app/v0"
	"github.com/netcloth/netcloth-chain/app/v0/auth"
	authexported "github.com/netcloth/netcloth-chain/app/v0/auth/exported"
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
	*BaseApp
	GenesisAccounts []authexported.Account
}

func NewNCHApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, invCheckPeriod uint, baseAppOptions ...func(*BaseApp)) *NCHApp {
	baseApp := NewBaseApp(appName, logger, db, baseAppOptions...)

	baseApp.SetCommitMultiStoreTracer(traceStore)
	baseApp.SetAppVersion(version.Version)

	protocolKeeper := sdk.NewProtocolKeeper(protocol.Keys[protocol.MainStoreKey])
	engine := protocol.NewProtocolEngine(protocolKeeper)
	baseApp.SetProtocolEngine(&engine)

	if !baseApp.fauxMerkleMode {
		baseApp.MountStore(protocol.Keys[protocol.MainStoreKey], sdk.StoreTypeIAVL)
	} else {
		baseApp.MountStore(protocol.Keys[protocol.MainStoreKey], sdk.StoreTypeDB)
	}

	baseApp.MountKVStores(protocol.Keys)
	baseApp.MountTransientStores(protocol.TKeys)

	if loadLatest {
		err := baseApp.LoadLatestVersion(protocol.Keys[protocol.MainStoreKey])
		if err != nil {
			cmn.Exit(err.Error())
		}
	}

	engine.Add(v0.NewProtocolV0(0, logger, protocolKeeper, baseApp.DeliverTx, invCheckPeriod, nil))

	loaded, current := engine.LoadCurrentProtocol(baseApp.Cms.GetKVStore(protocol.Keys[protocol.MainStoreKey]))
	if !loaded {
		cmn.Exit(fmt.Sprintf("Your software doesn't support the required protocol (version %d)!, to upgrade nchd", current))
	} else {
		fmt.Printf("blockchain current protocol version :%d\n", current)
	}

	baseApp.TxDecoder = auth.DefaultTxDecoder(engine.GetCurrentProtocol().GetCodec())

	var app = &NCHApp{baseApp, nil}

	return app
}

func MakeLatestCodec() *codec.Codec {
	return v0.MakeCodec()
}

func (app *NCHApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, protocol.Keys[protocol.MainStoreKey])
}

func (app *NCHApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	return app.Engine.GetCurrentProtocol().ExportAppStateAndValidators(ctx, forZeroHeight, jailWhiteList)
}

// for test
