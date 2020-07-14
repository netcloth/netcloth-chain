package v0

import (
	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/netcloth/netcloth-chain/app/protocol"
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

// ModuleBasics is in charge of setting up basic,
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

type ProtocolV0 struct {
	version uint64
	Cdc     *codec.Codec
	logger  log.Logger

	moduleManager *module.Manager

	AccountKeeper  auth.AccountKeeper
	RefundKeeper   auth.RefundKeeper
	BankKeeper     bank.Keeper
	slashingKeeper slashing.Keeper
	mintKeeper     mint.Keeper
	distrKeeper    distr.Keeper
	protocolKeeper sdk.ProtocolKeeper
	GovKeeper      gov.Keeper
	crisisKeeper   crisis.Keeper
	paramsKeeper   params.Keeper
	SupplyKeeper   supply.Keeper
	StakingKeeper  staking.Keeper
	ipalKeeper     ipal.Keeper
	cipalKeeper    cipal.Keeper
	vmKeeper       vm.Keeper
	upgradeKeeper  upgrade.Keeper
	guardianKeeper guardian.Keeper

	router      sdk.Router
	queryRouter sdk.QueryRouter

	anteHandler      sdk.AnteHandler
	feeRefundHandler sdk.FeeRefundHandler
	initChainer      sdk.InitChainer
	deliverTx        genutil.DeliverTxfn

	config *cfg.InstrumentationConfig

	invCheckPeriod uint
}

func (p *ProtocolV0) SetRouter(router sdk.Router) {
	p.router = router
}

func (p *ProtocolV0) SetQuearyRouter(queryRouter sdk.QueryRouter) {
	p.queryRouter = queryRouter
}

func (p *ProtocolV0) SetAnteHandler(anteHandler sdk.AnteHandler) {
	p.anteHandler = anteHandler
}

func (p *ProtocolV0) SetInitChainer(initChainer sdk.InitChainer) {
	p.initChainer = initChainer
}

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

func (p *ProtocolV0) GetVersion() uint64 {
	return p.version
}

func (p *ProtocolV0) GetRouter() sdk.Router {
	return p.router
}

func (p *ProtocolV0) GetQueryRouter() sdk.QueryRouter {
	return p.queryRouter
}

func (p *ProtocolV0) GetAnteHandler() sdk.AnteHandler {
	return p.anteHandler
}

func (p *ProtocolV0) GetFeeRefundHandler() sdk.FeeRefundHandler {
	return p.feeRefundHandler
}

func (p *ProtocolV0) LoadContext() {
	p.configCodec()
	p.configKeepers()
	p.configModuleManager()
	p.configRouters()
	p.configFeeHandlers()
}

func (p *ProtocolV0) Init() {
}

func (p *ProtocolV0) GetCodec() *codec.Codec {
	return p.Cdc
}

func (p *ProtocolV0) GetInitChainer() sdk.InitChainer {
	return p.InitChainer
}

func (p *ProtocolV0) GetBeginBlocker() sdk.BeginBlocker {
	return p.BeginBlocker
}

func (p *ProtocolV0) GetEndBlocker() sdk.EndBlocker {
	return p.EndBlocker
}

func (p *ProtocolV0) configCodec() {
	p.Cdc = MakeCodec()
}

func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)

	return cdc
}

func ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (p *ProtocolV0) configKeepers() {
	p.paramsKeeper = params.NewKeeper(p.Cdc, protocol.Keys[params.StoreKey], protocol.TKeys[params.TStoreKey])
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

	p.AccountKeeper = auth.NewAccountKeeper(p.Cdc, protocol.Keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)
	p.RefundKeeper = auth.NewRefundKeeper(p.Cdc, protocol.Keys[auth.RefundKey])
	p.BankKeeper = bank.NewBaseKeeper(p.AccountKeeper, bankSubspace, ModuleAccountAddrs())
	p.SupplyKeeper = supply.NewKeeper(p.Cdc, protocol.Keys[protocol.SupplyStoreKey], p.AccountKeeper, p.BankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(
		p.Cdc, protocol.Keys[staking.StoreKey], protocol.TKeys[staking.TStoreKey],
		p.SupplyKeeper, stakingSubspace)
	p.mintKeeper = mint.NewKeeper(p.Cdc, protocol.Keys[mint.StoreKey], mintSubspace, &stakingKeeper, p.SupplyKeeper, auth.FeeCollectorName)
	p.distrKeeper = distr.NewKeeper(p.Cdc, protocol.Keys[distr.StoreKey], distrSubspace, &stakingKeeper,
		p.SupplyKeeper, auth.FeeCollectorName, ModuleAccountAddrs())
	p.slashingKeeper = slashing.NewKeeper(
		p.Cdc, protocol.Keys[slashing.StoreKey], &stakingKeeper, slashingSubspace)
	p.crisisKeeper = crisis.NewKeeper(crisisSubspace, p.invCheckPeriod, p.SupplyKeeper, auth.FeeCollectorName)

	p.cipalKeeper = cipal.NewKeeper(
		protocol.Keys[cipal.StoreKey],
		p.Cdc,
		cipalSubspace)

	p.ipalKeeper = ipal.NewKeeper(
		protocol.Keys[ipal.StoreKey],
		p.Cdc,
		p.SupplyKeeper,
		ipalSubspace)

	p.vmKeeper = vm.NewKeeper(
		p.Cdc,
		protocol.Keys[protocol.VMStoreKey],
		protocol.Keys[protocol.VMCodeStoreKey],
		protocol.Keys[protocol.VMLogStoreKey],
		vmSubspace,
		p.AccountKeeper)

	p.guardianKeeper = guardian.NewKeeper(p.Cdc, protocol.Keys[protocol.GuardianStoreKey])

	p.GovKeeper = gov.NewKeeper(
		p.Cdc, protocol.Keys[gov.StoreKey], govSubspace, p.SupplyKeeper,
		&stakingKeeper, p.guardianKeeper, p.protocolKeeper,
	)

	govRouter := gov.NewRouter()
	govRouter.
		AddRoute(gov.RouterKey, gov.NewGovProposalHandler(p.GovKeeper)).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(p.paramsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(p.distrKeeper))

	p.GovKeeper.SetRouter(govRouter)

	p.StakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(p.distrKeeper.Hooks(), p.slashingKeeper.Hooks()),
	)

	p.upgradeKeeper = upgrade.NewKeeper(
		p.Cdc,
		protocol.Keys[protocol.UpgradeStoreKey],
		p.protocolKeeper,
		p.StakingKeeper)
}

func (p *ProtocolV0) configModuleManager() {
	moduleManager := module.NewManager(
		genaccounts.NewAppModule(p.AccountKeeper),
		genutil.NewAppModule(p.AccountKeeper, p.StakingKeeper, p.deliverTx),
		auth.NewAppModule(p.AccountKeeper),
		bank.NewAppModule(p.BankKeeper, p.AccountKeeper),
		crisis.NewAppModule(&p.crisisKeeper),
		supply.NewAppModule(p.SupplyKeeper, p.AccountKeeper),
		distr.NewAppModule(p.distrKeeper, p.SupplyKeeper),
		gov.NewAppModule(p.GovKeeper, p.SupplyKeeper),
		mint.NewAppModule(p.mintKeeper),
		slashing.NewAppModule(p.slashingKeeper, p.StakingKeeper),
		staking.NewAppModule(p.StakingKeeper, p.distrKeeper, p.AccountKeeper, p.SupplyKeeper),
		cipal.NewAppModule(p.cipalKeeper),
		ipal.NewAppModule(p.ipalKeeper),
		vm.NewAppModule(p.vmKeeper),
		upgrade.NewAppModule(p.upgradeKeeper),
		guardian.NewAppModule(p.guardianKeeper),
	)

	moduleManager.SetOrderBeginBlockers(mint.ModuleName, distr.ModuleName, slashing.ModuleName)

	moduleManager.SetOrderEndBlockers(types.ModuleName, crisis.ModuleName, gov.ModuleName, staking.ModuleName, ipal.ModuleName, vm.ModuleName) // TODO upgrade should be the first or the last?

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

func (p *ProtocolV0) configRouters() {
	p.moduleManager.RegisterRoutes(p.router, p.queryRouter)
}

func (p *ProtocolV0) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	p.Cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	return p.moduleManager.InitGenesis(ctx, genesisState)
}

func (p *ProtocolV0) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return p.moduleManager.BeginBlock(ctx, req)
}

func (p *ProtocolV0) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return p.moduleManager.EndBlock(ctx, req)
}

func (p *ProtocolV0) configFeeHandlers() {
	p.anteHandler = ante.NewAnteHandler(p.AccountKeeper, p.SupplyKeeper, ante.DefaultSigVerificationGasConsumer)
	p.feeRefundHandler = auth.NewFeeRefundHandler(p.AccountKeeper, p.SupplyKeeper, p.RefundKeeper)
}

// for test

// for simulation

func (p *ProtocolV0) GetSimulationManager() interface{} {
	return nil
}
