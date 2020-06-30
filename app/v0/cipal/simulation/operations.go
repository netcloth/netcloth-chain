package simulation

import (
	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"github.com/netcloth/netcloth-chain/app/v0/cipal/keeper"
	"github.com/netcloth/netcloth-chain/app/v0/cipal/types"
	"github.com/netcloth/netcloth-chain/app/v0/simulation"
	"github.com/netcloth/netcloth-chain/baseapp"
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/simapp/helpers"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdksimulation "github.com/netcloth/netcloth-chain/types/simulation"
	"math/rand"
)

func WeightedOperations(appParams sdksimulation.AppParams, cdc *codec.Codec, ak keeper.AccountKeeper, k keeper.Keeper) simulation.WeightedOperations {

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			100,
			SimulateMsgCreateCIpal(ak, k),
		),
	}
}

func SimulateMsgCreateCIpal(ak keeper.AccountKeeper, k keeper.Keeper) sdksimulation.Operation {

	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []sdksimulation.Account, chainID string) (sdksimulation.OperationMsg, []sdksimulation.FutureOperation, error) {
		acc, _ := sdksimulation.RandomAcc(r, accs)
		accountObj := ak.GetAccount(ctx, acc.Address)

		expiration := ctx.BlockHeader().Time.AddDate(0, 0, 1)
		adMsg := types.NewADParam(acc.Address.String(), acc.Address.String(), 1, expiration)
		sig, err := acc.PrivKey.Sign(adMsg.GetSignBytes())

		stdSig := auth.StdSignature{PubKey: acc.PubKey, Signature: sig}
		msg := types.NewMsgCIPALClaim(acc.Address, acc.Address.String(), acc.Address.String(), 1, expiration, stdSig)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			sdk.NewCoins(sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(1000000))),
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{accountObj.GetAccountNumber()},
			[]uint64{accountObj.GetSequence()},
			acc.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return sdksimulation.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return sdksimulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}
