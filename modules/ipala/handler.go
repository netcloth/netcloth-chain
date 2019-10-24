package ipala

import (
    "github.com/NetCloth/netcloth-chain/modules/ipala/keeper"
    sdk "github.com/NetCloth/netcloth-chain/types"
    abci "github.com/tendermint/tendermint/abci/types"
)

func NewHandler(k Keeper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
        ctx = ctx.WithEventManager(sdk.NewEventManager())

        switch msg := msg.(type) {
        case MsgServiceNodeClaim:
            return handleMsgServerNodeClaim(ctx, k, msg)
        default:
            errMsg := "Unrecognized Msg type: %s" + msg.Type()
            return sdk.ErrUnknownRequest(errMsg).Result()
        }
    }
}

func handleMsgServerNodeClaim(ctx sdk.Context, k Keeper, m MsgServiceNodeClaim) sdk.Result {
    err := m.ValidateBasic()
    if err != nil {
        return err.Result()
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
        k.DoUnBond(ctx, matureUnstaking)
    }
    return []abci.ValidatorUpdate{}
}