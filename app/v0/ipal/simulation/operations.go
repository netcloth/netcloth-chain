package simulation

import (
	"math/rand"

	"github.com/netcloth/netcloth-chain/app/simapp/helpers"
	"github.com/netcloth/netcloth-chain/app/v0/ipal/keeper"
	"github.com/netcloth/netcloth-chain/app/v0/ipal/types"
	"github.com/netcloth/netcloth-chain/app/v0/simulation"
	"github.com/netcloth/netcloth-chain/baseapp"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
	simtypes "github.com/netcloth/netcloth-chain/types/simulation"
)

func WeightedOperations(appParams simtypes.AppParams, cdc *codec.Codec, ak keeper.AccountKeeper, k keeper.Keeper) simulation.WeightedOperations {

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			100,
			SimulateMsgCreateIpal(ak, k),
		),
	}
}

func SimulateMsgCreateIpal(ak keeper.AccountKeeper, k keeper.Keeper) simtypes.Operation {

	return func(r *rand.Rand, app interface{}, ctx sdk.Context, accs []simtypes.Account, chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		var a *baseapp.BaseApp
		var ok = false
		if a, ok = app.(*baseapp.BaseApp); !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgIPALNodeClaim, "app invalid"), nil, nil
		}

		acc, _ := simtypes.RandomAcc(r, accs)
		_, found := k.GetIPALNode(ctx, acc.Address)
		if found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgIPALNodeClaim, "ipal exist"), nil, nil
		}

		minBond := k.GetParams(ctx).MinBond
		bondDenom := minBond.Denom

		accountObj := ak.GetAccount(ctx, acc.Address)
		amount := accountObj.GetCoins().AmountOf(bondDenom)
		if !amount.IsPositive() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgIPALNodeClaim, "balance is negative"), nil, nil
		}

		bondAmt := simtypes.RandomAmount(r, amount)
		for bondAmt.LT(minBond.Amount) || bondAmt.GT(minBond.Amount.Mul(sdk.NewInt(2))) {
			bondAmt = simtypes.RandomAmount(r, amount)
		}

		bond := sdk.NewCoin(bondDenom, bondAmt)

		var fees sdk.Coins
		var err error
		fees, err = simtypes.RandomFees(r, ctx, sdk.NewCoins(bond))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgIPALNodeClaim, "unable to generate fees"), nil, err
		}

		moniker, website, details, extension := simtypes.RandStringOfLength(r, 1000), simtypes.RandStringOfLength(r, 1000), simtypes.RandStringOfLength(r, 1000), simtypes.RandStringOfLength(r, 1000)

		endPoint := types.NewEndpoint(r.Uint64(), "192.168.1.100:666")
		msg := types.NewMsgIPALNodeClaim(acc.Address, moniker, website, details, extension, types.Endpoints{endPoint}, bond)

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
