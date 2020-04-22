package test2

import (
	"github.com/netcloth/netcloth-chain/app/v0/gov"
	"github.com/netcloth/netcloth-chain/app/v0/staking"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/netcloth/netcloth-chain/types"
)

func TestTickExpiredDepositPeriod2(t *testing.T) {
	input := getMockApp(t, 10, gov.GenesisState{}, nil)

	header := abci.Header{Height: input.mApp.LastBlockHeight() + 1}
	input.mApp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := input.mApp.BaseApp.NewContext(false, abci.Header{})
	initGenAccount(t, ctx, input.mApp)

	govHandler := gov.NewHandler(input.keeper)

	inactiveQueue := input.keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	newProposalMsg := gov.NewMsgSubmitProposal(
		gov.ContentFromProposalType("test", "test", gov.ProposalTypeText),
		sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 5)},
		input.addrs[0],
	)

	res, err := govHandler(ctx, newProposalMsg)
	require.Nil(t, err)
	require.True(t, res.IsOK())

	inactiveQueue = input.keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	newHeader := ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(time.Duration(1) * time.Second)
	ctx = ctx.WithBlockHeader(newHeader)

	inactiveQueue = input.keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	newHeader = ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(input.keeper.GetDepositParams(ctx).MaxDepositPeriod)
	ctx = ctx.WithBlockHeader(newHeader)

	inactiveQueue = input.keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.True(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	gov.EndBlocker(ctx, input.keeper)

	inactiveQueue = input.keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()
}

func TestTickMultipleExpiredDepositPeriod(t *testing.T) {
	input := getMockApp(t, 10, gov.GenesisState{}, nil)

	header := abci.Header{Height: input.mApp.LastBlockHeight() + 1}
	input.mApp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := input.mApp.BaseApp.NewContext(false, abci.Header{})

	initGenAccount(t, ctx, input.mApp)

	govHandler := gov.NewHandler(input.keeper)

	inactiveQueue := input.keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	newProposalMsg := gov.NewMsgSubmitProposal(
		gov.ContentFromProposalType("test", "test", gov.ProposalTypeText),
		sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 5)},
		input.addrs[0],
	)

	res, err := govHandler(ctx, newProposalMsg)
	require.Nil(t, err)
	require.True(t, res.IsOK())

	inactiveQueue = input.keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	newHeader := ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(time.Duration(2) * time.Second)
	ctx = ctx.WithBlockHeader(newHeader)

	inactiveQueue = input.keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	newProposalMsg2 := gov.NewMsgSubmitProposal(
		gov.ContentFromProposalType("test2", "test2", gov.ProposalTypeText),
		sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 5)},
		input.addrs[0],
	)

	res, err = govHandler(ctx, newProposalMsg2)
	require.Nil(t, err)
	require.True(t, res.IsOK())

	newHeader = ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(input.keeper.GetDepositParams(ctx).MaxDepositPeriod).Add(time.Duration(-1) * time.Second)
	ctx = ctx.WithBlockHeader(newHeader)

	inactiveQueue = input.keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.True(t, inactiveQueue.Valid())
	inactiveQueue.Close()
	gov.EndBlocker(ctx, input.keeper)
	inactiveQueue = input.keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	newHeader = ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(time.Duration(5) * time.Second)
	ctx = ctx.WithBlockHeader(newHeader)

	inactiveQueue = input.keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.True(t, inactiveQueue.Valid())
	inactiveQueue.Close()
	gov.EndBlocker(ctx, input.keeper)
	inactiveQueue = input.keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()
}

func TestTickPassedDepositPeriod(t *testing.T) {
	input := getMockApp(t, 10, gov.GenesisState{}, nil)

	header := abci.Header{Height: input.mApp.LastBlockHeight() + 1}
	input.mApp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := input.mApp.BaseApp.NewContext(false, abci.Header{})

	initGenAccount(t, ctx, input.mApp)

	govHandler := gov.NewHandler(input.keeper)

	inactiveQueue := input.keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()
	activeQueue := input.keeper.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, activeQueue.Valid())
	activeQueue.Close()

	newProposalMsg := gov.NewMsgSubmitProposal(
		gov.ContentFromProposalType("test2", "test2", gov.ProposalTypeText),
		sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 5)},
		input.addrs[0],
	)

	res, err := govHandler(ctx, newProposalMsg)
	require.Nil(t, err)
	require.True(t, res.IsOK())
	var proposalID uint64
	input.keeper.Getcdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &proposalID)

	inactiveQueue = input.keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	newHeader := ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(time.Duration(1) * time.Second)
	ctx = ctx.WithBlockHeader(newHeader)

	inactiveQueue = input.keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	newDepositMsg := gov.NewMsgDeposit(input.addrs[1], proposalID, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 5)})
	res, err = govHandler(ctx, newDepositMsg)
	require.Nil(t, err)
	require.True(t, res.IsOK())

	activeQueue = input.keeper.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, activeQueue.Valid())
	activeQueue.Close()
}

func TestTickPassedVotingPeriod(t *testing.T) {
	input := getMockApp(t, 10, gov.GenesisState{}, nil)
	SortAddresses(input.addrs)

	header := abci.Header{Height: input.mApp.LastBlockHeight() + 1}
	input.mApp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := input.mApp.BaseApp.NewContext(false, abci.Header{})

	initGenAccount(t, ctx, input.mApp)

	govHandler := gov.NewHandler(input.keeper)

	inactiveQueue := input.keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()
	activeQueue := input.keeper.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, activeQueue.Valid())
	activeQueue.Close()

	proposalCoins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(5))}
	newProposalMsg := gov.NewMsgSubmitProposal(testProposal(), proposalCoins, input.addrs[0])

	res, err := govHandler(ctx, newProposalMsg)
	require.Nil(t, err)
	require.True(t, res.IsOK())
	var proposalID uint64
	input.keeper.Getcdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &proposalID)

	newHeader := ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(time.Duration(1) * time.Second)
	ctx = ctx.WithBlockHeader(newHeader)

	newDepositMsg := gov.NewMsgDeposit(input.addrs[1], proposalID, proposalCoins)
	res, err = govHandler(ctx, newDepositMsg)
	require.Nil(t, err)
	require.True(t, res.IsOK())

	newHeader = ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(input.keeper.GetDepositParams(ctx).MaxDepositPeriod).Add(input.keeper.GetVotingParams(ctx).VotingPeriod)
	ctx = ctx.WithBlockHeader(newHeader)

	inactiveQueue = input.keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	activeQueue = input.keeper.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.True(t, activeQueue.Valid())

	var activeProposalID uint64

	require.NoError(t, input.keeper.Getcdc().UnmarshalBinaryLengthPrefixed(activeQueue.Value(), &activeProposalID))
	proposal, ok := input.keeper.GetProposal(ctx, activeProposalID)
	require.True(t, ok)
	require.Equal(t, gov.StatusVotingPeriod, proposal.Status)
	depositsIterator := input.keeper.GetDepositsIterator(ctx, proposalID)
	require.True(t, depositsIterator.Valid())
	depositsIterator.Close()
	activeQueue.Close()

	gov.EndBlocker(ctx, input.keeper)

	activeQueue = input.keeper.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, activeQueue.Valid())
	activeQueue.Close()
}

func TestProposalPassedEndblocker(t *testing.T) {
	input := getMockApp(t, 1, gov.GenesisState{}, nil)
	SortAddresses(input.addrs)

	handler := gov.NewHandler(input.keeper)
	stakingHandler := staking.NewHandler(input.sk)

	header := abci.Header{Height: input.mApp.LastBlockHeight() + 1}
	input.mApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := input.mApp.BaseApp.NewContext(false, abci.Header{})
	initGenAccount(t, ctx, input.mApp)

	valAddr := sdk.ValAddress(input.addrs[0])

	createValidators(t, stakingHandler, ctx, []sdk.ValAddress{valAddr}, []int64{10})
	staking.EndBlocker(ctx, input.sk)

	macc := input.keeper.GetGovernanceAccount(ctx)
	require.NotNil(t, macc)
	initialModuleAccCoins := macc.GetCoins()

	proposal, err := input.keeper.SubmitProposal(ctx, testProposal(), input.addrs[0])
	require.NoError(t, err)

	proposalCoins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(10))}
	newDepositMsg := gov.NewMsgDeposit(input.addrs[0], proposal.ProposalID, proposalCoins)
	res, err := handler(ctx, newDepositMsg)
	require.Nil(t, err)
	require.True(t, res.IsOK())

	macc = input.keeper.GetGovernanceAccount(ctx)
	require.NotNil(t, macc)
	moduleAccCoins := macc.GetCoins()

	deposits := initialModuleAccCoins.Add(proposal.TotalDeposit).Add(proposalCoins)
	require.True(t, moduleAccCoins.IsEqual(deposits))

	err = input.keeper.AddVote(ctx, proposal.ProposalID, input.addrs[0], gov.OptionYes)
	require.NoError(t, err)

	newHeader := ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(input.keeper.GetDepositParams(ctx).MaxDepositPeriod).Add(input.keeper.GetVotingParams(ctx).VotingPeriod)
	ctx = ctx.WithBlockHeader(newHeader)

	gov.EndBlocker(ctx, input.keeper)

	macc = input.keeper.GetGovernanceAccount(ctx)
	require.NotNil(t, macc)
	require.True(t, macc.GetCoins().IsEqual(initialModuleAccCoins))
}

func TestEndBlockerProposalHandlerFailed(t *testing.T) {
	input := getMockApp(t, 1, gov.GenesisState{}, nil)
	SortAddresses(input.addrs)

	// hijack the router to one that will fail in a proposal's handler
	//input.keeper.router = gov.NewRouter().AddRoute(gov.RouterKey, badProposalHandler)
	input.keeper.SetRouter(gov.NewRouter().AddRoute(gov.RouterKey, badProposalHandler))

	handler := gov.NewHandler(input.keeper)
	stakingHandler := staking.NewHandler(input.sk)

	header := abci.Header{Height: input.mApp.LastBlockHeight() + 1}
	input.mApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := input.mApp.BaseApp.NewContext(false, abci.Header{})
	initGenAccount(t, ctx, input.mApp)

	valAddr := sdk.ValAddress(input.addrs[0])

	createValidators(t, stakingHandler, ctx, []sdk.ValAddress{valAddr}, []int64{10})
	staking.EndBlocker(ctx, input.sk)

	// Create a proposal where the handler will pass for the test proposal
	// because the value of contextKeyBadProposal is true.
	ctx = ctx.WithValue(contextKeyBadProposal, true)
	proposal, err := input.keeper.SubmitProposal(ctx, testProposal(), input.addrs[0])
	require.NoError(t, err)

	proposalCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(10)))
	newDepositMsg := gov.NewMsgDeposit(input.addrs[0], proposal.ProposalID, proposalCoins)
	res, err := handler(ctx, newDepositMsg)
	require.Nil(t, err)
	require.True(t, res.IsOK())

	err = input.keeper.AddVote(ctx, proposal.ProposalID, input.addrs[0], gov.OptionYes)
	require.NoError(t, err)

	newHeader := ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(input.keeper.GetDepositParams(ctx).MaxDepositPeriod).Add(input.keeper.GetVotingParams(ctx).VotingPeriod)
	ctx = ctx.WithBlockHeader(newHeader)

	// Set the contextKeyBadProposal value to false so that the handler will fail
	// during the processing of the proposal in the EndBlocker.
	ctx = ctx.WithValue(contextKeyBadProposal, false)

	// validate that the proposal fails/has been rejected
	gov.EndBlocker(ctx, input.keeper)
}
