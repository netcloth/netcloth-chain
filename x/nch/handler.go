package nch

import (
	"fmt"
	"github.com/NetCloth/netcloth-chain/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "nch" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgTransfer:
			return handleMsgTransfer(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nch Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to transfer
func handleMsgTransfer(ctx sdk.Context, keeper Keeper, msg MsgTransfer) sdk.Result {
	// transfer coins
	if !msg.Value.IsValid() {
		return sdk.ErrInsufficientCoins("invalid coins").Result()
	}

	// check fee
	if msg.Fee.AmountOf(types.AppCoin).Int64() < MinTransferFee {
		return sdk.ErrInsufficientCoins("insufficient fee").Result()
	}

	// substract fee
	_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.From, msg.Fee)
	if err != nil {
		return sdk.ErrInsufficientCoins("does not have enough coins for fee").Result()
	}

	ctx.Logger().Info(fmt.Sprintf("transfer %s from %s to %s", msg.Value.String(), msg.From.String(), msg.To.String()))

	// transfer coin
	_, err = keeper.coinKeeper.SendCoins(ctx, msg.From, msg.To, msg.Value)
	if err != nil {
		return sdk.ErrInsufficientCoins("does not have enough coins").Result()
	}

	return sdk.Result{}
}