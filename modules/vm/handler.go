package vm

import (
	"fmt"
	"os"

	"github.com/netcloth/netcloth-chain/crypto"

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

	acc := k.GetAccount(ctx, msg.From)
	if acc == nil {
		return sdk.ErrInvalidAddress("account not found").Result()
	}

	//generate contract address
	contractAddr := crypto.CreateAddress(msg.From, acc.GetSequence())
	fmt.Fprintf(os.Stderr, fmt.Sprintf("contractAddr = %v, acc.GetSequence() = %v\n", contractAddr.String(), acc.GetSequence()))

	//check account's balance >= amount

	//check recur deep < 1024

	//from account seq += 1

	//crate contract account

	//transfer

	//create contract object

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

func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
