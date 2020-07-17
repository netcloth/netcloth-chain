package simulation

import (
	"math/rand"

	"github.com/netcloth/netcloth-chain/app/simapp/helpers"
	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"github.com/netcloth/netcloth-chain/app/v0/cipal/keeper"
	"github.com/netcloth/netcloth-chain/app/v0/cipal/types"
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
			SimulateMsgCreateCIpal(ak, k),
		),
	}
}

func SimulateMsgCreateCIpal(ak keeper.AccountKeeper, k keeper.Keeper) simtypes.Operation {

	return func(r *rand.Rand, app interface{}, ctx sdk.Context, accs []simtypes.Account, chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		a := baseapp.DereferenceBaseApp(app)
		if a == nil {
			return simtypes.NoOpMsg(types.ModuleName, "create_cipal", "app invalid"), nil, nil
		}

		acc, _ := simtypes.RandomAcc(r, accs)
		accountObj := ak.GetAccount(ctx, acc.Address)

		expiration := ctx.BlockHeader().Time.AddDate(0, 0, 1)
		adMsg := types.NewADParam(acc.Address.String(), acc.Address.String(), 1, expiration)
		sig, _ := acc.PrivKey.Sign(adMsg.GetSignBytes())

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

		_, _, err := a.Deliver(tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
