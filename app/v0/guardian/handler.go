package guardian

import (
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case MsgAddProfiler:
			return handleMsgAddProfiler(ctx, k, msg)
		case MsgDeleteProfiler:
			return handleMsgDeleteProfiler(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleMsgAddProfiler(ctx sdk.Context, k Keeper, msg MsgAddProfiler) (*sdk.Result, error) {
	if profiler, found := k.GetProfiler(ctx, msg.AddedBy); !found || profiler.AccountType != Genesis {
		return nil, ErrInvalidOperator(msg.AddedBy)
	}

	if _, found := k.GetProfiler(ctx, msg.Address); found {
		return nil, ErrProfilerExists(msg.Address)
	}

	profiler := NewGuardian(msg.Description, Ordinary, msg.Address, msg.AddedBy)
	err := k.AddProfiler(ctx, profiler)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{}, nil
}

func handleMsgDeleteProfiler(ctx sdk.Context, k Keeper, msg MsgDeleteProfiler) (*sdk.Result, error) {
	if profiler, found := k.GetProfiler(ctx, msg.DeletedBy); !found || profiler.AccountType != Genesis {
		return nil, ErrInvalidOperator(msg.DeletedBy)
	}

	profiler, found := k.GetProfiler(ctx, msg.Address)
	if !found {
		return nil, ErrProfilerNotExists(msg.Address)
	}

	if profiler.AccountType == Genesis {
		return nil, ErrDeleteGenesisProfiler(msg.Address)
	}

	err := k.DeleteProfiler(ctx, msg.Address)
	if err != nil {
		return nil, err
	}
	return &sdk.Result{}, nil
}
