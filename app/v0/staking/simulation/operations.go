package simulation

import (
	"fmt"
	"math/rand"

	"github.com/netcloth/netcloth-chain/app/simapp/helpers"
	simappparams "github.com/netcloth/netcloth-chain/app/simapp/params"
	"github.com/netcloth/netcloth-chain/app/v0/simulation"
	"github.com/netcloth/netcloth-chain/app/v0/staking/keeper"
	"github.com/netcloth/netcloth-chain/app/v0/staking/types"
	"github.com/netcloth/netcloth-chain/baseapp"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
	simtypes "github.com/netcloth/netcloth-chain/types/simulation"
)

// Simulation operation weights constants
const (
	OpWeightMsgCreateValidator = "op_weight_msg_create_validator"
	OpWeightMsgEditValidator   = "op_weight_msg_edit_validator"
	OpWeightMsgDelegate        = "op_weight_msg_delegate"
	OpWeightMsgUndelegate      = "op_weight_msg_undelegate"
	OpWeightMsgBeginRedelegate = "op_weight_msg_begin_redelegate"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams, cdc *codec.Codec, ak types.AccountKeeper, k keeper.Keeper) simulation.WeightedOperations {
	var (
		weightMsgCreateValidator int
		weightMsgEditValidator   int
		weightMsgDelegate        int
		weightMsgUndelegate      int
		weightMsgBeginRedelegate int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgCreateValidator, &weightMsgCreateValidator, nil,
		func(_ *rand.Rand) {
			weightMsgCreateValidator = simappparams.DefaultWeightMsgCreateValidator
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgEditValidator, &weightMsgEditValidator, nil,
		func(_ *rand.Rand) {
			weightMsgEditValidator = simappparams.DefaultWeightMsgEditValidator
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgDelegate, &weightMsgDelegate, nil,
		func(_ *rand.Rand) {
			weightMsgDelegate = simappparams.DefaultWeightMsgDelegate
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgUndelegate, &weightMsgUndelegate, nil,
		func(_ *rand.Rand) {
			weightMsgUndelegate = simappparams.DefaultWeightMsgUndelegate
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgBeginRedelegate, &weightMsgBeginRedelegate, nil,
		func(_ *rand.Rand) {
			weightMsgBeginRedelegate = simappparams.DefaultWeightMsgBeginRedelegate
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreateValidator,
			SimulateMsgCreateValidator(ak, k),
		),
		simulation.NewWeightedOperation(
			weightMsgEditValidator,
			SimulateMsgEditValidator(ak, k),
		),
		simulation.NewWeightedOperation(
			weightMsgDelegate,
			SimulateMsgDelegate(ak, k),
		),
		simulation.NewWeightedOperation(
			weightMsgUndelegate,
			SimulateMsgUndelegate(ak, k),
		),
		simulation.NewWeightedOperation(
			weightMsgBeginRedelegate,
			SimulateMsgBeginRedelegate(ak, k),
		),
	}
}

// SimulateMsgCreateValidator generates a MsgCreateValidator with random values
// nolint: interfacer
func SimulateMsgCreateValidator(ak types.AccountKeeper, k keeper.Keeper) simtypes.Operation {

	return func(r *rand.Rand, app interface{}, ctx sdk.Context, accs []simtypes.Account, chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		a := baseapp.DereferenceBaseApp(app)
		if a == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateValidator, "app invalid"), nil, nil
		}

		acc, _ := simtypes.RandomAcc(r, accs)
		address := sdk.ValAddress(acc.Address)

		_, found := k.GetValidator(ctx, address)
		if found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateValidator, "validator exist"), nil, nil
		}

		bondDenom := k.GetParams(ctx).BondDenom

		accountObj := ak.GetAccount(ctx, acc.Address)
		amount := accountObj.GetCoins().AmountOf(bondDenom)
		if !amount.IsPositive() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateValidator, "balance is negative"), nil, nil
		}

		amount = amount.Quo(sdk.NewInt(2))

		selfDelegation := sdk.NewCoin(bondDenom, amount)

		var fees sdk.Coins
		var err error
		coins, hasNeg := accountObj.GetCoins().SafeSub(sdk.Coins{selfDelegation})
		if !hasNeg {
			fees, err = simtypes.RandomFees(r, ctx, coins)
			if err != nil {
				return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateValidator, "unable to generate fees"), nil, err
			}
		}

		description := types.NewDescription(
			simtypes.RandStringOfLength(r, 10),
			simtypes.RandStringOfLength(r, 10),
			simtypes.RandStringOfLength(r, 10),
			simtypes.RandStringOfLength(r, 10),
		)

		maxCommission := sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 0, 100)), 2)
		commission := types.NewCommissionRates(
			simtypes.RandomDecAmount(r, maxCommission),
			maxCommission,
			simtypes.RandomDecAmount(r, maxCommission),
		)

		msg := types.NewMsgCreateValidator(address, acc.PubKey, selfDelegation, description, commission, sdk.OneInt())

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{accountObj.GetAccountNumber()},
			[]uint64{accountObj.GetSequence()},
			acc.PrivKey,
		)

		_, _, err = a.Deliver(tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgEditValidator generates a MsgEditValidator with random values
// nolint: interfacer
func SimulateMsgEditValidator(ak types.AccountKeeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app interface{}, ctx sdk.Context, accs []simtypes.Account, chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		a := baseapp.DereferenceBaseApp(app)
		if a == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgEditValidator, "app invalid"), nil, nil
		}

		if len(k.GetAllValidators(ctx)) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgEditValidator, "number of validators equal zero"), nil, nil
		}

		val, ok := keeper.RandomValidator(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgEditValidator, "unable to pick a validator"), nil, nil
		}

		address := val.GetOperator()

		newCommissionRate := simtypes.RandomDecAmount(r, val.Commission.MaxRate)

		if err := val.Commission.ValidateNewRate(newCommissionRate, ctx.BlockHeader().Time); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgEditValidator, "invalid commission rate"), nil, nil
		}

		acc, found := simtypes.FindAccount(accs, sdk.AccAddress(address))
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgEditValidator, "unable to find account"), nil, fmt.Errorf("validator %s not found", val.GetOperator())
		}

		accountObj := ak.GetAccount(ctx, acc.Address)
		coins := accountObj.GetCoins()

		fees, err := simtypes.RandomFees(r, ctx, coins)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgEditValidator, "unable to generate fees"), nil, err
		}

		description := types.NewDescription(
			simtypes.RandStringOfLength(r, 10),
			simtypes.RandStringOfLength(r, 10),
			simtypes.RandStringOfLength(r, 10),
			simtypes.RandStringOfLength(r, 10),
		)

		msg := types.NewMsgEditValidator(address, description, &newCommissionRate, nil)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{accountObj.GetAccountNumber()},
			[]uint64{accountObj.GetSequence()},
			acc.PrivKey,
		)

		_, _, err = a.Deliver(tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgDelegate generates a MsgDelegate with random values
// nolint: interfacer
func SimulateMsgDelegate(ak types.AccountKeeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app interface{}, ctx sdk.Context, accs []simtypes.Account, chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		a := baseapp.DereferenceBaseApp(app)
		if a == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDelegate, "app invalid"), nil, nil
		}

		bondDenom := k.GetParams(ctx).BondDenom

		if len(k.GetAllValidators(ctx)) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDelegate, "number of validators equal zero"), nil, nil
		}

		val, ok := keeper.RandomValidator(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDelegate, "unable to pick a validator"), nil, nil
		}

		if val.InvalidExRate() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDelegate, "validator's invalid echange rate"), nil, nil
		}

		acc, _ := simtypes.RandomAcc(r, accs)
		accountObj := ak.GetAccount(ctx, acc.Address)
		coins := accountObj.GetCoins()
		if coins.IsAnyNegative() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDelegate, "balance is negative"), nil, nil
		}

		bondAmt := sdk.NewCoin(bondDenom, sdk.NewInt(1))

		msg := types.NewMsgDelegate(acc.Address, val.GetOperator(), bondAmt)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			sdk.NewCoins(sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(1000000))),
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{accountObj.GetAccountNumber()},
			[]uint64{accountObj.GetSequence()},
			acc.PrivKey,
		)

		_, _, err := a.Deliver(tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgUndelegate generates a MsgUndelegate with random values
// nolint: interfacer
func SimulateMsgUndelegate(ak types.AccountKeeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app interface{}, ctx sdk.Context, accs []simtypes.Account, chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		a := baseapp.DereferenceBaseApp(app)
		if a == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUndelegate, "app invalid"), nil, nil
		}

		validator, ok := keeper.RandomValidator(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUndelegate, "validator is not ok"), nil, nil
		}

		valAddr := validator.GetOperator()
		delegations := k.GetValidatorDelegations(ctx, validator.OperatorAddress)

		delegation := delegations[r.Intn(len(delegations))]
		delAddr := delegation.GetDelegatorAddr()

		if k.HasMaxUnbondingDelegationEntries(ctx, delAddr, valAddr) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUndelegate, "keeper does have a max unbonding delegation entries"), nil, nil
		}

		totalBond := validator.TokensFromShares(delegation.GetShares()).TruncateInt()
		if !totalBond.IsPositive() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUndelegate, "total bond is negative"), nil, nil
		}

		unbondAmt, err := simtypes.RandPositiveInt(r, totalBond)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUndelegate, "invalid unbond amount"), nil, err
		}

		if unbondAmt.IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUndelegate, "unbond amount is zero"), nil, nil
		}

		msg := types.NewMsgUndelegate(
			delAddr, valAddr, sdk.NewCoin(k.BondDenom(ctx), unbondAmt),
		)

		// need to retrieve the simulation account associated with delegation to retrieve PrivKey
		var simAccount simtypes.Account

		for _, simAcc := range accs {
			if simAcc.Address.Equals(delAddr) {
				simAccount = simAcc
				break
			}
		}
		// if simaccount.PrivKey == nil, delegation address does not exist in accs. Return error
		if simAccount.PrivKey == nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "account private key is nil"), nil, fmt.Errorf("delegation addr: %s does not exist in simulation accounts", delAddr)
		}

		accountObj := ak.GetAccount(ctx, delAddr)
		coins := accountObj.GetCoins()

		fees, err := simtypes.RandomFees(r, ctx, coins)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{accountObj.GetAccountNumber()},
			[]uint64{accountObj.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = a.Deliver(tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgBeginRedelegate generates a MsgBeginRedelegate with random values
// nolint: interfacer
func SimulateMsgBeginRedelegate(ak types.AccountKeeper, k keeper.Keeper) simtypes.Operation {

	return func(r *rand.Rand, app interface{}, ctx sdk.Context, accs []simtypes.Account, chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		a := baseapp.DereferenceBaseApp(app)
		if a == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBeginRedelegate, "app invalid"), nil, nil
		}

		srcVal, ok := keeper.RandomValidator(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBeginRedelegate, "unable to pick validator"), nil, nil
		}

		srcAddr := srcVal.GetOperator()
		delegations := k.GetValidatorDelegations(ctx, srcAddr)

		// get random delegator from src validator
		delegation := delegations[r.Intn(len(delegations))]
		delAddr := delegation.GetDelegatorAddr()

		if k.HasReceivingRedelegation(ctx, delAddr, srcAddr) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBeginRedelegate, "receveing redelegation is not allowed"), nil, nil // skip
		}

		// get random destination validator
		destVal, ok := keeper.RandomValidator(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBeginRedelegate, "unable to pick validator"), nil, nil
		}

		destAddr := destVal.GetOperator()
		if srcAddr.Equals(destAddr) || destVal.InvalidExRate() || k.HasMaxRedelegationEntries(ctx, delAddr, srcAddr, destAddr) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBeginRedelegate, "checks failed"), nil, nil
		}

		totalBond := srcVal.TokensFromShares(delegation.GetShares()).TruncateInt()
		if !totalBond.IsPositive() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBeginRedelegate, "total bond is negative"), nil, nil
		}

		destVal.SelfDelegation.TruncateInt()
		redAmt, err := simtypes.RandPositiveInt(r, sdk.NewInt(1))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBeginRedelegate, "unable to generate positive amount"), nil, err
		}

		if redAmt.IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBeginRedelegate, "amount is zero"), nil, nil
		}

		// check if the shares truncate to zero
		shares, err := srcVal.SharesFromTokens(redAmt)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBeginRedelegate, "invalid shares"), nil, err
		}

		if srcVal.TokensFromShares(shares).TruncateInt().IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBeginRedelegate, "shares truncate to zero"), nil, nil // skip
		}

		// need to retrieve the simulation account associated with delegation to retrieve PrivKey
		var simAccount simtypes.Account

		for _, simAcc := range accs {
			if simAcc.Address.Equals(delAddr) {
				simAccount = simAcc
				break
			}
		}

		// if simaccount.PrivKey == nil, delegation address does not exist in accs. Return error
		if simAccount.PrivKey == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBeginRedelegate, "account private key is nil"), nil, fmt.Errorf("delegation addr: %s does not exist in simulation accounts", delAddr)
		}

		accountObj := ak.GetAccount(ctx, delAddr)
		spendable := accountObj.GetCoins()

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBeginRedelegate, "unable to generate fees"), nil, err
		}

		msg := types.NewMsgBeginRedelegate(
			delAddr, srcAddr, destAddr,
			sdk.NewCoin(k.BondDenom(ctx), redAmt),
		)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{accountObj.GetAccountNumber()},
			[]uint64{accountObj.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = a.Deliver(tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
