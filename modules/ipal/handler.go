package ipal

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/modules/ipal/keeper"
	"github.com/netcloth/netcloth-chain/modules/ipal/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgServiceNodeClaim:
			return handleMsgServiceNodeClaim(ctx, k, msg)
		default:
			errMsg := "Unrecognized Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgServiceNodeClaim(ctx sdk.Context, k Keeper, m MsgServiceNodeClaim) sdk.Result {
	m.TrimSpace()

	err := m.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	acc, monikerExist := k.GetServiceNodeAddByMoniker(ctx, m.Moniker)
	if monikerExist && !acc.Equals(m.OperatorAddress) {
		return types.ErrMonikerExist(fmt.Sprintf("moniker: [%s] already exist", m.Moniker)).Result()
	}

	err = k.DoServiceNodeClaim(ctx, m)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	)

	return sdk.Result{}
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	matureUnstakings := k.DequeueAllMatureUnBondingQueue(ctx, ctx.BlockHeader().Time)
	for _, matureUnstaking := range matureUnstakings {
		k.DoUnbond(ctx, matureUnstaking)
	}
	return []abci.ValidatorUpdate{}
}
