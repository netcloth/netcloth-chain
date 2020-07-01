package staking_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/netcloth/netcloth-chain/app/mock"
	v0 "github.com/netcloth/netcloth-chain/app/mock/p0"
	"github.com/netcloth/netcloth-chain/app/protocol"
	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"github.com/netcloth/netcloth-chain/app/v0/gov"
	"github.com/netcloth/netcloth-chain/app/v0/staking"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
)

var (
	bondedTokens = sdk.TokensFromConsensusPower(1000)
	initTokens   = sdk.TokensFromConsensusPower(100000)
	bondCoin     = sdk.NewCoin(sdk.DefaultBondDenom, bondedTokens)
	initCoin     = sdk.NewCoin(sdk.DefaultBondDenom, initTokens)
	initCoins    = sdk.NewCoins(initCoin)
)

var (
	pubkeys = []crypto.PubKey{
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
	}

	testDescription     = staking.NewDescription("T", "E", "S", "T")
	testCommissionRates = staking.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
)

type testInput struct {
	mApp     *mock.NCHApp
	keeper   gov.Keeper
	sk       staking.Keeper
	ak       auth.AccountKeeper
	addrs    []sdk.AccAddress
	pubKeys  []crypto.PubKey
	privKeys []crypto.PrivKey
}

func getProtocolV0(t *testing.T, app *mock.NCHApp) *v0.ProtocolV0 {
	curProtocol := app.Engine.GetCurrentProtocol()
	protocolV0, ok := curProtocol.(*v0.ProtocolV0)
	require.True(t, ok)
	return protocolV0
}

func getMockApp(t *testing.T, numGenAccs int, genAccs []auth.Account) testInput {
	mApp := NewNCHApp(t)

	var (
		addrs    []sdk.AccAddress
		pubKeys  []crypto.PubKey
		privKeys []crypto.PrivKey
	)

	if genAccs == nil || len(genAccs) == 0 {
		genAccs, addrs, pubKeys, privKeys = mock.CreateGenAccounts(numGenAccs, initCoins)
	}

	protocolV0 := getProtocolV0(t, mApp)

	err := setGenesis(mApp, protocolV0.Cdc, genAccs)
	require.Nil(t, err)

	header := abci.Header{Height: mApp.LastBlockHeight() + 1}
	mApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mApp.BaseApp.NewContext(false, abci.Header{})
	initGenAccount(t, ctx, mApp)
	mApp.Commit()

	return testInput{mApp, protocolV0.GovKeeper, protocolV0.StakingKeeper, protocolV0.AccountKeeper, addrs, pubKeys, privKeys}
}

func NewNCHApp(t *testing.T) *mock.NCHApp {
	logger := log.NewNopLogger()
	db := dbm.NewMemDB()
	baseApp := mock.NewBaseApp("nchmock", logger, db)

	baseApp.SetCommitMultiStoreTracer(nil)
	baseApp.SetAppVersion("v0")

	protocolKeeper := sdk.NewProtocolKeeper(protocol.Keys[protocol.MainStoreKey])
	engine := protocol.NewProtocolEngine(protocolKeeper)
	baseApp.SetProtocolEngine(&engine)

	baseApp.MountKVStores(protocol.Keys)
	baseApp.MountTransientStores(protocol.TKeys)

	err := baseApp.LoadLatestVersion(protocol.Keys[protocol.MainStoreKey])
	require.Nil(t, err)

	engine.Add(v0.NewProtocolV0(0, logger, protocolKeeper, baseApp.DeliverTx, 10, nil))

	engine.LoadProtocol(0)

	baseApp.TxDecoder = auth.DefaultTxDecoder(engine.GetCurrentProtocol().GetCodec())

	return &mock.NCHApp{BaseApp: baseApp}
}

func setGenesis(app *mock.NCHApp, cdc *codec.Codec, accs []auth.Account) error {
	app.GenesisAccounts = accs

	genesisState := v0.NewDefaultGenesisState()

	stateBytes, err := codec.MarshalJSONIndent(cdc, genesisState)
	if err != nil {
		return err
	}

	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	app.Commit()

	return nil
}

func initGenAccount(t *testing.T, ctx sdk.Context, app *mock.NCHApp) {
	p0 := getProtocolV0(t, app)
	for _, genAcc := range app.GenesisAccounts {
		acc := p0.AccountKeeper.NewAccountWithAddress(ctx, genAcc.GetAddress())
		acc.SetCoins(genAcc.GetCoins())
		p0.AccountKeeper.SetAccount(ctx, acc)
	}
}
