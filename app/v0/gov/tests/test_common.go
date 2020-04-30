package tests

import (
	"bytes"
	"errors"
	"github.com/netcloth/netcloth-chain/app/protocol"
	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"github.com/netcloth/netcloth-chain/app/v0/gov"
	"github.com/netcloth/netcloth-chain/app/v0/gov/types"
	"github.com/netcloth/netcloth-chain/app/v0/staking"
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/simapp/app"
	v0 "github.com/netcloth/netcloth-chain/simapp/p0"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"sort"
	"testing"
)

var (
	valTokens  = sdk.TokensFromConsensusPower(1000)
	initTokens = sdk.TokensFromConsensusPower(100000)
	valCoins   = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, valTokens))
	initCoins  = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens))
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
	mApp     *app.NCHApp
	keeper   gov.Keeper
	sk       staking.Keeper
	ak       auth.AccountKeeper
	addrs    []sdk.AccAddress
	pubKeys  []crypto.PubKey
	privKeys []crypto.PrivKey
}

func getProtocolV0(t *testing.T, app *app.NCHApp) *v0.ProtocolV0 {
	curProtocol := app.Engine.GetCurrentProtocol()
	protocolV0, ok := curProtocol.(*v0.ProtocolV0)
	require.True(t, ok)
	return protocolV0
}

func getMockApp(t *testing.T, numGenAccs int, genState gov.GenesisState, genAccs []auth.Account) testInput {
	mApp := NewNCHApp(t)

	var (
		addrs    []sdk.AccAddress
		pubKeys  []crypto.PubKey
		privKeys []crypto.PrivKey
	)

	if genAccs == nil || len(genAccs) == 0 {
		genAccs, addrs, pubKeys, privKeys = CreateGenAccounts(numGenAccs, valCoins)
	}

	protocolV0 := getProtocolV0(t, mApp)

	err := setGenesis(mApp, protocolV0.Cdc, genAccs, genState)
	require.Nil(t, err)

	return testInput{mApp, protocolV0.GovKeeper, protocolV0.StakingKeeper, protocolV0.AccountKeeper, addrs, pubKeys, privKeys}
}

func NewNCHApp(t *testing.T) *app.NCHApp {
	logger := log.NewNopLogger()
	db := dbm.NewMemDB()
	baseApp := app.NewBaseApp("nchmock", logger, db)

	baseApp.SetCommitMultiStoreTracer(nil)
	baseApp.SetAppVersion("v0")

	protocolKeeper := sdk.NewProtocolKeeper(protocol.MainKVStoreKey)
	engine := protocol.NewProtocolEngine(protocolKeeper)
	baseApp.SetProtocolEngine(&engine)

	baseApp.MountStore(protocol.MainKVStoreKey, sdk.StoreTypeDB)

	baseApp.MountKVStores(protocol.Keys)
	baseApp.MountTransientStores(protocol.TKeys)

	err := baseApp.LoadLatestVersion(protocol.MainKVStoreKey)
	require.Nil(t, err)

	engine.Add(v0.NewProtocolV0(0, logger, protocolKeeper, baseApp.DeliverTx, 10, nil))

	engine.LoadProtocol(0)

	baseApp.TxDecoder = auth.DefaultTxDecoder(engine.GetCurrentProtocol().GetCodec())

	return &app.NCHApp{BaseApp: baseApp}
}

func setGenesis(app *app.NCHApp, cdc *codec.Codec, accs []auth.Account, genState gov.GenesisState) error {
	app.GenesisAccounts = accs

	genesisState := v0.NewDefaultGenesisState()
	if !genState.IsEmpty() {
		govState := cdc.MustMarshalJSON(genState)
		genesisState["gov"] = govState
	}

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

func initGenAccount(t *testing.T, ctx sdk.Context, app *app.NCHApp) {
	p0 := getProtocolV0(t, app)
	for _, genAcc := range app.GenesisAccounts {
		acc := p0.AccountKeeper.NewAccountWithAddress(ctx, genAcc.GetAddress())
		acc.SetCoins(genAcc.GetCoins())
		p0.AccountKeeper.SetAccount(ctx, acc)
	}
}

// implement `Interface` in sort package.
type sortByteArrays [][]byte

func (b sortByteArrays) Len() int {
	return len(b)
}

func (b sortByteArrays) Less(i, j int) bool {
	// bytes package already implements Comparable for []byte.
	switch bytes.Compare(b[i], b[j]) {
	case -1:
		return true
	case 0, 1:
		return false
	default:
		return false
	}
}

func (b sortByteArrays) Swap(i, j int) {
	b[j], b[i] = b[i], b[j]
}

// Public
func SortByteArrays(src [][]byte) [][]byte {
	sorted := sortByteArrays(src)
	sort.Sort(sorted)
	return sorted
}

// Sorts Addresses
func SortAddresses(addrs []sdk.AccAddress) {
	var byteAddrs [][]byte
	for _, addr := range addrs {
		byteAddrs = append(byteAddrs, addr.Bytes())
	}
	SortByteArrays(byteAddrs)
	for i, byteAddr := range byteAddrs {
		addrs[i] = byteAddr
	}
}

func testProposal() gov.Content {
	return gov.NewTextProposal("Test", "description")
}

func createValidators(t *testing.T, stakingHandler sdk.Handler, ctx sdk.Context, addrs []sdk.ValAddress, powerAmt []int64) {
	require.True(t, len(addrs) <= len(pubkeys), "Not enough pubkeys specified at top of file.")

	for i := 0; i < len(addrs); i++ {

		valTokens := sdk.TokensFromConsensusPower(powerAmt[i])
		valCreateMsg := staking.NewMsgCreateValidator(
			addrs[i], pubkeys[i], sdk.NewCoin(sdk.DefaultBondDenom, valTokens),
			testDescription, testCommissionRates, sdk.OneInt(),
		)

		res, err := stakingHandler(ctx, valCreateMsg)
		require.NoError(t, err)
		require.NotNil(t, res)
	}
}

const contextKeyBadProposal = "contextKeyBadProposal"

func badProposalHandler(ctx sdk.Context, c gov.Content, pid uint64, proposer sdk.AccAddress) error {
	switch c.ProposalType() {
	case gov.ProposalTypeText, gov.ProposalTypeSoftwareUpgrade:
		v := ctx.Value(contextKeyBadProposal)

		if v == nil || !v.(bool) {
			return errors.New("proposal failed")
		}

		return nil

	default:
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized gov proposal type: %s", c.ProposalType())
	}
}

func ProposalEqual(proposalA gov.Proposal, proposalB gov.Proposal) bool {
	return bytes.Equal(types.ModuleCdc.MustMarshalBinaryBare(proposalA),
		types.ModuleCdc.MustMarshalBinaryBare(proposalB))
}
