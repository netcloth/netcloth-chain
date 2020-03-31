package app

import (
	"fmt"
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/netcloth/netcloth-chain/app/protocol"
	v0 "github.com/netcloth/netcloth-chain/app/v0"
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/auth"
	"github.com/netcloth/netcloth-chain/modules/auth/ante"
	"github.com/netcloth/netcloth-chain/modules/bank"
	"github.com/netcloth/netcloth-chain/modules/cipal"
	"github.com/netcloth/netcloth-chain/modules/crisis"
	distr "github.com/netcloth/netcloth-chain/modules/distribution"
	distrclient "github.com/netcloth/netcloth-chain/modules/distribution/client"
	"github.com/netcloth/netcloth-chain/modules/genaccounts"
	"github.com/netcloth/netcloth-chain/modules/genutil"
	"github.com/netcloth/netcloth-chain/modules/gov"
	"github.com/netcloth/netcloth-chain/modules/ipal"
	"github.com/netcloth/netcloth-chain/modules/mint"
	"github.com/netcloth/netcloth-chain/modules/params"
	paramsclient "github.com/netcloth/netcloth-chain/modules/params/client"
	"github.com/netcloth/netcloth-chain/modules/slashing"
	"github.com/netcloth/netcloth-chain/modules/staking"
	"github.com/netcloth/netcloth-chain/modules/supply"
	"github.com/netcloth/netcloth-chain/modules/vm"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/types/module"
	"github.com/netcloth/netcloth-chain/version"
)

const (
	appName = "nch"
)

var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.nchcli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.nchd")

	// The module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		genaccounts.AppModuleBasic{},
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsclient.ProposalHandler, distrclient.ProposalHandler),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
		cipal.AppModuleBasic{},
		ipal.AppModuleBasic{},
		vm.AppModuleBasic{},
	)

	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		ipal.ModuleName:           {supply.Staking},
	}
)

func CreateCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)

	return cdc
}

type NCHApp struct {
	*BaseApp

	invCheckPeriod uint

	accountKeeper  auth.AccountKeeper
	refundKeeper   auth.RefundKeeper
	bankKeeper     bank.Keeper
	supplyKeeper   supply.Keeper
	stakingKeeper  staking.Keeper
	slashingKeeper slashing.Keeper
	mintKeeper     mint.Keeper
	distrKeeper    distr.Keeper
	govKeeper      gov.Keeper
	crisisKeeper   crisis.Keeper
	paramsKeeper   params.Keeper
	cipalKeeper    cipal.Keeper
	ipalKeeper     ipal.Keeper
	vmKeeper       vm.Keeper

	mm *module.Manager
}

func NewNCHApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, invCheckPeriod uint, baseAppOptions ...func(*BaseApp)) *NCHApp {
	bApp := NewBaseApp(appName, logger, db, baseAppOptions...)

	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	//bApp.SetAnteHandler(ante.NewAnteHandler(bApp.accountKeeper, bApp.supplyKeeper, ante.DefaultSigVerificationGasConsumer))
	//bApp.SetFeeRefundHandler(auth.NewFeeRefundHandler(app.accountKeeper, app.supplyKeeper, app.refundKeeper))

	protocolKeeper := sdk.NewProtocolKeeper(protocol.MainKVStoreKey)
	engine := protocol.NewProtocolEngine(protocolKeeper)
	bApp.SetProtocolEngine(&engine)

	if !bApp.fauxMerkleMode {
		bApp.MountStore(protocol.MainKVStoreKey, sdk.StoreTypeIAVL)
	} else {
		bApp.MountStore(protocol.MainKVStoreKey, sdk.StoreTypeDB)
	}

	bApp.MountKVStores(protocol.Keys)
	bApp.MountTransientStores(protocol.TKeys)

	if loadLatest {
		err := bApp.LoadLatestVersion(protocol.MainKVStoreKey)
		if err != nil {
			cmn.Exit(err.Error())
		}
	}

	engine.Add(v0.NewProtocolV0(0, logger, protocolKeeper, bApp.DeliverTx, nil))
	loaded, current := engine.LoadCurrentProtocol(bApp.cms.GetKVStore(protocol.MainKVStoreKey))
	if !loaded {
		cmn.Exit(fmt.Sprintf("Your software doesn't support the required protocol (version %d)!", current))
	}
	bApp.txDecoder = auth.DefaultTxDecoder(engine.GetCurrentProtocol().GetCodec())

	var app = &NCHApp{
		BaseApp:        bApp,
		invCheckPeriod: invCheckPeriod,
	}

	return app
}

func NewNCHAppForReplay(logger log.Logger, db dbm.DB, traceStore io.Writer, loadInit, loadLatest bool, invCheckPeriod uint, baseAppOptions ...func(*BaseApp)) *NCHApp {
	bApp := NewBaseApp(appName, logger, db, baseAppOptions...)

	protocolKeeper := sdk.NewProtocolKeeper(protocol.MainKVStoreKey)
	engine := protocol.NewProtocolEngine(protocolKeeper)
	bApp.SetProtocolEngine(&engine)

	if !bApp.fauxMerkleMode {
		bApp.MountStore(protocol.MainKVStoreKey, sdk.StoreTypeIAVL)
	} else {
		bApp.MountStore(protocol.MainKVStoreKey, sdk.StoreTypeDB)
	}

	engine.Add(v0.NewProtocolV0(0, logger, protocolKeeper, bApp.DeliverTx, nil))
	loaded, current := engine.LoadCurrentProtocol(bApp.cms.GetKVStore(protocol.Keys[MainStoreKey]))
	if !loaded {
		cmn.Exit(fmt.Sprintf("Your software doesn't support the required protocol (version %d)!", current))
	}

	bApp.txDecoder = auth.DefaultTxDecoder(engine.GetCurrentProtocol().GetCodec())

	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	var app = &NCHApp{
		BaseApp:        bApp,
		invCheckPeriod: invCheckPeriod,
	}

	app.MountKVStores(protocol.Keys)
	app.MountTransientStores(protocol.TKeys)

	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(ante.NewAnteHandler(app.accountKeeper, app.supplyKeeper, ante.DefaultSigVerificationGasConsumer))
	app.SetFeeRefundHandler(auth.NewFeeRefundHandler(app.accountKeeper, app.supplyKeeper, app.refundKeeper))
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion(protocol.MainKVStoreKey)
		if err != nil {
			cmn.Exit(err.Error())
		}
	} else if loadInit {
		err := app.LoadVersion(0, protocol.MainKVStoreKey)
		if err != nil {
			cmn.Exit(err.Error())
		}
	}

	return app
}

func (app *NCHApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *NCHApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func SetBech32AddressPrefixes(config *sdk.Config) {
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
}

func (app *NCHApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, protocol.MainKVStoreKey)
}

func (app *NCHApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}
