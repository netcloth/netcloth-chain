package ipal

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/modules/ipal/keeper"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgServiceNodeClaim:
			return handleMsgServiceNodeClaim(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Unrecognized Msg type: %s"+msg.Type())
		}
	}
}

func handleMsgServiceNodeClaim(ctx sdk.Context, k Keeper, m MsgServiceNodeClaim) (*sdk.Result, error) {
	m.TrimSpace()

	err := m.ValidateBasic()
	if err != nil {
		return nil, err
	}

	acc, monikerExist := k.GetServiceNodeAddByMoniker(ctx, m.Moniker)
	if monikerExist && !acc.Equals(m.OperatorAddress) {
		return nil, sdkerrors.Wrapf(ErrMonikerExist, "moniker: [%s] already exist", m.Moniker)
	}

	err = k.DoServiceNodeClaim(ctx, m)
	if err != nil {
		return nil, err
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
	matureUnstakings := k.DequeueAllMatureUnBondingQueue(ctx, ctx.BlockHeader().Time)
	for _, matureUnstaking := range matureUnstakings {
		k.DoUnbond(ctx, matureUnstaking)
	}
	return []abci.ValidatorUpdate{}
}
