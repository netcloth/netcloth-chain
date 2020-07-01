package simulation

import (
	"math/rand"

	"github.com/tendermint/tendermint/crypto"

	"github.com/netcloth/netcloth-chain/app/simapp/helpers"
	"github.com/netcloth/netcloth-chain/app/v0/bank/internal/keeper"
	"github.com/netcloth/netcloth-chain/app/v0/bank/internal/types"
	"github.com/netcloth/netcloth-chain/app/v0/simulation"
	"github.com/netcloth/netcloth-chain/baseapp"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
	simtypes "github.com/netcloth/netcloth-chain/types/simulation"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(appParams simtypes.AppParams, cdc *codec.Codec, ak types.AccountKeeper, bk keeper.Keeper) simulation.WeightedOperations {
	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			80,
			SimulateMsgSend(ak, bk),
		),
	}
}

// SimulateMsgSend tests and runs a single msg send where both
// accounts already exist.
// nolint: funlen
func SimulateMsgSend(ak types.AccountKeeper, bk keeper.Keeper) simtypes.Operation {

	return func(r *rand.Rand, app interface{}, ctx sdk.Context, accs []simtypes.Account, chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var a *baseapp.BaseApp
		var ok = false
		if a, ok = app.(*baseapp.BaseApp); !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSend, "app invalid"), nil, nil
		}

		if !bk.GetSendEnabled(ctx) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSend, "transfers are not enabled"), nil, nil
		}

		simAccount, toSimAcc, coins, skip := randomSendFields(r, ctx, accs, bk, ak)

		if skip {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSend, "skip all transfers"), nil, nil
		}

		msg := types.NewMsgSend(simAccount.Address, toSimAcc.Address, coins)

		err := sendMsgSend(r, a, bk, ak, msg, ctx, chainID, []crypto.PrivKey{simAccount.PrivKey})
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "invalid transfers"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func sendMsgSend(
	r *rand.Rand,
	app *baseapp.BaseApp,
	bk keeper.Keeper,
	ak types.AccountKeeper,
	msg types.MsgSend,
	ctx sdk.Context,
	chainID string,
	privkeys []crypto.PrivKey,
) error {

	var (
		fees sdk.Coins
		err  error
	)

	account := ak.GetAccount(ctx, msg.FromAddress)
	spendable := bk.GetCoins(ctx, account.GetAddress())

	coins, hasNeg := spendable.SafeSub(msg.Amount)
	if !hasNeg {
		fees, err = simtypes.RandomFees(r, ctx, coins)
		if err != nil {
			return err
		}
	}

	tx := helpers.GenTx(
		[]sdk.Msg{msg},
		fees,
		helpers.DefaultGenTxGas,
		chainID,
		[]uint64{account.GetAccountNumber()},
		[]uint64{account.GetSequence()},
		privkeys...,
	)

	_, _, err = app.Deliver(tx)
	if err != nil {
		return err
	}

	return nil
}

func randomSendFields(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account, bk keeper.Keeper, ak types.AccountKeeper) (simtypes.Account, simtypes.Account, sdk.Coins, bool) {

	simAccount, _ := simtypes.RandomAcc(r, accs)
	toSimAcc, _ := simtypes.RandomAcc(r, accs)

	// disallow sending money to yourself
	for simAccount.PubKey.Equals(toSimAcc.PubKey) {
		toSimAcc, _ = simtypes.RandomAcc(r, accs)
	}

	acc := ak.GetAccount(ctx, simAccount.Address)
	if acc == nil {
		return simAccount, toSimAcc, nil, true
	}

	spendable := bk.GetCoins(ctx, acc.GetAddress())

	sendCoins := simtypes.RandSubsetCoins(r, spendable)
	if sendCoins.Empty() {
		return simAccount, toSimAcc, nil, true
	}

	return simAccount, toSimAcc, sendCoins, false
}
