package vm

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/modules/vm/keeper"
	"github.com/netcloth/netcloth-chain/modules/vm/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgContract:
			return handleMsgContract(ctx, msg, k)
		default:
			errMsg := fmt.Sprintf("Unrecognized Msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgContract(ctx sdk.Context, msg MsgContract, k Keeper) sdk.Result {
	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	gasLimit := ctx.GasMeter().Limit() - ctx.GasMeter().GasConsumed()
	_, res := DoStateTransition(ctx, msg, k, gasLimit, false)
	if !res.IsOK() {
		return res
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Data: res.Data, GasUsed: res.GasUsed, Events: ctx.EventManager().Events()}
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	k.StateDB.UpdateAccounts() //update account balance for fee refund when create/call contract
	k.StateDB.WithContext(ctx).Commit(true)
	return []abci.ValidatorUpdate{}
}
