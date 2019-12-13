package vm

import (
	"fmt"

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
	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	st := StateTransition{
		Sender:    msg.From,
		Recipient: nil,
		Price:     sdk.NewInt(1000000),
		GasLimit:  10000000,
		Amount:    msg.Amount.Amount,
		Payload:   msg.Code,
		CSDB:      k.CSDB.WithContext(ctx),
	}
	_, res := st.TransitionCSDB(ctx)
	return res
}

func handleMsgContractCall(ctx sdk.Context, msg MsgContractCall, k Keeper) sdk.Result {
	ctx.Logger().Info("handleMsgContractCall ...")

	st := StateTransition{
		Sender:    msg.From,
		Recipient: msg.Recipient,
		Price:     sdk.NewInt(1000000),
		GasLimit:  10000000,
		Payload:   msg.Payload,
		Amount:    msg.Amount.Amount,
		CSDB:      k.CSDB.WithContext(ctx),
	}
	_, res := st.TransitionCSDB(ctx)
	return res
}

func EndBlocker(ctx sdk.Context, keeper Keeper) []abci.ValidatorUpdate {
	// Gas costs are handled within msg handler so costs should be ignored
	ebCtx := ctx.WithBlockGasMeter(sdk.NewInfiniteGasMeter())

	// Update account balances before committing other parts of state
	keeper.CSDB.UpdateAccounts()

	// Commit state objects to KV store
	_, err := keeper.CSDB.WithContext(ebCtx).Commit(true)
	if err != nil {
		panic(err)
	}

	// Clear accounts cache after account data has been committed
	keeper.CSDB.ClearStateObjects()

	return []abci.ValidatorUpdate{}
}
