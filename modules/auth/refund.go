package auth

import (
	"github.com/netcloth/netcloth-chain/codec"
	auth "github.com/netcloth/netcloth-chain/modules/auth/types"
	"github.com/netcloth/netcloth-chain/types"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

type RefundKeeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

func NewRefundKeeper(cdc *codec.Codec, key sdk.StoreKey) RefundKeeper {
	return RefundKeeper{
		storeKey: key,
		cdc:      cdc,
	}
}

func NewFeeRefundHandler(am AccountKeeper, supplyKeeper auth.SupplyKeeper, rk RefundKeeper) types.FeeRefundHandler {
	return func(ctx sdk.Context, tx sdk.Tx, txResult sdk.Result) (actualCostFee sdk.Coin, err error) {
		txAccount := GetFeePayers(ctx)
		if txAccount == nil {
			return sdk.Coin{}, nil
		}

		stdTx, ok := tx.(StdTx)
		if !ok {
			return sdk.Coin{}, nil
		}
		ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())

		fee := getFee(stdTx.Fee.Amount)

		// if all gas has been consumed, then there is no need to run the fee refund process
		if txResult.GasWanted <= txResult.GasUsed {
			actualCostFee = fee
			return actualCostFee, nil
		}

		unusedGas := txResult.GasWanted - txResult.GasUsed
		refundCoin := sdk.NewCoin(fee.Denom, fee.Amount.Mul(sdk.NewInt(int64(unusedGas))).Quo(sdk.NewInt(int64(txResult.GasWanted))))
		acc := am.GetAccount(ctx, txAccount.GetAddress())

		if ctx.BlockHeight() == 0 { // fee for genesis block is 0
			return sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(0)), nil
		}
		_, err = RefundFees(supplyKeeper, ctx, acc, refundCoin)
		if err != nil {
			return sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(0)), err
		}

		return actualCostFee, nil
	}
}

func RefundFees(supplyKeeper auth.SupplyKeeper, ctx sdk.Context, acc Account, fees sdk.Coin) (*sdk.Result, error) {
	if !fees.IsValid() {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	//TODO add more validation
	err := supplyKeeper.SendCoinsFromModuleToAccount(ctx, auth.FeeCollectorName, acc.GetAddress(), sdk.NewCoins(fees))
	if err != nil {
		return nil, err
	}

	return &sdk.Result{}, nil
}

func getFee(coins sdk.Coins) sdk.Coin {
	if coins == nil || coins.Empty() {
		return sdk.NewCoin(sdk.NativeTokenName, sdk.ZeroInt())
	}
	return sdk.NewCoin(sdk.NativeTokenName, coins.AmountOf(sdk.NativeTokenName))
}
