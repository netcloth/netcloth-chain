package auth

import (
	auth "github.com/netcloth/netcloth-chain/modules/auth/types"
	"github.com/netcloth/netcloth-chain/types"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
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

		if ctx.BlockHeight() == 0 { // fee for genesis block is 0
			return sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(0)), nil
		}
		_, err = RefundFees(supplyKeeper, ctx, acc, refundCoin)
		if err != nil {
			return actualCostFee, err
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
