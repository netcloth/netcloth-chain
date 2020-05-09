package staking

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtypes "github.com/tendermint/tendermint/types"

	keep "github.com/netcloth/netcloth-chain/app/v0/staking/keeper"
	"github.com/netcloth/netcloth-chain/app/v0/staking/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// retrieve params which are instant
func setInstantUnbondPeriod(keeper keep.Keeper, ctx sdk.Context) types.Params {
	params := keeper.GetParams(ctx)
	params.UnbondingTime = 0
	keeper.SetParams(ctx, params)
	return params
}

func TestValidatorByPowerIndex(t *testing.T) {
	validatorAddr, validatorAddr3 := sdk.ValAddress(keep.Addrs[0]), sdk.ValAddress(keep.Addrs[1])

	initPower := int64(1000000)
	initBond := sdk.TokensFromConsensusPower(initPower)
	ctx, _, keeper, _ := keep.CreateTestInput(t, false, initPower)
	_ = setInstantUnbondPeriod(keeper, ctx)

	// create validator
	msgCreateValidator := NewTestMsgCreateValidator(validatorAddr, keep.PKs[0], initBond)
	_, err := handleMsgCreateValidator(ctx, msgCreateValidator, keeper)
	require.Nil(t, err, "expected create-validator to be ok, got %v", err)

	// must end-block
	updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.Equal(t, 1, len(updates))

	// verify the self-delegation exists
	bond, found := keeper.GetDelegation(ctx, sdk.AccAddress(validatorAddr), validatorAddr)
	require.True(t, found)
	gotBond := bond.Shares.RoundInt()
	require.Equal(t, initBond, gotBond)

	// verify that the by power index exists
	validator, found := keeper.GetValidator(ctx, validatorAddr)
	require.True(t, found)
	power := GetValidatorsByPowerIndexKey(validator)
	require.True(t, keep.ValidatorByPowerIndexExists(ctx, keeper, power))

	// create a second validator keep it bonded
	msgCreateValidator = NewTestMsgCreateValidator(validatorAddr3, keep.PKs[2], initBond)
	_, err = handleMsgCreateValidator(ctx, msgCreateValidator, keeper)
	require.Nil(t, err, "expected create-validator to be ok, got %v", err)

	// must end-block
	updates = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.Equal(t, 1, len(updates))

	// slash and jail the first validator
	consAddr0 := sdk.ConsAddress(keep.PKs[0].Address())
	keeper.Slash(ctx, consAddr0, 0, initPower, sdk.NewDecWithPrec(5, 1))
	keeper.Jail(ctx, consAddr0)
	keeper.ApplyAndReturnValidatorSetUpdates(ctx)

	validator, found = keeper.GetValidator(ctx, validatorAddr)
	require.True(t, found)
	require.Equal(t, sdk.Unbonding, validator.Status)      // ensure is unbonding
	require.Equal(t, initBond.QuoRaw(2), validator.Tokens) // ensure tokens slashed
	keeper.Unjail(ctx, consAddr0)

	// the old power record should have been deleted as the power changed
	require.False(t, keep.ValidatorByPowerIndexExists(ctx, keeper, power))

	// but the new power record should have been created
	validator, found = keeper.GetValidator(ctx, validatorAddr)
	require.True(t, found)
	power2 := GetValidatorsByPowerIndexKey(validator)
	require.True(t, keep.ValidatorByPowerIndexExists(ctx, keeper, power2))

	// now the new record power index should be the same as the original record
	power3 := GetValidatorsByPowerIndexKey(validator)
	require.Equal(t, power2, power3)

	// unbond self-delegation
	totalBond := validator.TokensFromShares(bond.GetShares()).TruncateInt()
	unbondAmt := sdk.NewCoin(sdk.DefaultBondDenom, totalBond)
	msgUndelegate := NewMsgUndelegate(sdk.AccAddress(validatorAddr), validatorAddr, unbondAmt)

	got, err := handleMsgUndelegate(ctx, msgUndelegate, keeper)
	require.Nil(t, err, "expected msg to be ok, got %v", err)

	var finishTime time.Time
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(got.Data, &finishTime)

	ctx = ctx.WithBlockTime(finishTime)
	EndBlocker(ctx, keeper)
	EndBlocker(ctx, keeper)

	// verify that by power key nolonger exists
	_, found = keeper.GetValidator(ctx, validatorAddr)
	require.False(t, found)
	require.False(t, keep.ValidatorByPowerIndexExists(ctx, keeper, power3))
}

func TestDuplicatesMsgCreateValidator(t *testing.T) {
	ctx, _, keeper, _ := keep.CreateTestInput(t, false, 1000)

	addr1, addr2 := sdk.ValAddress(keep.Addrs[0]), sdk.ValAddress(keep.Addrs[1])
	pk1, pk2 := keep.PKs[0], keep.PKs[1]

	valTokens := sdk.TokensFromConsensusPower(10)
	msgCreateValidator1 := NewTestMsgCreateValidator(addr1, pk1, valTokens)
	_, err := handleMsgCreateValidator(ctx, msgCreateValidator1, keeper)
	require.Nil(t, err, "%v", err)

	keeper.ApplyAndReturnValidatorSetUpdates(ctx)

	validator, found := keeper.GetValidator(ctx, addr1)
	require.True(t, found)
	assert.Equal(t, sdk.Bonded, validator.Status)
	assert.Equal(t, addr1, validator.OperatorAddress)
	assert.Equal(t, pk1, validator.ConsPubKey)
	assert.Equal(t, valTokens, validator.BondedTokens())
	assert.Equal(t, valTokens.ToDec(), validator.DelegatorShares)
	assert.Equal(t, Description{}, validator.Description)
	assert.Equal(t, valTokens.ToDec(), validator.SelfDelegation)

	// two validators can't have the same operator address
	msgCreateValidator2 := NewTestMsgCreateValidator(addr1, pk2, valTokens)
	_, err = handleMsgCreateValidator(ctx, msgCreateValidator2, keeper)
	require.NotNil(t, err, "%v", err)
	require.Equal(t, err, ErrValidatorOwnerExists)

	// two validators can't have the same pubkey
	msgCreateValidator3 := NewTestMsgCreateValidator(addr2, pk1, valTokens)
	_, err = handleMsgCreateValidator(ctx, msgCreateValidator3, keeper)
	require.NotNil(t, err, "%v", err)
	require.Equal(t, err, ErrValidatorPubKeyExists)

	// must have different pubkey and operator
	msgCreateValidator4 := NewTestMsgCreateValidator(addr2, pk2, valTokens)
	_, err = handleMsgCreateValidator(ctx, msgCreateValidator4, keeper)
	require.Nil(t, err, "%v", err)

	// must end-block
	updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.Equal(t, 1, len(updates))

	validator, found = keeper.GetValidator(ctx, addr2)
	require.True(t, found)
	assert.Equal(t, sdk.Bonded, validator.Status)
	assert.Equal(t, addr2, validator.OperatorAddress)
	assert.Equal(t, pk2, validator.ConsPubKey)
	assert.True(sdk.IntEq(t, valTokens, validator.Tokens))
	assert.True(sdk.DecEq(t, valTokens.ToDec(), validator.DelegatorShares))
	assert.Equal(t, Description{}, validator.Description)
	assert.Equal(t, valTokens.ToDec(), validator.SelfDelegation)
}

func TestInvalidPubKeyTypeMsgCreateValidator(t *testing.T) {
	ctx, _, keeper, _ := keep.CreateTestInput(t, false, 1000)

	addr := sdk.ValAddress(keep.Addrs[0])
	invalidPk := secp256k1.GenPrivKey().PubKey()

	// invalid pukKey type should not be allowed
	msgCreateValidator := NewTestMsgCreateValidator(addr, invalidPk, sdk.NewInt(10))
	_, err := handleMsgCreateValidator(ctx, msgCreateValidator, keeper)
	require.NotNil(t, err, "%v", err)

	ctx = ctx.WithConsensusParams(&abci.ConsensusParams{
		Validator: &abci.ValidatorParams{PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeSecp256k1}},
	})

	_, err = handleMsgCreateValidator(ctx, msgCreateValidator, keeper)
	require.Nil(t, err, "%v", err)
}

func TestLegacyValidatorDelegations(t *testing.T) {
	ctx, _, keeper, _ := keep.CreateTestInput(t, false, int64(1000))
	setInstantUnbondPeriod(keeper, ctx)

	bondAmount := sdk.TokensFromConsensusPower(10)
	valAddr := sdk.ValAddress(keep.Addrs[0])
	valConsPubKey, valConsAddr := keep.PKs[0], sdk.ConsAddress(keep.PKs[0].Address())
	delAddr := keep.Addrs[1]

	// create validator
	msgCreateVal := NewTestMsgCreateValidator(valAddr, valConsPubKey, bondAmount)
	_, err := handleMsgCreateValidator(ctx, msgCreateVal, keeper)
	require.Nil(t, err, "expected create validator msg to be ok, got %v", err)

	// must end-block
	updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.Equal(t, 1, len(updates))

	// verify the validator exists and has the correct attributes
	validator, found := keeper.GetValidator(ctx, valAddr)
	require.True(t, found)
	require.Equal(t, sdk.Bonded, validator.Status)
	require.Equal(t, bondAmount, validator.DelegatorShares.RoundInt())
	require.Equal(t, bondAmount, validator.BondedTokens())
	require.Equal(t, bondAmount.ToDec(), validator.SelfDelegation)

	// delegate tokens to the validator
	msgDelegate := NewTestMsgDelegate(delAddr, valAddr, bondAmount)
	_, err = handleMsgDelegate(ctx, msgDelegate, keeper)
	require.Nil(t, err, "expected delegation to be ok, got %v", err)

	// verify validator bonded shares
	validator, found = keeper.GetValidator(ctx, valAddr)
	require.True(t, found)
	require.Equal(t, bondAmount.MulRaw(2), validator.DelegatorShares.RoundInt())
	require.Equal(t, bondAmount.MulRaw(2), validator.BondedTokens())
	// verify self delegation
	require.Equal(t, bondAmount.ToDec(), validator.SelfDelegation)

	// unbond validator total self-delegations (which should jail the validator)
	unbondAmt := sdk.NewCoin(sdk.DefaultBondDenom, bondAmount)
	msgUndelegate := NewMsgUndelegate(sdk.AccAddress(valAddr), valAddr, unbondAmt)

	got, err := handleMsgUndelegate(ctx, msgUndelegate, keeper)
	require.Nil(t, err, "expected begin unbonding validator msg to be ok, got %v", err)

	var finishTime time.Time
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(got.Data, &finishTime)
	ctx = ctx.WithBlockTime(finishTime)
	EndBlocker(ctx, keeper)

	// verify the validator record still exists, is jailed, and has correct tokens
	validator, found = keeper.GetValidator(ctx, valAddr)
	require.True(t, found)
	require.True(t, validator.Jailed)
	require.Equal(t, bondAmount, validator.Tokens)
	require.Equal(t, sdk.ZeroDec(), validator.SelfDelegation)

	// verify delegation still exists
	bond, found := keeper.GetDelegation(ctx, delAddr, valAddr)
	require.True(t, found)
	require.Equal(t, bondAmount, bond.Shares.RoundInt())
	require.Equal(t, bondAmount, validator.DelegatorShares.RoundInt())

	// verify the validator can still self-delegate
	msgSelfDelegate := NewTestMsgDelegate(sdk.AccAddress(valAddr), valAddr, bondAmount)
	_, err = handleMsgDelegate(ctx, msgSelfDelegate, keeper)
	require.Nil(t, err, "expected delegation to be ok, got %v", err)

	// verify validator bonded shares
	validator, found = keeper.GetValidator(ctx, valAddr)
	require.True(t, found)
	require.Equal(t, bondAmount.MulRaw(2), validator.DelegatorShares.RoundInt())
	require.Equal(t, bondAmount.MulRaw(2), validator.Tokens)
	require.Equal(t, bondAmount.ToDec(), validator.SelfDelegation)

	// unjail the validator now that is has non-zero self-delegated shares
	keeper.Unjail(ctx, valConsAddr)

	// verify the validator can now accept delegations
	msgDelegate = NewTestMsgDelegate(delAddr, valAddr, bondAmount)
	_, err = handleMsgDelegate(ctx, msgDelegate, keeper)
	require.Nil(t, err, "expected delegation to be ok, got %v", err)

	// verify validator bonded shares
	validator, found = keeper.GetValidator(ctx, valAddr)
	require.True(t, found)
	require.Equal(t, bondAmount.MulRaw(3), validator.DelegatorShares.RoundInt())
	require.Equal(t, bondAmount.MulRaw(3), validator.Tokens)
	require.Equal(t, bondAmount.ToDec(), validator.SelfDelegation)

	// verify new delegation
	bond, found = keeper.GetDelegation(ctx, delAddr, valAddr)
	require.True(t, found)
	require.Equal(t, bondAmount.MulRaw(2), bond.Shares.RoundInt())
	require.Equal(t, bondAmount.MulRaw(3), validator.DelegatorShares.RoundInt())
}

func TestIncrementsMsgDelegate(t *testing.T) {
	initPower := int64(1000)
	initBond := sdk.TokensFromConsensusPower(initPower)
	ctx, accMapper, keeper, _ := keep.CreateTestInput(t, false, initPower)
	params := keeper.GetParams(ctx)

	bondAmount := sdk.TokensFromConsensusPower(10)
	validatorAddr, delegatorAddr := sdk.ValAddress(keep.Addrs[0]), keep.Addrs[1]

	// first create validator
	msgCreateValidator := NewTestMsgCreateValidator(validatorAddr, keep.PKs[0], bondAmount)
	_, err := handleMsgCreateValidator(ctx, msgCreateValidator, keeper)
	require.Nil(t, err, "expected create validator msg to be ok, got %v", err)

	// apply TM updates
	keeper.ApplyAndReturnValidatorSetUpdates(ctx)

	validator, found := keeper.GetValidator(ctx, validatorAddr)
	require.True(t, found)
	require.Equal(t, sdk.Bonded, validator.Status)
	require.Equal(t, bondAmount, validator.DelegatorShares.RoundInt())
	require.Equal(t, bondAmount, validator.BondedTokens(), "validator: %v", validator)
	require.Equal(t, bondAmount.ToDec(), validator.SelfDelegation)

	_, found = keeper.GetDelegation(ctx, delegatorAddr, validatorAddr)
	require.False(t, found)

	bond, found := keeper.GetDelegation(ctx, sdk.AccAddress(validatorAddr), validatorAddr)
	require.True(t, found)
	require.Equal(t, bondAmount, bond.Shares.RoundInt())

	bondedTokens := keeper.TotalBondedTokens(ctx)
	require.Equal(t, bondAmount.Int64(), bondedTokens.Int64())

	// just send the same msgbond multiple times
	msgDelegate := NewTestMsgDelegate(delegatorAddr, validatorAddr, bondAmount)

	for i := int64(0); i < 5; i++ {
		ctx = ctx.WithBlockHeight(int64(i))

		_, err := handleMsgDelegate(ctx, msgDelegate, keeper)
		require.Nil(t, err, "expected msg %d to be ok, got %v", i, err)

		//Check that the accounts and the bond account have the appropriate values
		validator, found := keeper.GetValidator(ctx, validatorAddr)
		require.True(t, found)
		bond, found := keeper.GetDelegation(ctx, delegatorAddr, validatorAddr)
		require.True(t, found)

		expBond := bondAmount.MulRaw(i + 1)
		expDelegatorShares := bondAmount.MulRaw(i + 2) // (1 self delegation)
		expDelegatorAcc := initBond.Sub(expBond)

		gotBond := bond.Shares.RoundInt()
		gotDelegatorShares := validator.DelegatorShares.RoundInt()
		gotDelegatorAcc := accMapper.GetAccount(ctx, delegatorAddr).GetCoins().AmountOf(params.BondDenom)

		// verify self delegation
		require.Equal(t, bondAmount.ToDec(), validator.SelfDelegation)

		require.Equal(t, expBond, gotBond,
			"i: %v\nexpBond: %v\ngotBond: %v\nvalidator: %v\nbond: %v\n",
			i, expBond, gotBond, validator, bond)
		require.Equal(t, expDelegatorShares, gotDelegatorShares,
			"i: %v\nexpDelegatorShares: %v\ngotDelegatorShares: %v\nvalidator: %v\nbond: %v\n",
			i, expDelegatorShares, gotDelegatorShares, validator, bond)
		require.Equal(t, expDelegatorAcc, gotDelegatorAcc,
			"i: %v\nexpDelegatorAcc: %v\ngotDelegatorAcc: %v\nvalidator: %v\nbond: %v\n",
			i, expDelegatorAcc, gotDelegatorAcc, validator, bond)
	}
}

func TestEditValidatorDecreaseMinSelfDelegation(t *testing.T) {
	validatorAddr := sdk.ValAddress(keep.Addrs[0])

	initPower := int64(100)
	initBond := sdk.TokensFromConsensusPower(100)
	ctx, _, keeper, _ := keep.CreateTestInput(t, false, initPower)
	_ = setInstantUnbondPeriod(keeper, ctx)

	// create validator
	msgCreateValidator := NewTestMsgCreateValidator(validatorAddr, keep.PKs[0], initBond)
	msgCreateValidator.MinSelfDelegation = sdk.NewInt(2)
	_, err := handleMsgCreateValidator(ctx, msgCreateValidator, keeper)
	require.Nil(t, err, "expected create-validator to be ok, got %v", err)

	// must end-block
	updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.Equal(t, 1, len(updates))

	// verify the self-delegation exists
	bond, found := keeper.GetDelegation(ctx, sdk.AccAddress(validatorAddr), validatorAddr)
	require.True(t, found)
	gotBond := bond.Shares.RoundInt()
	require.Equal(t, initBond, gotBond,
		"initBond: %v\ngotBond: %v\nbond: %v\n",
		initBond, gotBond, bond)

	newMinSelfDelegation := sdk.OneInt()
	msgEditValidator := NewMsgEditValidator(validatorAddr, Description{}, nil, &newMinSelfDelegation)
	_, err = handleMsgEditValidator(ctx, msgEditValidator, keeper)
	require.NotNil(t, err, "should not be able to decrease minSelfDelegation")
	require.Equal(t, err, ErrMinSelfDelegationDecreased)
}

func TestEditValidatorIncreaseMinSelfDelegationBeyondCurrentBond(t *testing.T) {
	validatorAddr := sdk.ValAddress(keep.Addrs[0])

	initPower := int64(100)
	initBond := sdk.TokensFromConsensusPower(100)
	ctx, _, keeper, _ := keep.CreateTestInput(t, false, initPower)
	_ = setInstantUnbondPeriod(keeper, ctx)

	// create validator
	msgCreateValidator := NewTestMsgCreateValidator(validatorAddr, keep.PKs[0], initBond)
	msgCreateValidator.MinSelfDelegation = sdk.NewInt(2)
	_, err := handleMsgCreateValidator(ctx, msgCreateValidator, keeper)
	require.Nil(t, err, "expected create-validator to be ok, got %v", err)

	// must end-block
	updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.Equal(t, 1, len(updates))

	// verify the self-delegation exists
	bond, found := keeper.GetDelegation(ctx, sdk.AccAddress(validatorAddr), validatorAddr)
	require.True(t, found)
	gotBond := bond.Shares.RoundInt()
	require.Equal(t, initBond, gotBond,
		"initBond: %v\ngotBond: %v\nbond: %v\n",
		initBond, gotBond, bond)

	newMinSelfDelegation := initBond.Add(sdk.OneInt())
	msgEditValidator := NewMsgEditValidator(validatorAddr, Description{}, nil, &newMinSelfDelegation)
	_, err = handleMsgEditValidator(ctx, msgEditValidator, keeper)
	require.NotNil(t, err, "should not be able to increase minSelfDelegation above current self delegation")
	require.Equal(t, err, ErrSelfDelegationBelowMinimum)
}

func TestIncrementsMsgUnbond(t *testing.T) {
	initPower := int64(1000)
	initBond := sdk.TokensFromConsensusPower(initPower)
	ctx, accMapper, keeper, _ := keep.CreateTestInput(t, false, initPower)
	params := setInstantUnbondPeriod(keeper, ctx)
	denom := params.BondDenom

	// create validator, delegate
	validatorAddr, delegatorAddr := sdk.ValAddress(keep.Addrs[0]), keep.Addrs[1]

	msgCreateValidator := NewTestMsgCreateValidator(validatorAddr, keep.PKs[0], initBond)
	_, err := handleMsgCreateValidator(ctx, msgCreateValidator, keeper)
	require.Nil(t, err, "expected create-validator to be ok, got %v", err)

	// initial balance
	amt1 := accMapper.GetAccount(ctx, delegatorAddr).GetCoins().AmountOf(denom)

	msgDelegate := NewTestMsgDelegate(delegatorAddr, validatorAddr, initBond)
	_, err = handleMsgDelegate(ctx, msgDelegate, keeper)
	require.Nil(t, err, "expected delegation to be ok, got %v", err)

	// balance should have been subtracted after delegation
	amt2 := accMapper.GetAccount(ctx, delegatorAddr).GetCoins().AmountOf(denom)
	require.True(sdk.IntEq(t, amt1.Sub(initBond), amt2))

	// apply TM updates
	keeper.ApplyAndReturnValidatorSetUpdates(ctx)

	validator, found := keeper.GetValidator(ctx, validatorAddr)
	require.True(t, found)
	require.Equal(t, initBond.MulRaw(2), validator.DelegatorShares.RoundInt())
	require.Equal(t, initBond.MulRaw(2), validator.BondedTokens())

	// just send the same msgUnbond multiple times
	// TODO use decimals here
	unbondAmt := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))
	msgUndelegate := NewMsgUndelegate(delegatorAddr, validatorAddr, unbondAmt)
	numUnbonds := int64(5)
	for i := int64(0); i < numUnbonds; i++ {

		got, err := handleMsgUndelegate(ctx, msgUndelegate, keeper)
		require.Nil(t, err, "expected msg %d to be ok, got %v", i, err)
		var finishTime time.Time
		types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(got.Data, &finishTime)
		ctx = ctx.WithBlockTime(finishTime)
		EndBlocker(ctx, keeper)

		// check that the accounts and the bond account have the appropriate values
		validator, found = keeper.GetValidator(ctx, validatorAddr)
		require.True(t, found)
		bond, found := keeper.GetDelegation(ctx, delegatorAddr, validatorAddr)
		require.True(t, found)

		expBond := initBond.Sub(unbondAmt.Amount.Mul(sdk.NewInt(i + 1)))
		expDelegatorShares := initBond.MulRaw(2).Sub(unbondAmt.Amount.Mul(sdk.NewInt(i + 1)))
		expDelegatorAcc := initBond.Sub(expBond)

		gotBond := bond.Shares.RoundInt()
		gotDelegatorShares := validator.DelegatorShares.RoundInt()
		gotDelegatorAcc := accMapper.GetAccount(ctx, delegatorAddr).GetCoins().AmountOf(params.BondDenom)

		require.Equal(t, expBond.Int64(), gotBond.Int64(),
			"i: %v\nexpBond: %v\ngotBond: %v\nvalidator: %v\nbond: %v\n",
			i, expBond, gotBond, validator, bond)
		require.Equal(t, expDelegatorShares.Int64(), gotDelegatorShares.Int64(),
			"i: %v\nexpDelegatorShares: %v\ngotDelegatorShares: %v\nvalidator: %v\nbond: %v\n",
			i, expDelegatorShares, gotDelegatorShares, validator, bond)
		require.Equal(t, expDelegatorAcc.Int64(), gotDelegatorAcc.Int64(),
			"i: %v\nexpDelegatorAcc: %v\ngotDelegatorAcc: %v\nvalidator: %v\nbond: %v\n",
			i, expDelegatorAcc, gotDelegatorAcc, validator, bond)
	}

	// these are more than we have bonded now
	errorCases := []sdk.Int{
		//1<<64 - 1, // more than int64 power
		//1<<63 + 1, // more than int64 power
		sdk.TokensFromConsensusPower(1<<63 - 1),
		sdk.TokensFromConsensusPower(1 << 31),
		initBond,
	}

	for i, c := range errorCases {
		unbondAmt := sdk.NewCoin(sdk.DefaultBondDenom, c)
		msgUndelegate := NewMsgUndelegate(delegatorAddr, validatorAddr, unbondAmt)
		_, err = handleMsgUndelegate(ctx, msgUndelegate, keeper)
		require.NotNil(t, err, "expected unbond msg to fail, index: %v", i)
		require.Equal(t, err, ErrNotEnoughDelegationShares)
	}

	leftBonded := initBond.Sub(unbondAmt.Amount.Mul(sdk.NewInt(numUnbonds)))

	// should be able to unbond remaining
	unbondAmt = sdk.NewCoin(sdk.DefaultBondDenom, leftBonded)
	msgUndelegate = NewMsgUndelegate(delegatorAddr, validatorAddr, unbondAmt)
	got, err := handleMsgUndelegate(ctx, msgUndelegate, keeper)
	require.Nil(t, err,
		"got: %v\nmsgUnbond: %v\nshares: %s\nleftBonded: %s\n", got.Log, msgUndelegate, unbondAmt, leftBonded)
}

func TestMultipleMsgCreateValidator(t *testing.T) {
	initPower := int64(1000)
	initTokens := sdk.TokensFromConsensusPower(initPower)
	ctx, accMapper, keeper, _ := keep.CreateTestInput(t, false, initPower)
	params := setInstantUnbondPeriod(keeper, ctx)

	validatorAddrs := []sdk.ValAddress{
		sdk.ValAddress(keep.Addrs[0]),
		sdk.ValAddress(keep.Addrs[1]),
		sdk.ValAddress(keep.Addrs[2]),
	}
	delegatorAddrs := []sdk.AccAddress{
		keep.Addrs[0],
		keep.Addrs[1],
		keep.Addrs[2],
	}

	// bond them all
	for i, validatorAddr := range validatorAddrs {
		valTokens := sdk.TokensFromConsensusPower(10)
		msgCreateValidatorOnBehalfOf := NewTestMsgCreateValidator(validatorAddr, keep.PKs[i], valTokens)

		_, err := handleMsgCreateValidator(ctx, msgCreateValidatorOnBehalfOf, keeper)
		require.Nil(t, err, "expected msg %d to be ok, got %v", i, err)

		// verify that the account is bonded
		validators := keeper.GetValidators(ctx, 100)
		require.Equal(t, (i + 1), len(validators))

		val := validators[i]
		balanceExpd := initTokens.Sub(valTokens)
		balanceGot := accMapper.GetAccount(ctx, delegatorAddrs[i]).GetCoins().AmountOf(params.BondDenom)

		require.Equal(t, i+1, len(validators), "expected %d validators got %d, validators: %v", i+1, len(validators), validators)
		require.Equal(t, valTokens, val.DelegatorShares.RoundInt(), "expected %d shares, got %d", 10, val.DelegatorShares)
		require.Equal(t, balanceExpd, balanceGot, "expected account to have %d, got %d", balanceExpd, balanceGot)
	}

	// unbond them all by removing delegation
	for i, validatorAddr := range validatorAddrs {
		_, found := keeper.GetValidator(ctx, validatorAddr)
		require.True(t, found)

		unbondAmt := sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(10))
		msgUndelegate := NewMsgUndelegate(delegatorAddrs[i], validatorAddr, unbondAmt) // remove delegation
		got, err := handleMsgUndelegate(ctx, msgUndelegate, keeper)

		require.Nil(t, err, "expected msg %d to be ok, got %v", i, err)
		var finishTime time.Time

		// Jump to finishTime for unbonding period and remove from unbonding queue
		types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(got.Data, &finishTime)
		ctx = ctx.WithBlockTime(finishTime)

		EndBlocker(ctx, keeper)

		// Check that the validator is deleted from state
		validators := keeper.GetValidators(ctx, 100)
		require.Equal(t, len(validatorAddrs)-(i+1), len(validators),
			"expected %d validators got %d", len(validatorAddrs)-(i+1), len(validators))

		_, found = keeper.GetValidator(ctx, validatorAddr)
		require.False(t, found)

		gotBalance := accMapper.GetAccount(ctx, delegatorAddrs[i]).GetCoins().AmountOf(params.BondDenom)
		require.Equal(t, initTokens, gotBalance, "expected account to have %d, got %d", initTokens, gotBalance)
	}
}

func TestMultipleMsgDelegate(t *testing.T) {
	ctx, _, keeper, _ := keep.CreateTestInput(t, false, 1000)
	validatorAddr, delegatorAddrs := sdk.ValAddress(keep.Addrs[0]), keep.Addrs[1:]
	_ = setInstantUnbondPeriod(keeper, ctx)

	// first make a validator
	bondAmount := sdk.TokensFromConsensusPower(1000)
	msgCreateValidator := NewTestMsgCreateValidator(validatorAddr, keep.PKs[0], bondAmount)
	_, err := handleMsgCreateValidator(ctx, msgCreateValidator, keeper)
	require.Nil(t, err, "expected msg to be ok, got %v", err)

	// apply TM updates
	keeper.ApplyAndReturnValidatorSetUpdates(ctx)

	validator, found := keeper.GetValidator(ctx, validatorAddr)
	require.True(t, found)
	require.Equal(t, sdk.Bonded, validator.Status)
	require.Equal(t, bondAmount, validator.DelegatorShares.RoundInt())
	require.Equal(t, bondAmount, validator.BondedTokens(), "validator: %v", validator)
	require.Equal(t, bondAmount.ToDec(), validator.SelfDelegation)
	require.Equal(t, sdk.OneDec(), validator.BondedLever(true, sdk.ZeroDec()))

	// delegate multiple parties
	for i, delegatorAddr := range delegatorAddrs {
		msgDelegate := NewTestMsgDelegate(delegatorAddr, validatorAddr, sdk.NewInt(10))
		_, err := handleMsgDelegate(ctx, msgDelegate, keeper)
		require.Nil(t, err, "expected msg %d to be ok, got %v", i, err)

		// check that the account is bonded
		bond, found := keeper.GetDelegation(ctx, delegatorAddr, validatorAddr)
		require.True(t, found)
		require.NotNil(t, bond, "expected delegatee bond %d to exist", bond)
	}

	// unbond them all
	for i, delegatorAddr := range delegatorAddrs {
		unbondAmt := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))
		msgUndelegate := NewMsgUndelegate(delegatorAddr, validatorAddr, unbondAmt)

		got, err := handleMsgUndelegate(ctx, msgUndelegate, keeper)
		require.Nil(t, err, "expected msg %d to be ok, got %v", i, err)

		var finishTime time.Time
		types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(got.Data, &finishTime)

		ctx = ctx.WithBlockTime(finishTime)
		EndBlocker(ctx, keeper)

		// check that the account is unbonded
		_, found := keeper.GetDelegation(ctx, delegatorAddr, validatorAddr)
		require.False(t, found)
	}
}

func TestJailValidator(t *testing.T) {
	ctx, _, keeper, _ := keep.CreateTestInput(t, false, 1000)
	validatorAddr, delegatorAddr := sdk.ValAddress(keep.Addrs[0]), keep.Addrs[1]
	_ = setInstantUnbondPeriod(keeper, ctx)
	bondAmount := sdk.TokensFromConsensusPower(10)

	// create the validator
	msgCreateValidator := NewTestMsgCreateValidator(validatorAddr, keep.PKs[0], bondAmount)
	_, err := handleMsgCreateValidator(ctx, msgCreateValidator, keeper)
	require.Nil(t, err, "expected no error on runMsgCreateValidator")

	// apply TM updates
	updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.Equal(t, 1, len(updates))

	// bond a delegator
	msgDelegate := NewTestMsgDelegate(delegatorAddr, validatorAddr, bondAmount)
	_, err = handleMsgDelegate(ctx, msgDelegate, keeper)
	require.Nil(t, err, "expected ok, got %v", err)

	// unbond the validators bond portion
	unbondAmt := sdk.NewCoin(sdk.DefaultBondDenom, bondAmount)
	msgUndelegateValidator := NewMsgUndelegate(sdk.AccAddress(validatorAddr), validatorAddr, unbondAmt)
	got, err := handleMsgUndelegate(ctx, msgUndelegateValidator, keeper)
	require.Nil(t, err, "expected no error: %v", err)

	var finishTime time.Time
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(got.Data, &finishTime)

	ctx = ctx.WithBlockTime(finishTime)
	EndBlocker(ctx, keeper)

	validator, found := keeper.GetValidator(ctx, validatorAddr)
	require.True(t, found)
	require.True(t, validator.Jailed, "%v", validator)

	// test that the delegator can still withdraw their bonds
	msgUndelegateDelegator := NewMsgUndelegate(delegatorAddr, validatorAddr, unbondAmt)

	got, err = handleMsgUndelegate(ctx, msgUndelegateDelegator, keeper)
	require.Nil(t, err, "expected no error")
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(got.Data, &finishTime)

	ctx = ctx.WithBlockTime(finishTime)
	EndBlocker(ctx, keeper)

	// verify that the pubkey can now be reused
	_, err = handleMsgCreateValidator(ctx, msgCreateValidator, keeper)
	require.Nil(t, err, "expected ok, got %v", err)
}

func TestValidatorBondedLever(t *testing.T) {
	initPower := int64(1000)
	ctx, _, keeper, _ := keep.CreateTestInput(t, false, initPower)

	bondAmount := sdk.TokensFromConsensusPower(10)
	validatorAddr, delegatorAddr := sdk.ValAddress(keep.Addrs[0]), keep.Addrs[1]

	// first create validator
	msgCreateValidator := NewTestMsgCreateValidator(validatorAddr, keep.PKs[0], bondAmount)
	_, err := handleMsgCreateValidator(ctx, msgCreateValidator, keeper)
	require.Nil(t, err, "expected create validator msg to be ok, got %v", err)

	// apply TM updates
	keeper.ApplyAndReturnValidatorSetUpdates(ctx)

	validator, found := keeper.GetValidator(ctx, validatorAddr)
	require.True(t, found)
	require.Equal(t, sdk.Bonded, validator.Status)
	require.Equal(t, bondAmount, validator.DelegatorShares.RoundInt())
	require.Equal(t, bondAmount, validator.BondedTokens(), "validator: %v", validator)
	require.Equal(t, bondAmount.ToDec(), validator.SelfDelegation)
	require.Equal(t, sdk.OneDec(), validator.BondedLever(true, sdk.ZeroDec()))

	bondedTokens := keeper.TotalBondedTokens(ctx)
	require.Equal(t, bondAmount.Int64(), bondedTokens.Int64())

	// just send the same msgbond multiple times
	msgDelegate := NewTestMsgDelegate(delegatorAddr, validatorAddr, bondAmount.MulRaw(3))
	_, err = handleMsgDelegate(ctx, msgDelegate, keeper)
	require.Nil(t, err, "expected msg to be ok, got %v", err)

	//Check that the accounts and the bond account have the appropriate values
	validator, found = keeper.GetValidator(ctx, validatorAddr)
	require.True(t, found)
	require.Equal(t, sdk.Bonded, validator.Status)
	require.Equal(t, bondAmount.MulRaw(4), validator.DelegatorShares.RoundInt())
	require.Equal(t, bondAmount.MulRaw(4), validator.BondedTokens(), "validator: %v", validator)
	require.Equal(t, bondAmount.ToDec(), validator.SelfDelegation)
	require.Equal(t, sdk.OneDec().MulInt64(4), validator.BondedLever(true, sdk.ZeroDec()))

	// unbond self-delegation
	unbondAmt := sdk.NewCoin(sdk.DefaultBondDenom, bondAmount)
	msgUndelegate := NewMsgUndelegate(sdk.AccAddress(validatorAddr), validatorAddr, unbondAmt)
	got, err := handleMsgUndelegate(ctx, msgUndelegate, keeper)
	require.Nil(t, err, "expected msg to be ok, got %v", err)

	validator, found = keeper.GetValidator(ctx, validatorAddr)
	require.True(t, found)
	require.Equal(t, bondAmount.MulRaw(3), validator.DelegatorShares.RoundInt())
	require.Equal(t, bondAmount.MulRaw(3), validator.BondedTokens(), "validator: %v", validator)
	require.Equal(t, sdk.ZeroDec(), validator.SelfDelegation)
	require.Equal(t, sdk.NewDec(int64(^uint32(0))), validator.BondedLever(true, sdk.ZeroDec()))

	var finishTime time.Time
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(got.Data, &finishTime)
	ctx = ctx.WithBlockTime(finishTime)
	EndBlocker(ctx, keeper)
	EndBlocker(ctx, keeper)

	validator, found = keeper.GetValidator(ctx, validatorAddr)
	require.True(t, found)
	require.Equal(t, bondAmount.MulRaw(3), validator.DelegatorShares.RoundInt())
	require.Equal(t, sdk.ZeroInt(), validator.BondedTokens(), "validator: %v", validator)
	require.Equal(t, sdk.ZeroDec(), validator.SelfDelegation)
	require.Equal(t, sdk.NewDec(int64(^uint32(0))), validator.BondedLever(true, sdk.ZeroDec()))
}

func TestMaxValidators(t *testing.T) {
	ctx, _, keeper, _ := keep.CreateTestInput(t, false, int64(0))
	params := keeper.GetParams(ctx)

	maxValidators := params.MaxValidators
	maxValidatorsLimit := params.MaxValidatorsExtendingLimit
	delta := params.MaxValidatorsExtendingSpeed
	nextExtendingTime := params.NextExtendingTime
	for i := 0; maxValidators < maxValidatorsLimit; i++ {
		// end blocker
		ctx = ctx.WithBlockTime(nextExtendingTime)
		EndBlocker(ctx, keeper)

		// end blocker
		ctx = ctx.WithBlockTime(nextExtendingTime.Add(time.Hour * 1))
		EndBlocker(ctx, keeper)

		params = keeper.GetParams(ctx)
		// verify max validators
		require.Equal(t, maxValidators+delta, params.MaxValidators)

		maxValidators = params.MaxValidators
		maxValidatorsLimit = params.MaxValidatorsExtendingLimit
		delta = params.MaxValidatorsExtendingSpeed
		nextExtendingTime = params.NextExtendingTime
		//fmt.Println(fmt.Sprintf("round %v, nextExtendingtime: %v, maxValidators: %v, delta: %v, maxValidatorsLimit: %v",
		//	i, nextExtendingTime, maxValidators, delta, maxValidatorsLimit))
	}

	// end blocker
	ctx = ctx.WithBlockTime(nextExtendingTime)
	EndBlocker(ctx, keeper)

	params = keeper.GetParams(ctx)
	// verify max validators with upper limit
	require.Equal(t, params.MaxValidatorsExtendingLimit, params.MaxValidators)
}
