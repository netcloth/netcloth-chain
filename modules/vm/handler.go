package vm

import (
	"fmt"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgContractCreate:
			return handleMMsgContractCreate(ctx, msg, k)

		case MsgContractCall:
			return handleMsgContractCall(ctx, msg, k)
		default:
			errMsg := fmt.Sprintf("Unrecognized Msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMMsgContractCreate(ctx sdk.Context, msg MsgContractCreate, k Keeper) sdk.Result {
	ctx.Logger().Info("handleMMsgContractCreate ...")

	st := StateTransition{}
	_, res := st.TransitionCSDB(ctx)
	return res
}

func handleMsgContractCall(ctx sdk.Context, msg MsgContractCall, k Keeper) sdk.Result {
	ctx.Logger().Info("handleMsgContractCall ...")

	st := StateTransition{}
	_, res := st.TransitionCSDB(ctx)
	return res
}
