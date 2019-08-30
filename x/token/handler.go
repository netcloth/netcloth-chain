package token

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgTokenIssue:
			return handleMsgTokenIssue(ctx, k, msg)
		default:
			errMsg := "Unrecognized bank Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle handleMsgTokenIssue.
func handleMsgTokenIssue(ctx sdk.Context, k Keeper, msg MsgTokenIssue) sdk.Result {
	fmt.Println("handleMsgTokenIssue ...")
}