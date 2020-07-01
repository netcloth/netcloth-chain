package slashing

import (
	"errors"
	"math/rand"

	"github.com/netcloth/netcloth-chain/app/simapp/helpers"
	simappparams "github.com/netcloth/netcloth-chain/app/simapp/params"
	"github.com/netcloth/netcloth-chain/app/v0/simulation"
	"github.com/netcloth/netcloth-chain/app/v0/slashing/types"
	stakingkeeper "github.com/netcloth/netcloth-chain/app/v0/staking/keeper"
	"github.com/netcloth/netcloth-chain/baseapp"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
	simtypes "github.com/netcloth/netcloth-chain/types/simulation"
)

// Simulation operation weights constants
const (
	OpWeightMsgUnjail = "op_weight_msg_unjail"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc *codec.Codec,
	ak AccountKeeper,
	k Keeper,
	sk stakingkeeper.Keeper) simulation.WeightedOperations {

	var weightMsgUnjail int
	appParams.GetOrGenerate(cdc, OpWeightMsgUnjail, &weightMsgUnjail, nil,
		func(_ *rand.Rand) {
			weightMsgUnjail = simappparams.DefaultWeightMsgUnjail
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgUnjail,
			SimulateMsgUnjail(ak, k, sk),
		),
	}
}

// SimulateMsgUnjail generates a MsgUnjail with random values
// nolint: interfacer
func SimulateMsgUnjail(
	ak AccountKeeper,
	k Keeper,
	sk stakingkeeper.Keeper) simtypes.Operation {

	return func(r *rand.Rand, app interface{}, ctx sdk.Context, accs []simtypes.Account, chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		var a *baseapp.BaseApp
		var ok = false
		if a, ok = app.(*baseapp.BaseApp); !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUnjail, "app invalid"), nil, nil
		}

		validator, ok := stakingkeeper.RandomValidator(r, sk, ctx)
		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUnjail, "validator is not ok"), nil, nil // skip
		}

		simAccount, found := simtypes.FindAccount(accs, sdk.AccAddress(validator.GetOperator()))
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUnjail, "unable to find account"), nil, nil // skip
		}

		if !validator.IsJailed() {
			// TODO: due to this condition this message is almost, if not always, skipped !
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUnjail, "validator is not jailed"), nil, nil
		}

		consAddr := sdk.ConsAddress(validator.GetConsPubKey().Address())
		info, found := k.GetValidatorSigningInfo(ctx, consAddr)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUnjail, "unable to find validator signing info"), nil, nil // skip
		}

		selfDel := sk.Delegation(ctx, simAccount.Address, validator.GetOperator())
		if selfDel == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUnjail, "self delegation is nil"), nil, nil // skip
		}

		account := ak.GetAccount(ctx, sdk.AccAddress(validator.GetOperator()))
		spendable := account.GetCoins()

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUnjail, "unable to generate fees"), nil, err
		}

		msg := types.NewMsgUnjail(validator.GetOperator())

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, res, err := a.Deliver(tx)

		// result should fail if:
		// - validator cannot be unjailed due to tombstone
		// - validator is still in jailed period
		// - self delegation too low
		if info.Tombstoned ||
			ctx.BlockHeader().Time.Before(info.JailedUntil) ||
			validator.TokensFromShares(selfDel.GetShares()).TruncateInt().LT(validator.GetMinSelfDelegation()) {
			if res != nil && err == nil {
				if info.Tombstoned {
					return simtypes.NewOperationMsg(msg, true, ""), nil, errors.New("validator should not have been unjailed if validator tombstoned")
				}
				if ctx.BlockHeader().Time.Before(info.JailedUntil) {
					return simtypes.NewOperationMsg(msg, true, ""), nil, errors.New("validator unjailed while validator still in jail period")
				}
				if validator.TokensFromShares(selfDel.GetShares()).TruncateInt().LT(validator.GetMinSelfDelegation()) {
					return simtypes.NewOperationMsg(msg, true, ""), nil, errors.New("validator unjailed even though self-delegation too low")
				}
			}
			// msg failed as expected
			return simtypes.NewOperationMsg(msg, false, ""), nil, nil
		}

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, errors.New(res.Log)
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
