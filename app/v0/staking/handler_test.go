package staking

import (
	"testing"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtypes "github.com/tendermint/tendermint/types"

	keep "github.com/netcloth/netcloth-chain/app/v0/staking/keeper"
	"github.com/netcloth/netcloth-chain/app/v0/staking/types"
	sdk "github.com/netcloth/netcloth-chain/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
