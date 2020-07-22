package v0

import (
	"github.com/netcloth/netcloth-chain/app/protocol"
	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"github.com/netcloth/netcloth-chain/app/v0/auth/ante"
	"github.com/netcloth/netcloth-chain/app/v0/bank"
	"github.com/netcloth/netcloth-chain/app/v0/cipal"
	"github.com/netcloth/netcloth-chain/app/v0/crisis"
	distr "github.com/netcloth/netcloth-chain/app/v0/distribution"
	distrclient "github.com/netcloth/netcloth-chain/app/v0/distribution/client"
	"github.com/netcloth/netcloth-chain/app/v0/genaccounts"
	"github.com/netcloth/netcloth-chain/app/v0/genutil"
	"github.com/netcloth/netcloth-chain/app/v0/gov"
	"github.com/netcloth/netcloth-chain/app/v0/guardian"
	"github.com/netcloth/netcloth-chain/app/v0/ipal"
	"github.com/netcloth/netcloth-chain/app/v0/mint"
	"github.com/netcloth/netcloth-chain/app/v0/params"
	paramsclient "github.com/netcloth/netcloth-chain/app/v0/params/client"
	"github.com/netcloth/netcloth-chain/app/v0/slashing"
	"github.com/netcloth/netcloth-chain/app/v0/staking"
	"github.com/netcloth/netcloth-chain/app/v0/supply"
	"github.com/netcloth/netcloth-chain/app/v0/upgrade"
	"github.com/netcloth/netcloth-chain/app/v0/upgrade/types"
	"github.com/netcloth/netcloth-chain/app/v0/vm"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/types/module"
)

var _ protocol.Protocol = (*ProtocolV0)(nil)

// ModuleBasics - The module BasicManager is in charge of setting up basic,
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
	gov.NewAppModuleBasic(paramsclient.ProposalHandler, distrclient.ProposalHandler),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	supply.AppModuleBasic{},
	cipal.AppModuleBasic{},
	ipal.AppModuleBasic{},
	vm.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	guardian.AppModuleBasic{},
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

// ProtocolV0 is the struct of the original protocol
type ProtocolV0 struct {
	version uint64
	cdc     *codec.Codec
	logger  log.Logger

	moduleManager *module.Manager
	simManager    *module.SimulationManager

	accountKeeper  auth.AccountKeeper
	refundKeeper   auth.RefundKeeper
	bankKeeper     bank.Keeper
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
	upgradeKeeper  upgrade.Keeper
	guardianKeeper guardian.Keeper

	router      sdk.Router
	queryRouter sdk.QueryRouter

	anteHandler      sdk.AnteHandler
	feeRefundHandler sdk.FeeRefundHandler

	initChainer sdk.InitChainer
	deliverTx   genutil.DeliverTxfn

	config *cfg.InstrumentationConfig

	invCheckPeriod uint
}

// NewProtocolV0 creates a new instance of ProtocolV0
func NewProtocolV0(version uint64, log log.Logger, pk sdk.ProtocolKeeper, deliverTx genutil.DeliverTxfn, invCheckPeriod uint, config *cfg.InstrumentationConfig) *ProtocolV0 {
	p0 := ProtocolV0{
		version:        version,
		logger:         log,
		protocolKeeper: pk,
		router:         protocol.NewRouter(),
		queryRouter:    protocol.NewQueryRouter(),
		config:         config,
		deliverTx:      deliverTx,
		invCheckPeriod: invCheckPeriod,
	}

	return &p0
}

// GetVersion gets the version of this protocol
func (p *ProtocolV0) GetVersion() uint64 {
	return p.version
}

// GetRouter
func (p *ProtocolV0) GetRouter() sdk.Router {
	return p.router
}

// GetQueryRouter
func (p *ProtocolV0) GetQueryRouter() sdk.QueryRouter {
	return p.queryRouter
}

// GetAnteHandler
func (p *ProtocolV0) GetAnteHandler() sdk.AnteHandler {
	return p.anteHandler
}

// GetFeeRefundHandler
func (p *ProtocolV0) GetFeeRefundHandler() sdk.FeeRefundHandler {
	return p.feeRefundHandler
}

// LoadContext updates the context for the app after the upgrade of protocol
func (p *ProtocolV0) LoadContext() {
	p.configCodec()
	p.configKeepers()
	p.configModuleManager()
	p.configSimulationManager()
	p.configRouters()
	p.configFeeHandlers()
}

// Init
func (p *ProtocolV0) Init() {
}

// GetCodec gets tx codec
func (p *ProtocolV0) GetCodec() *codec.Codec {
	return p.cdc
}

// GetInitChainer
func (p *ProtocolV0) GetInitChainer() sdk.InitChainer {
	return p.InitChainer
}

// GetBeginBlocker
func (p *ProtocolV0) GetBeginBlocker() sdk.BeginBlocker {
	return p.BeginBlocker
}

// GetEndBlocker
func (p *ProtocolV0) GetEndBlocker() sdk.EndBlocker {
	return p.EndBlocker
}

func (p *ProtocolV0) configCodec() {
	p.cdc = MakeCodec()
}

// MakeCodec registers codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)

	return cdc
}

// ModuleAccountAddrs returns all the module account addresses
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
	p.supplyKeeper = supply.NewKeeper(p.cdc, protocol.Keys[protocol.SupplyStoreKey], p.accountKeeper, p.bankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(
		p.cdc, protocol.Keys[staking.StoreKey], protocol.TKeys[staking.TStoreKey],
		p.supplyKeeper, stakingSubspace)
	p.mintKeeper = mint.NewKeeper(p.cdc, protocol.Keys[mint.StoreKey], mintSubspace, &stakingKeeper, p.supplyKeeper, auth.FeeCollectorName)
	p.distrKeeper = distr.NewKeeper(p.cdc, protocol.Keys[distr.StoreKey], distrSubspace, &stakingKeeper,
		p.supplyKeeper, auth.FeeCollectorName, ModuleAccountAddrs())
	p.slashingKeeper = slashing.NewKeeper(
		p.cdc, protocol.Keys[slashing.StoreKey], &stakingKeeper, slashingSubspace)
	p.crisisKeeper = crisis.NewKeeper(crisisSubspace, p.invCheckPeriod, p.supplyKeeper, auth.FeeCollectorName)

	p.cipalKeeper = cipal.NewKeeper(
		protocol.Keys[cipal.StoreKey],
		p.cdc,
		cipalSubspace,
	)

	p.ipalKeeper = ipal.NewKeeper(
		protocol.Keys[ipal.StoreKey],
		p.cdc,
		p.supplyKeeper,
		ipalSubspace,
	)

	p.vmKeeper = vm.NewKeeper(
		p.cdc,
		protocol.Keys[protocol.VMStoreKey],
		vmSubspace,
		p.accountKeeper,
	)

	p.guardianKeeper = guardian.NewKeeper(p.cdc, protocol.Keys[protocol.GuardianStoreKey])

	p.govKeeper = gov.NewKeeper(
		p.cdc, protocol.Keys[gov.StoreKey], govSubspace, p.supplyKeeper,
		&stakingKeeper, p.guardianKeeper, p.protocolKeeper,
	)

	govRouter := gov.NewRouter()
	govRouter.
		AddRoute(gov.RouterKey, gov.NewGovProposalHandler(p.govKeeper)).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(p.paramsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(p.distrKeeper))

	p.govKeeper.SetRouter(govRouter)

	p.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(p.distrKeeper.Hooks(), p.slashingKeeper.Hooks()),
	)

	p.upgradeKeeper = upgrade.NewKeeper(
		p.cdc,
		protocol.Keys[protocol.UpgradeStoreKey],
		p.protocolKeeper,
		p.stakingKeeper)
}

func (p *ProtocolV0) configModuleManager() {
	moduleManager := module.NewManager(
		genaccounts.NewAppModule(p.accountKeeper),
		genutil.NewAppModule(p.accountKeeper, p.stakingKeeper, p.deliverTx),
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
		upgrade.NewAppModule(p.upgradeKeeper),
		guardian.NewAppModule(p.guardianKeeper),
	)

	moduleManager.SetOrderBeginBlockers(
		mint.ModuleName,
		distr.ModuleName,
		slashing.ModuleName)

	moduleManager.SetOrderEndBlockers(
		crisis.ModuleName,
		gov.ModuleName,
		staking.ModuleName,
		ipal.ModuleName,
		vm.ModuleName,
		upgrade.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	moduleManager.SetOrderInitGenesis(
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
		types.ModuleName,
		guardian.ModuleName,
		upgrade.ModuleName,
	)

	p.moduleManager = moduleManager
}

func (p *ProtocolV0) configSimulationManager() {
	slashingModule := slashing.NewAppModule(p.slashingKeeper, p.stakingKeeper)
	slashingModuleP := slashingModule.WithAccountKeeper(p.accountKeeper).WithStakingKeeper(p.stakingKeeper)

	distrModule := distr.NewAppModule(p.distrKeeper, p.supplyKeeper)
	distrModuleP := distrModule.WithAccountKeeper(p.accountKeeper).WithStakingKeeper(p.stakingKeeper)

	govModule := gov.NewAppModule(p.govKeeper, p.supplyKeeper)
	govModuleP := govModule.WithAccountKeeper(p.accountKeeper)

	ipalModule := ipal.NewAppModule(p.ipalKeeper)
	ipalModuleP := ipalModule.WithAccountKeeper(p.accountKeeper)

	vmModule := vm.NewAppModule(p.vmKeeper)
	vmModuleP := vmModule.WithAccountKeeper(p.accountKeeper)

	cipalModule := cipal.NewAppModule(p.cipalKeeper)
	cipalModuleP := cipalModule.WithAccountKeeper(p.accountKeeper)

	simManager := module.NewSimulationManager(
		genaccounts.NewSimAppModule(p.accountKeeper),
		auth.NewAppModule(p.accountKeeper),
		bank.NewAppModule(p.bankKeeper, p.accountKeeper),
		staking.NewAppModule(p.stakingKeeper, p.distrKeeper, p.accountKeeper, p.supplyKeeper),
		slashingModuleP,
		mint.NewAppModule(p.mintKeeper),
		mint.NewAppModule(p.mintKeeper),
		distrModuleP,
		govModuleP,
		ipalModuleP,
		cipalModuleP,
		vmModuleP,
	)
	p.simManager = simManager
}

func (p *ProtocolV0) configRouters() {
	p.moduleManager.RegisterRoutes(p.router, p.queryRouter)
}

// InitChainer initializes application state at genesis as a hook
func (p *ProtocolV0) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState sdk.GenesisState
	p.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	return p.moduleManager.InitGenesis(ctx, genesisState)
}

// BeginBlocker set function to BaseApp as a hook
func (p *ProtocolV0) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return p.moduleManager.BeginBlock(ctx, req)
}

// EndBlocker sets function to BaseApp as a hook
func (p *ProtocolV0) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return p.moduleManager.EndBlock(ctx, req)
}

func (p *ProtocolV0) configFeeHandlers() {
	p.anteHandler = ante.NewAnteHandler(p.accountKeeper, p.supplyKeeper, ante.DefaultSigVerificationGasConsumer)
	p.feeRefundHandler = auth.NewFeeRefundHandler(p.accountKeeper, p.supplyKeeper, p.refundKeeper)
}

//for test

// SetInitChainer
func (p *ProtocolV0) SetInitChainer(initChainer sdk.InitChainer) {
	p.initChainer = initChainer
}

// SetRouter
func (p *ProtocolV0) SetRouter(router sdk.Router) {
	p.router = router
}

// SetQuearyRouter
func (p *ProtocolV0) SetQuearyRouter(queryRouter sdk.QueryRouter) {
	p.queryRouter = queryRouter
}

// SetAnteHandler
func (p *ProtocolV0) SetAnteHandler(anteHandler sdk.AnteHandler) {
	p.anteHandler = anteHandler
}

// GetSimulationManager - for simulation
func (p *ProtocolV0) GetSimulationManager() interface{} {
	return p.simManager
}
