package test2

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/app/v0/gov"
	"github.com/netcloth/netcloth-chain/app/v0/gov/types"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
)

const custom = "custom"

func getQueriedParams(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) (gov.DepositParams, gov.VotingParams, gov.TallyParams) {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, gov.QuerierRoute, gov.QueryParams, gov.ParamDeposit}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{gov.QueryParams, gov.ParamDeposit}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var depositParams gov.DepositParams
	err2 := cdc.UnmarshalJSON(bz, &depositParams)
	require.Nil(t, err2)

	query = abci.RequestQuery{
		Path: strings.Join([]string{custom, gov.QuerierRoute, gov.QueryParams, gov.ParamVoting}, "/"),
		Data: []byte{},
	}

	bz, err = querier(ctx, []string{gov.QueryParams, gov.ParamVoting}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var votingParams gov.VotingParams
	err2 = cdc.UnmarshalJSON(bz, &votingParams)
	require.Nil(t, err2)

	query = abci.RequestQuery{
		Path: strings.Join([]string{custom, gov.QuerierRoute, gov.QueryParams, gov.ParamTallying}, "/"),
		Data: []byte{},
	}

	bz, err = querier(ctx, []string{gov.QueryParams, gov.ParamTallying}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var tallyParams gov.TallyParams
	err2 = cdc.UnmarshalJSON(bz, &tallyParams)
	require.Nil(t, err2)

	return depositParams, votingParams, tallyParams
}

func getQueriedProposal(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, proposalID uint64) gov.Proposal {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, gov.QuerierRoute, gov.QueryProposal}, "/"),
		Data: cdc.MustMarshalJSON(gov.NewQueryProposalParams(proposalID)),
	}

	bz, err := querier(ctx, []string{gov.QueryProposal}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var proposal gov.Proposal
	err2 := cdc.UnmarshalJSON(bz, proposal)
	require.Nil(t, err2)
	return proposal
}

func getQueriedProposals(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, depositor, voter sdk.AccAddress, status gov.ProposalStatus, limit uint64) []gov.Proposal {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, gov.QuerierRoute, gov.QueryProposals}, "/"),
		Data: cdc.MustMarshalJSON(gov.NewQueryProposalsParams(status, limit, voter, depositor)),
	}

	bz, err := querier(ctx, []string{gov.QueryProposals}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var proposals gov.Proposals
	err2 := cdc.UnmarshalJSON(bz, &proposals)
	require.Nil(t, err2)
	return proposals
}

func getQueriedDeposit(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, proposalID uint64, depositor sdk.AccAddress) gov.Deposit {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, gov.QuerierRoute, gov.QueryDeposit}, "/"),
		Data: cdc.MustMarshalJSON(gov.NewQueryDepositParams(proposalID, depositor)),
	}

	bz, err := querier(ctx, []string{gov.QueryDeposit}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var deposit gov.Deposit
	err2 := cdc.UnmarshalJSON(bz, &deposit)
	require.Nil(t, err2)
	return deposit
}

func getQueriedDeposits(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, proposalID uint64) []gov.Deposit {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, gov.QuerierRoute, gov.QueryDeposits}, "/"),
		Data: cdc.MustMarshalJSON(gov.NewQueryProposalParams(proposalID)),
	}

	bz, err := querier(ctx, []string{gov.QueryDeposits}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var deposits []gov.Deposit
	err2 := cdc.UnmarshalJSON(bz, &deposits)
	require.Nil(t, err2)
	return deposits
}

func getQueriedVote(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, proposalID uint64, voter sdk.AccAddress) gov.Vote {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, gov.QuerierRoute, gov.QueryVote}, "/"),
		Data: cdc.MustMarshalJSON(gov.NewQueryVoteParams(proposalID, voter)),
	}

	bz, err := querier(ctx, []string{gov.QueryVote}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var vote gov.Vote
	err2 := cdc.UnmarshalJSON(bz, &vote)
	require.Nil(t, err2)
	return vote
}

func getQueriedVotes(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, proposalID uint64) []gov.Vote {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, gov.QuerierRoute, gov.QueryVote}, "/"),
		Data: cdc.MustMarshalJSON(gov.NewQueryProposalParams(proposalID)),
	}

	bz, err := querier(ctx, []string{gov.QueryVotes}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var votes []gov.Vote
	err2 := cdc.UnmarshalJSON(bz, &votes)
	require.Nil(t, err2)
	return votes
}

func getQueriedTally(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, proposalID uint64) gov.TallyResult {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, gov.QuerierRoute, gov.QueryTally}, "/"),
		Data: cdc.MustMarshalJSON(gov.NewQueryProposalParams(proposalID)),
	}

	bz, err := querier(ctx, []string{gov.QueryTally}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var tally gov.TallyResult
	err2 := cdc.UnmarshalJSON(bz, &tally)
	require.Nil(t, err2)
	return tally
}

func TestQueryParams(t *testing.T) {
	cdc := codec.New()
	input := getMockApp(t, 1000, gov.GenesisState{}, nil)
	querier := gov.NewQuerier(input.keeper)

	header := abci.Header{Height: input.mApp.LastBlockHeight() + 1}
	input.mApp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := input.mApp.NewContext(false, abci.Header{})

	getQueriedParams(t, ctx, cdc, querier)
}

func TestQueries(t *testing.T) {
	cdc := codec.New()
	input := getMockApp(t, 1000, gov.GenesisState{}, nil)
	querier := gov.NewQuerier(input.keeper)
	handler := gov.NewHandler(input.keeper)

	types.RegisterCodec(cdc)

	header := abci.Header{Height: input.mApp.LastBlockHeight() + 1}
	input.mApp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := input.mApp.NewContext(false, abci.Header{})

	initGenAccount(t, ctx, input.mApp)

	depositParams, _, _ := getQueriedParams(t, ctx, cdc, querier)

	// input.addrs[0] proposes (and deposits) proposals #1 and #2
	res, err := handler(ctx, gov.NewMsgSubmitProposal(testProposal(), sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 1)}, input.addrs[0]))
	require.Nil(t, err)
	var proposalID1 uint64
	require.True(t, res.IsOK())
	cdc.MustUnmarshalBinaryLengthPrefixed(res.Data, &proposalID1)

	res, err = handler(ctx, gov.NewMsgSubmitProposal(testProposal(), sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000000)}, input.addrs[0]))
	require.Nil(t, err)
	var proposalID2 uint64
	require.True(t, res.IsOK())
	cdc.MustUnmarshalBinaryLengthPrefixed(res.Data, &proposalID2)

	// input.addrs[1] proposes (and deposits) proposals #3
	res, err = handler(ctx, gov.NewMsgSubmitProposal(testProposal(), sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 1)}, input.addrs[1]))
	require.Nil(t, err)
	var proposalID3 uint64
	require.True(t, res.IsOK())
	cdc.MustUnmarshalBinaryLengthPrefixed(res.Data, &proposalID3)

	// input.addrs[1] deposits on proposals #2 & #3
	res, err = handler(ctx, gov.NewMsgDeposit(input.addrs[1], proposalID2, depositParams.MinDeposit))
	require.Nil(t, err)
	res, err = handler(ctx, gov.NewMsgDeposit(input.addrs[1], proposalID3, depositParams.MinDeposit))
	require.Nil(t, err)

	// check deposits on proposal1 match individual deposits
	deposits := getQueriedDeposits(t, ctx, cdc, querier, proposalID1)
	require.Len(t, deposits, 1)
	deposit := getQueriedDeposit(t, ctx, cdc, querier, proposalID1, input.addrs[0])
	require.Equal(t, deposit, deposits[0])

	// check deposits on proposal2 match individual deposits
	deposits = getQueriedDeposits(t, ctx, cdc, querier, proposalID2)
	require.Len(t, deposits, 2)
	deposit = getQueriedDeposit(t, ctx, cdc, querier, proposalID2, input.addrs[0])
	require.True(t, deposit.Equals(deposits[0]))
	deposit = getQueriedDeposit(t, ctx, cdc, querier, proposalID2, input.addrs[1])
	require.True(t, deposit.Equals(deposits[1]))

	// check deposits on proposal3 match individual deposits
	deposits = getQueriedDeposits(t, ctx, cdc, querier, proposalID3)
	require.Len(t, deposits, 1)
	deposit = getQueriedDeposit(t, ctx, cdc, querier, proposalID3, input.addrs[1])
	require.Equal(t, deposit, deposits[0])

	// Only proposal #1 should be in Deposit Period
	proposals := getQueriedProposals(t, ctx, cdc, querier, nil, nil, gov.StatusDepositPeriod, 0)
	require.Len(t, proposals, 1)
	require.Equal(t, proposalID1, proposals[0].ProposalID)

	// Only proposals #2 and #3 should be in Voting Period
	proposals = getQueriedProposals(t, ctx, cdc, querier, nil, nil, gov.StatusVotingPeriod, 0)
	require.Len(t, proposals, 2)
	require.Equal(t, proposalID2, proposals[0].ProposalID)
	require.Equal(t, proposalID3, proposals[1].ProposalID)

	// Addrs[0] votes on proposals #2 & #3
	res, err = handler(ctx, gov.NewMsgVote(input.addrs[0], proposalID2, gov.OptionYes))
	require.Nil(t, err)
	require.True(t, res.IsOK())

	res, err = handler(ctx, gov.NewMsgVote(input.addrs[0], proposalID3, gov.OptionYes))
	require.Nil(t, err)
	require.True(t, res.IsOK())

	// Addrs[1] votes on proposal #3
	handler(ctx, gov.NewMsgVote(input.addrs[1], proposalID3, gov.OptionYes))

	// Test query voted by input.addrs[0]
	proposals = getQueriedProposals(t, ctx, cdc, querier, nil, input.addrs[0], gov.StatusNil, 0)
	require.Equal(t, proposalID2, (proposals[0]).ProposalID)
	require.Equal(t, proposalID3, (proposals[1]).ProposalID)

	// Test query votes on Proposal 2
	votes := getQueriedVotes(t, ctx, cdc, querier, proposalID2)
	require.Len(t, votes, 1)
	require.Equal(t, input.addrs[0], votes[0].Voter)

	vote := getQueriedVote(t, ctx, cdc, querier, proposalID2, input.addrs[0])
	require.Equal(t, vote, votes[0])

	// Test query votes on Proposal 3
	votes = getQueriedVotes(t, ctx, cdc, querier, proposalID3)
	require.Len(t, votes, 2)
	require.True(t, input.addrs[0].String() == votes[0].Voter.String())
	require.True(t, input.addrs[1].String() == votes[1].Voter.String())

	// Test proposals queries with filters

	// Test query all proposals
	proposals = getQueriedProposals(t, ctx, cdc, querier, nil, nil, gov.StatusNil, 0)
	require.Equal(t, proposalID1, (proposals[0]).ProposalID)
	require.Equal(t, proposalID2, (proposals[1]).ProposalID)
	require.Equal(t, proposalID3, (proposals[2]).ProposalID)

	// Test query voted by input.addrs[1]
	proposals = getQueriedProposals(t, ctx, cdc, querier, nil, input.addrs[1], gov.StatusNil, 0)
	require.Equal(t, proposalID3, (proposals[0]).ProposalID)

	// Test query deposited by input.addrs[0]
	proposals = getQueriedProposals(t, ctx, cdc, querier, input.addrs[0], nil, gov.StatusNil, 0)
	require.Equal(t, proposalID1, (proposals[0]).ProposalID)

	// Test query deposited by addr2
	proposals = getQueriedProposals(t, ctx, cdc, querier, input.addrs[1], nil, gov.StatusNil, 0)
	require.Equal(t, proposalID2, (proposals[0]).ProposalID)
	require.Equal(t, proposalID3, (proposals[1]).ProposalID)

	// Test query voted AND deposited by addr1
	proposals = getQueriedProposals(t, ctx, cdc, querier, input.addrs[0], input.addrs[0], gov.StatusNil, 0)
	require.Equal(t, proposalID2, (proposals[0]).ProposalID)
}
