package token

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgIssue:
			return handleMsgIssue(ctx, k, msg)
		default:
			errMsg := "Unrecognized Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle handleMsgIssue
func handleMsgIssue(ctx sdk.Context, k Keeper, msg MsgIssue) sdk.Result {
	// check issue amount
	if !msg.Amount.IsValid() {
		return sdk.ErrInsufficientCoins("invalid coins").Result()
	}

	// check coin name

	ctx.Logger().Debug("%s issue %s to %s ", msg.Banker.String(), msg.Amount.String(), msg.Address.String())

	newCoins := sdk.NewCoins(msg.Amount)
	_, err := k.coinKeeper.AddCoins(ctx, msg.Address, newCoins)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}