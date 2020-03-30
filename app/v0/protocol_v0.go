package v0

import (
	"github.com/netcloth/netcloth-chain/app/protocol"
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/auth"
	"github.com/netcloth/netcloth-chain/modules/bank"
	"github.com/netcloth/netcloth-chain/modules/cipal"
	"github.com/netcloth/netcloth-chain/modules/crisis"
	distr "github.com/netcloth/netcloth-chain/modules/distribution"
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
	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
)

var _ protocol.Protocol = (*ProtocolV0)(nil)

// The module BasicManager is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.
var ModuleBasics = module.NewBasicManager(
	genaccounts.AppModuleBasic{},
	genutil.AppModuleBasic{},
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distr.AppModuleBasic{},
	gov.NewAppModuleBasic(paramsclient.ProposalHandler, distr.ProposalHandler),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	supply.AppModuleBasic{},
	cipal.AppModuleBasic{},
	ipal.AppModuleBasic{},
	vm.AppModuleBasic{},
)

var maccPerms = map[string][]string{
	auth.FeeCollectorName:     nil,
	distr.ModuleName:          nil,
	mint.ModuleName:           {supply.Minter},
	staking.BondedPoolName:    {supply.Burner, supply.Staking},
	staking.NotBondedPoolName: {supply.Burner, supply.Staking},
	gov.ModuleName:            {supply.Burner},
	ipal.ModuleName:           {supply.Staking},
}

var invCheckPeriod uint

type ProtocolV0 struct {
	version        uint64
	cdc            *codec.Codec
	logger         log.Logger
	invariantLevel string
	checkInvariant bool
	trackCoinFlow  bool

	// Manage getting and setting accounts
	accountKeeper  auth.AccountKeeper
	refundKeeper   auth.RefundKeeper
	bankKeeper     bank.Keeper
	StakingKeeper  staking.Keeper
	slashingKeeper slashing.Keeper
	mintKeeper     mint.Keeper
	distrKeeper    distr.Keeper
	protocolKeeper sdk.ProtocolKeeper
	govKeeper      gov.Keeper
	crisisKeeper   crisis.Keeper
	paramsKeeper   params.Keeper
	supplyKeeper   supply.Keeper
	stakingKeeper  staking.Keeper
	ipalKeeper     ipal.Keeper
	cipalKeeper    cipal.Keeper
	vmKeeper       vm.Keeper

	//router      protocol.Router      // handle any kind of message
	//queryRouter protocol.QueryRouter // router for redirecting query calls
	router      sdk.Router      // handle any kind of message
	queryRouter sdk.QueryRouter // router for redirecting query calls

	anteHandlers     []sdk.AnteHandler    // ante handlers for fee and auth
	feeRefundHandler sdk.FeeRefundHandler // fee handler for fee refund
	//feePreprocessHandler sdk.FeePreprocessHandler // fee handler for fee preprocessor

	// may be nil
	initChainer  sdk.InitChainer  // initialize state with validators and state blob
	beginBlocker sdk.BeginBlocker // logic to run before any txs
	endBlocker   sdk.EndBlocker   // logic to run after all txs, and to determine valset changes
	config       *cfg.InstrumentationConfig

	mm *module.Manager
}

func NewProtocolV0(version uint64, log log.Logger, pk sdk.ProtocolKeeper, checkInvariant bool, trackCoinFlow bool, config *cfg.InstrumentationConfig) *ProtocolV0 {
	p0 := ProtocolV0{
		version:        version,
		logger:         log,
		protocolKeeper: pk,
		checkInvariant: checkInvariant,
		trackCoinFlow:  trackCoinFlow,
		config:         config,
	}

	return &p0
}

func (p *ProtocolV0) GetVersion() uint64 {
	return p.version
}

func (p *ProtocolV0) GetRouter() sdk.Router {
	panic("implement me")
}

func (p *ProtocolV0) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState protocol.GenesisState
	p.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	return p.mm.InitGenesis(ctx, genesisState)
}

func (p *ProtocolV0) GetBeginBlocker() sdk.BeginBlocker {
	return p.beginBlocker
}

func (p *ProtocolV0) GetEndBlocker() sdk.EndBlocker {
	return p.endBlocker
}

func (p *ProtocolV0) Load() {
	p.configCodec()
	p.configKeepers()
	//p.configRouters()
	//p.configFeeHandlers()
	//p.configParams()
}

func (p *ProtocolV0) Init(ctx sdk.Context) {
}

func (p *ProtocolV0) GetInitChainer() sdk.InitChainer {
	return p.InitChainer
}

func (p *ProtocolV0) configCodec() {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)

	p.cdc = cdc
}

func ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (p *ProtocolV0) configKeepers() {
	p.paramsKeeper = params.NewKeeper(p.cdc, protocol.Keys[params.StoreKey], protocol.TKeys[params.TStoreKey])
	authSubspace := p.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := p.paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := p.paramsKeeper.Subspace(staking.DefaultParamspace)
	mintSubspace := p.paramsKeeper.Subspace(mint.DefaultParamspace)
	distrSubspace := p.paramsKeeper.Subspace(distr.DefaultParamspace)
	slashingSubspace := p.paramsKeeper.Subspace(slashing.DefaultParamspace)
	govSubspace := p.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	crisisSubspace := p.paramsKeeper.Subspace(crisis.DefaultParamspace)
	cipalSubspace := p.paramsKeeper.Subspace(cipal.DefaultParamspace)
	ipalSubspace := p.paramsKeeper.Subspace(ipal.DefaultParamspace)
	vmSubspace := p.paramsKeeper.Subspace(vm.DefaultParamspace)

	p.accountKeeper = auth.NewAccountKeeper(p.cdc, protocol.Keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)
	p.refundKeeper = auth.NewRefundKeeper(p.cdc, protocol.Keys[auth.RefundKey])
	p.bankKeeper = bank.NewBaseKeeper(p.accountKeeper, bankSubspace, ModuleAccountAddrs())
	p.supplyKeeper = supply.NewKeeper(p.cdc, protocol.Keys[supply.StoreKey], p.accountKeeper, p.bankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(
		p.cdc, protocol.Keys[staking.StoreKey], protocol.TKeys[staking.TStoreKey],
		p.supplyKeeper, stakingSubspace)
	p.mintKeeper = mint.NewKeeper(p.cdc, protocol.Keys[mint.StoreKey], mintSubspace, &stakingKeeper, p.supplyKeeper, auth.FeeCollectorName)
	p.distrKeeper = distr.NewKeeper(p.cdc, protocol.Keys[distr.StoreKey], distrSubspace, &stakingKeeper,
		p.supplyKeeper, auth.FeeCollectorName, ModuleAccountAddrs())
	p.slashingKeeper = slashing.NewKeeper(
		p.cdc, protocol.Keys[slashing.StoreKey], &stakingKeeper, slashingSubspace)
	p.crisisKeeper = crisis.NewKeeper(crisisSubspace, invCheckPeriod, p.supplyKeeper, auth.FeeCollectorName)

	p.cipalKeeper = cipal.NewKeeper(
		protocol.Keys[cipal.StoreKey],
		p.cdc,
		cipalSubspace)

	p.ipalKeeper = ipal.NewKeeper(
		protocol.Keys[ipal.StoreKey],
		p.cdc,
		p.supplyKeeper,
		ipalSubspace)

	p.vmKeeper = vm.NewKeeper(
		p.cdc,
		protocol.Keys[vm.StoreKey],
		protocol.Keys[vm.CodeKey],
		protocol.Keys[vm.StoreDebugKey],
		vmSubspace,
		auth.NewAccountKeeperCopy(p.accountKeeper, true))

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.
		AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(p.paramsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(p.distrKeeper))

	p.govKeeper = gov.NewKeeper(
		p.cdc, protocol.Keys[gov.StoreKey], govSubspace, p.supplyKeeper,
		&stakingKeeper, govRouter,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	p.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(p.distrKeeper.Hooks(), p.slashingKeeper.Hooks()),
	)
}

func (p *ProtocolV0) LoadMM() {
	mm := module.NewManager(
		genaccounts.NewAppModule(p.accountKeeper),
		//genutil.NewAppModule(p.accountKeeper, p.stakingKeeper, app.Basep.DeliverTx),
		auth.NewAppModule(p.accountKeeper),
		bank.NewAppModule(p.bankKeeper, p.accountKeeper),
		crisis.NewAppModule(&p.crisisKeeper),
		supply.NewAppModule(p.supplyKeeper, p.accountKeeper),
		distr.NewAppModule(p.distrKeeper, p.supplyKeeper),
		gov.NewAppModule(p.govKeeper, p.supplyKeeper),
		mint.NewAppModule(p.mintKeeper),
		slashing.NewAppModule(p.slashingKeeper, p.stakingKeeper),
		staking.NewAppModule(p.stakingKeeper, p.distrKeeper, p.accountKeeper, p.supplyKeeper),
		cipal.NewAppModule(p.cipalKeeper),
		ipal.NewAppModule(p.ipalKeeper),
		vm.NewAppModule(p.vmKeeper),
	)

	mm.SetOrderBeginBlockers(mint.ModuleName, distr.ModuleName, slashing.ModuleName)

	mm.SetOrderEndBlockers(crisis.ModuleName, gov.ModuleName, staking.ModuleName, ipal.ModuleName, vm.ModuleName)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	mm.SetOrderInitGenesis(
		genaccounts.ModuleName,
		distr.ModuleName,
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,
		gov.ModuleName,
		mint.ModuleName,
		supply.ModuleName,
		crisis.ModuleName,
		genutil.ModuleName,
		ipal.ModuleName,
		cipal.ModuleName,
		vm.ModuleName,
	)

	mm.RegisterRoutes(p.router, p.queryRouter)

	p.mm = mm
}
