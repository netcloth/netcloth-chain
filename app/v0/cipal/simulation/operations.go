package simulation

//
//import (
//	"github.com/spf13/viper"
//	"math/rand"
//	"time"
//
//	"github.com/netcloth/netcloth-chain/app/v0/cipal/keeper"
//	"github.com/netcloth/netcloth-chain/app/v0/cipal/types"
//	"github.com/netcloth/netcloth-chain/app/v0/simulation"
//	"github.com/netcloth/netcloth-chain/baseapp"
//	"github.com/netcloth/netcloth-chain/codec"
//	"github.com/netcloth/netcloth-chain/simapp/helpers"
//	sdk "github.com/netcloth/netcloth-chain/types"
//	sdksimulation "github.com/netcloth/netcloth-chain/types/simulation"
//)
//
//func WeightedOperations(appParams sdksimulation.AppParams, cdc *codec.Codec, ak keeper.AccountKeeper, k keeper.Keeper) simulation.WeightedOperations {
//
//	return simulation.WeightedOperations{
//		simulation.NewWeightedOperation(
//			100,
//			SimulateMsgCreateCIpal(ak, k),
//		),
//	}
//}
//
//func SimulateMsgCreateCIpal(ak keeper.AccountKeeper, k keeper.Keeper) sdksimulation.Operation {
//
//	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []sdksimulation.Account, chainID string) (sdksimulation.OperationMsg, []sdksimulation.FutureOperation, error) {
//		acc, _ := sdksimulation.RandomAcc(r, accs)
//
//		accountObj := ak.GetAccount(ctx, acc.Address)
//		amount := accountObj.GetCoins().AmountOf(bondDenom)
//		if !amount.IsPositive() {
//			return sdksimulation.NoOpMsg(types.ModuleName, types.TypeMsgIPALNodeClaim, "balance is negative"), nil, nil
//		}
//
//		bondAmt := sdksimulation.RandomAmount(r, amount)
//		for bondAmt.LT(minBond.Amount) || bondAmt.GT(minBond.Amount.Mul(sdk.NewInt(2))) {
//			bondAmt = sdksimulation.RandomAmount(r, amount)
//		}
//
//		bond := sdk.NewCoin(bondDenom, bondAmt)
//
//		var fees sdk.Coins
//		var err error
//		fees, err = sdksimulation.RandomFees(r, ctx, sdk.NewCoins(bond))
//		if err != nil {
//			return sdksimulation.NoOpMsg(types.ModuleName, types.TypeMsgIPALNodeClaim, "unable to generate fees"), nil, err
//		}
//
//		moniker, website, details, extension := sdksimulation.RandStringOfLength(r, 1000), sdksimulation.RandStringOfLength(r, 1000), sdksimulation.RandStringOfLength(r, 1000), sdksimulation.RandStringOfLength(r, 1000)
//
//		endPoint := types.NewEndpoint(r.Uint64(), "192.168.1.100:666")
//
//		info, err := txBldr.Keybase().Get(cliCtxUser.GetFromName())
//		if err != nil {
//			return err
//		}
//		userAddress := info.GetAddress().String()
//
//		expiration := time.Now().UTC().AddDate(0, 0, 1)
//		adMsg := types.NewADParam(userAddress, serviceAddress, serviceType, expiration)
//
//		msg := types.NewMsgCIPALClaim()
//
//		tx := helpers.GenTx(
//			[]sdk.Msg{msg},
//			fees,
//			helpers.DefaultGenTxGas,
//			chainID,
//			[]uint64{accountObj.GetAccountNumber()},
//			[]uint64{accountObj.GetSequence()},
//			acc.PrivKey,
//		)
//
//		_, _, err = app.Deliver(tx)
//		if err != nil {
//			return sdksimulation.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
//		}
//
//		return sdksimulation.NewOperationMsg(msg, true, ""), nil, nil
//	}
//}
