package ipal

import (
	"fmt"

	"github.com/NetCloth/netcloth-chain/modules/ipal/keeper"
	sdk "github.com/NetCloth/netcloth-chain/types"
)

func NewHandler (k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgIPALClaim:
			return handleMsgIPALClaim(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized ipal claim message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgIPALClaim(ctx sdk.Context, k Keeper, msg MsgIPALClaim) sdk.Result {
	ctx.Logger().Info("handleMsgIPALClaim ...")

	// check to see if the userAddress and serverIP has been registered before
	//if _, found := k.GetIPALObject(ctx, msg.UserAddress, msg.ServerIP); found {
	//	return ErrIPALObjectExists(k.CodeSpace()).Result()
	//}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}

}