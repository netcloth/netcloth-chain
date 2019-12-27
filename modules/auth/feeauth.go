package auth

import (
	auth "github.com/netcloth/netcloth-chain/modules/auth/types"
	"github.com/netcloth/netcloth-chain/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

type FeeAuth struct {
	NativeFeeDenom string `json:"native_fee_denom"`
}

func NewFeeAuth(nativeFeeDenon string) FeeAuth {
	return FeeAuth{NativeFeeDenom: nativeFeeDenon}
}

func InitialFeeAuth() FeeAuth {
	return NewFeeAuth(sdk.NativeTokenName)
}

// NewFeeRefundHandler creates a fee token refund handler
func NewFeeRefundHandler(am AccountKeeper, supplyKeeper auth.SupplyKeeper, fk FeeKeeper) types.FeeRefundHandler {
	return func(ctx sdk.Context, tx sdk.Tx, txResult sdk.Result) (actualCostFee sdk.Coin, err error) {
		txAccounts := GetSigners(ctx)
		if len(txAccounts) < 1 {
			return sdk.Coin{}, nil
		}
		firstAccount := txAccounts[0]

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
		acc := am.GetAccount(ctx, firstAccount.GetAddress())

		res := RefundFees(supplyKeeper, ctx, acc, refundCoin)
		if !res.IsOK() {
			return actualCostFee, nil
		}

		return actualCostFee, nil
	}
}

func getFee(coins sdk.Coins) sdk.Coin {
	if coins == nil || coins.Empty() {
		return sdk.NewCoin(sdk.NativeTokenName, sdk.ZeroInt())
	}
	return sdk.NewCoin(sdk.NativeTokenName, coins.AmountOf(sdk.NativeTokenName))
}
