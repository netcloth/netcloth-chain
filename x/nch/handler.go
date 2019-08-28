package nch

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "nch" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSend:
			return handleMsgSend(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nch Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to transfer
func handleMsgSend(ctx sdk.Context, keeper Keeper, msg MsgSend) sdk.Result {
	// transfer coins
	if !msg.Amount.IsValid() {
		return sdk.ErrInsufficientCoins("invalid coins").Result()
	}

	ctx.Logger().Info(fmt.Sprintf("transfer %s from %s to %s", msg.Amount.String(), msg.From.String(), msg.To.String()))

	// transfer coin
	err := keeper.coinKeeper.SendCoins(ctx, msg.From, msg.To, msg.Amount)
	if err != nil {
		return sdk.ErrInsufficientCoins("does not have enough coins").Result()
	}

	return sdk.Result{}
}