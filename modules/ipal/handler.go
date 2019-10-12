package ipal

import (
	"fmt"
	"github.com/NetCloth/netcloth-chain/modules/ipal/types"
	sdk "github.com/NetCloth/netcloth-chain/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgIPALClaim:
			return handleMsgIPALClaim(ctx, k, msg)
		default:
			errMsg := "Unrecognized Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgIPALClaim(ctx sdk.Context, k Keeper, msg MsgIPALClaim) sdk.Result {
	// check user request expiration
	if ctx.BlockHeader().Time.After(msg.UserRequest.Params.Expiration) {
		return ErrIPALClaimExpired(DefaultCodespace).Result()
	}

	// check to see if the userAddress and serverIP has been registered before
	obj, found := k.GetIPALObject(ctx, msg.UserRequest.Params.UserAddress, msg.UserRequest.Params.ServerIP)
	if found {
		// update ipal object
		ctx.Logger().Info(fmt.Sprintf("update ipal object, key: %s", msg.UserRequest.Params.UserAddress))
		obj.ServerIP = msg.UserRequest.Params.ServerIP
		k.SetIPALObject(ctx, obj)
	} else {
		// create new ipal object
		ctx.Logger().Info(fmt.Sprintf("create ipal object, key: %s", msg.UserRequest.Params.UserAddress))
		obj = NewIPALObject(msg.UserRequest.Params.UserAddress, msg.UserRequest.Params.ServerIP)
		k.SetIPALObject(ctx, obj)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}