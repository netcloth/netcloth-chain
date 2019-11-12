package cipal

import (
	"github.com/NetCloth/netcloth-chain/modules/cipal/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/NetCloth/netcloth-chain/modules/cipal/keeper"
	sdk "github.com/NetCloth/netcloth-chain/types"
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
	if ctx.BlockHeader().Time.After(msg.UserRequest.Params.Expiration) {
		return ErrIPALClaimUserRequestExpired("user request expired").Result()
	}

	sigVerifyPass := msg.UserRequest.Sig.VerifyBytes(msg.UserRequest.Params.GetSignBytes(), msg.UserRequest.Sig.Signature)
	if !sigVerifyPass {
		return ErrCIPALClaimUserRequestSigVerify("user signature verify failed").Result()
	}

	obj, found := k.GetIPALObject(ctx, msg.UserRequest.Params.UserAddress)
	if found {
		updateIndex := -1
		var si types.ServiceInfo
		for i, v := range obj.ServiceInfos {
			if v.Type == msg.UserRequest.Params.ServiceInfo.Type {
				updateIndex = i
				si = v
				break
			}
		}

		if updateIndex != -1 {
			if si.Address != msg.UserRequest.Params.ServiceInfo.Address {
				obj.ServiceInfos[updateIndex].Address = msg.UserRequest.Params.ServiceInfo.Address
			}
		} else {
			obj.ServiceInfos = append(obj.ServiceInfos, msg.UserRequest.Params.ServiceInfo)
		}

		k.SetIPALObject(ctx, obj)
	} else {
		obj = NewIPALObject(msg.UserRequest.Params.UserAddress, msg.UserRequest.Params.ServiceInfo.Address, msg.UserRequest.Params.ServiceInfo.Type)
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
