package auth

import (
	"fmt"
	"math"

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

func (fa FeeAuth) getNativeFeeToken(ctx sdk.Context, coins sdk.Coins) sdk.Coin {
	if coins == nil {
		return sdk.NewCoin(sdk.NativeTokenName, sdk.ZeroInt())
	}
	return sdk.NewCoin(sdk.NativeTokenName, coins.AmountOf(sdk.NativeTokenName))
}

func (fa FeeAuth) feePreprocess(ctx sdk.Context, coins sdk.Coins, gasLimit uint64) sdk.Error {
	if gasLimit == 0 || int64(gasLimit) < 0 {
		return ErrInvalidGas(fmt.Sprintf("gaslimit %d should be positive and no more than %d", gasLimit, math.MaxInt64))
	}
	return nil
}

// NewFeePreprocessHandler creates a fee token refund handler
func NewFeeRefundHandler(am AccountKeeper, fk FeeKeeper) types.FeeRefundHandler {
	return func(ctx sdk.Context, tx sdk.Tx, txResult sdk.Result) (actualCostFee sdk.Coin, err error) {
		//TODO

		actualCostFee = sdk.NewCoin(sdk.NativeTokenName, sdk.ZeroInt())
		return actualCostFee, nil
	}
}

// NewFeePreprocessHandler creates a fee token preprocesser
func NewFeePreprocessHandler(fk FeeKeeper) types.FeePreprocessHandler {
	return func(ctx sdk.Context, tx sdk.Tx) sdk.Error {
		stdTx, ok := tx.(StdTx)
		if !ok {
			return sdk.ErrInternal("tx must be StdTx")
		}

		fa := fk.GetFeeAuth(ctx)
		totalNativeFee := fa.getNativeFeeToken(ctx, stdTx.Fee.Amount)

		return fa.feePreprocess(ctx, sdk.Coins{totalNativeFee}, stdTx.Fee.Gas)
	}
}

func getFee(coins sdk.Coins) sdk.Coin {
	if coins == nil || coins.Empty() {
		return sdk.NewCoin(sdk.NativeTokenName, sdk.ZeroInt())
	}

	return sdk.NewCoin(sdk.NativeTokenName, coins.AmountOf(sdk.NativeTokenName))
}

func checkFee(params Params, coins sdk.Coins, gasLimit uint64) sdk.Error {
	if gasLimit == 0 || int64(gasLimit) < 0 {
		return ErrInvalidGas(fmt.Sprintf("gaslimit %d should be positive and no more than %d", gasLimit, math.MaxInt64))
	}
	//TODO

	return nil
}
