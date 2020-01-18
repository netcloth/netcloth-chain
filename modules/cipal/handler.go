package cipal

import (
	"github.com/netcloth/netcloth-chain/modules/cipal/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/modules/cipal/keeper"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgIPALClaim:
			return handleMsgIPALClaim(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleMsgIPALClaim(ctx sdk.Context, k Keeper, msg MsgIPALClaim) (*sdk.Result, error) {
	if ctx.BlockHeader().Time.After(msg.UserRequest.Params.Expiration) {
		return nil, sdkerrors.Wrap(ErrIPALClaimUserRequestExpired, "user request expired")
	}

	sigVerifyPass := msg.UserRequest.Sig.VerifyBytes(msg.UserRequest.Params.GetSignBytes(), msg.UserRequest.Sig.Signature)
	if !sigVerifyPass {
		return nil, sdkerrors.Wrap(ErrCIPALClaimUserRequestSigVerify, "user signature verify failed")
	}

	obj, found := k.GetCIPALObject(ctx, msg.UserRequest.Params.UserAddress)
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

		k.SetCIPALObject(ctx, obj)
	} else {
		obj = NewIPALObject(msg.UserRequest.Params.UserAddress, msg.UserRequest.Params.ServiceInfo.Address, msg.UserRequest.Params.ServiceInfo.Type)
		k.SetCIPALObject(ctx, obj)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
