package vm

import (
	"fmt"
	"os"

	"github.com/netcloth/netcloth-chain/modules/vm/types"

	"github.com/netcloth/netcloth-chain/modules/vm/keeper"
	sdk "github.com/netcloth/netcloth-chain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgContractCreate:
			return handleMsgContractCreate(ctx, msg, k)
		case MsgContractCall:
			return handleMsgContractCall(ctx, msg, k)
		default:
			errMsg := fmt.Sprintf("Unrecognized Msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgContractCreate(ctx sdk.Context, msg MsgContractCreate, k Keeper) sdk.Result {
	ctx.Logger().Info("handleMsgContractCreate ...")
	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	//err = k.DoContractCreate(ctx, msg)
	//if err != nil {
	//	return err.Result()
	//}

	evmCtx := NewEVMContext(ctx, msg.From)
	db := types.NewCommitStateDB(ctx, k.AK, k.BK, k)
	cfg := Config{}
	evm := NewEVM(evmCtx, *db, cfg)
	d1, d2, d3, e := evm.Create(msg.From, msg.Code, 100000000, sdk.NewInt(10000).BigInt())

	fmt.Fprint(os.Stderr, fmt.Sprintf("d1 = %v, d2 = %v, d3 = %v, e = %v\n", d1, d2, d3, e))
	if e != nil {
		return sdk.ErrInternal("contract deploy err").Result()
	}
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgContractCall(ctx sdk.Context, msg MsgContractCall, k Keeper) sdk.Result {
	ctx.Logger().Info("handleMsgContractCall ...")

	st := StateTransition{}
	_, res := st.TransitionCSDB(ctx)
	return res
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
