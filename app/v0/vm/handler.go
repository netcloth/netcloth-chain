package vm

import (
	"github.com/netcloth/netcloth-chain/app/v0/vm/types"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgContract:
			return handleMsgContract(ctx, msg, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleMsgContract(ctx sdk.Context, msg MsgContract, k Keeper) (*sdk.Result, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	_, res, err := DoStateTransition(ctx, msg, k, ctx.Simulate)
	if err != nil {
		return &sdk.Result{Data: res.Data, GasUsed: res.GasUsed}, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return &sdk.Result{Data: res.Data, GasUsed: res.GasUsed, Events: ctx.EventManager().Events()}, nil
}
