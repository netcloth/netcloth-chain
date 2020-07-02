package gov_test

import (
	"bytes"
	"errors"
	"sort"
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
	"github.com/netcloth/netcloth-chain/app/v0/gov/types"
	"github.com/netcloth/netcloth-chain/app/v0/staking"
	"github.com/netcloth/netcloth-chain/app/v0/supply"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

var (
	valTokens = sdk.TokensFromConsensusPower(1000)
	valCoins  = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, valTokens))
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

func getMockApp(t *testing.T, numGenAccs int, genState gov.GenesisState, genAccs []auth.Account) testInput {
	mApp := NewNCHApp(t)

	var (
		addrs    []sdk.AccAddress
		pubKeys  []crypto.PubKey
		privKeys []crypto.PrivKey
	)

	if len(genAccs) == 0 {
		genAccs, addrs, pubKeys, privKeys = mock.CreateGenAccounts(numGenAccs, valCoins)
	}

	protocolV0 := getProtocolV0(t, mApp)

	err := setGenesis(mApp, protocolV0.Cdc, genAccs, genState)
	require.Nil(t, err)

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

func setGenesis(app *mock.NCHApp, cdc *codec.Codec, accs []auth.Account, genState gov.GenesisState) error {
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

func setTotalSupply(t *testing.T, app *mock.NCHApp, ctx sdk.Context, accAmt sdk.Int, totalAccounts int) {
	p0 := getProtocolV0(t, app)
	totalSupply := sdk.NewCoins(sdk.NewCoin(p0.StakingKeeper.BondDenom(ctx), accAmt.MulRaw(int64(totalAccounts))))
	prevSupply := p0.SupplyKeeper.GetSupply(ctx)
	p0.SupplyKeeper.SetSupply(ctx, supply.NewSupply(prevSupply.GetTotal().Add(totalSupply)))
}

func initGenAccount(t *testing.T, ctx sdk.Context, app *mock.NCHApp) {
	p0 := getProtocolV0(t, app)
	accAmt := sdk.NewInt(0)
	for _, genAcc := range app.GenesisAccounts {
		acc := p0.AccountKeeper.NewAccountWithAddress(ctx, genAcc.GetAddress())
		acc.SetCoins(genAcc.GetCoins())
		p0.AccountKeeper.SetAccount(ctx, acc)
		accAmt = accAmt.Add(genAcc.GetCoins()[0].Amount)
	}

	setTotalSupply(t, app, ctx, accAmt, 1)
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
	byteAddrs := make([][]byte, 0, len(addrs))
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

////exports of gov

const (
	ProposalTypeText    = gov.ProposalTypeText
	StatusVotingPeriod  = gov.StatusVotingPeriod
	StatusDepositPeriod = gov.StatusDepositPeriod
	StatusRejected      = gov.StatusRejected
	OptionYes           = gov.OptionYes
	RouterKey           = gov.RouterKey
	OptionAbstain       = gov.OptionAbstain
	OptionNoWithVeto    = gov.OptionNoWithVeto
	QuerierRoute        = gov.QuerierRoute
	QueryParams         = gov.QueryParams
	ParamDeposit        = gov.ParamDeposit
	OptionNo            = gov.OptionNo
	ParamVoting         = gov.ParamVoting
	ParamTallying       = gov.ParamTallying
	QueryProposal       = gov.QueryProposal
	QueryProposals      = gov.QueryProposals
	QueryDeposit        = gov.QueryDeposit
	QueryDeposits       = gov.QueryDeposits
	QueryVote           = gov.QueryVote
	QueryVotes          = gov.QueryVotes
	QueryTally          = gov.QueryTally
	StatusNil           = gov.StatusNil
)

type (
	GenesisState   = gov.GenesisState
	Proposal       = gov.Proposal
	Content        = gov.Content
	DepositParams  = gov.DepositParams
	VotingParams   = gov.VotingParams
	TallyParams    = gov.TallyParams
	Proposals      = gov.Proposals
	Deposit        = gov.Deposit
	Vote           = gov.Vote
	TallyResult    = gov.TallyResult
	ProposalStatus = gov.ProposalStatus
)

var (
	ContentFromProposalType    = gov.ContentFromProposalType
	NewQueryProposalParams     = gov.NewQueryProposalParams
	NewQueryProposalsParams    = gov.NewQueryProposalsParams
	EndBlocker                 = gov.EndBlocker
	NewMsgDeposit              = gov.NewMsgDeposit
	NewMsgSubmitProposal       = gov.NewMsgSubmitProposal
	NewHandler                 = gov.NewHandler
	NewRouter                  = gov.NewRouter
	ExportGenesis              = gov.ExportGenesis
	NewQueryVoteParams         = gov.NewQueryVoteParams
	ErrNoProposalHandlerExists = gov.ErrNoProposalHandlerExists
	Tally                      = gov.Tally
	NewQueryDepositParams      = gov.NewQueryDepositParams
	NewQuerier                 = gov.NewQuerier
	NewMsgVote                 = gov.NewMsgVote
	EmptyTallyResult           = gov.EmptyTallyResult
)
