package vm

import (
	"fmt"
	"os"

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
	var (
		ret         []byte
		leftOverGas uint64
		addr        sdk.AccAddress
		err         sdk.Error
	)

	const (
		gas = 100000000
	)

	err = msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	evmCtx := Context{}
	cfg := Config{}
	evm := NewEVM(evmCtx, ctx, k, cfg)

	ret, addr, leftOverGas, err = evm.Create(msg.From, msg.Code, gas, msg.Amount.Amount.BigInt())
	fmt.Fprint(os.Stderr, fmt.Sprintf("contractAddr = %s, leftOverGas = %v, err = %v\n", addr, leftOverGas, err))

	if err != nil {
		return sdk.ErrInternal("contract deploy err").Result()
	}

	return sdk.Result{Data: ret, GasUsed: gas - leftOverGas}
}

func handleMsgContractCall(ctx sdk.Context, msg MsgContractCall, k Keeper) sdk.Result {
	var (
		ret         []byte
		leftOverGas uint64
		err         sdk.Error
	)
	const (
		gas = 100000000
	)

	err = msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	evmCtx := Context{}
	cfg := Config{}
	evm := NewEVM(evmCtx, ctx, k, cfg)

	ret, leftOverGas, err = evm.Call(msg.From, msg.Recipient, msg.Payload, gas, msg.Amount.Amount.BigInt())
	fmt.Fprint(os.Stderr, fmt.Sprintf("ret = %x, leftOverGas = %v, err = %v\n", ret, leftOverGas, err))

	return sdk.Result{Data: ret, GasUsed: gas - leftOverGas}
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
