package ipal

import (
	"github.com/NetCloth/netcloth-chain/modules/ipal/keeper"
	sdk "github.com/NetCloth/netcloth-chain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgIPALClaim:
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
		return ErrIPALClaimUserRequestExpired("user request expired").Result()
	}

	// check to see if the userAddress and serverIP has been registered before
	obj, found := k.GetIPALObject(ctx, msg.UserRequest.Params.UserAddress, msg.UserRequest.Params.ServerIP)
	if found {
		// update ipal object
		obj.ServerIP = msg.UserRequest.Params.ServerIP
		k.SetIPALObject(ctx, obj)
	} else {
		// create new ipal object
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

func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}