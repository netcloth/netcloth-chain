package types

import "math/big"
import sdk "github.com/netcloth/netcloth-chain/types"


// StateTransition defines data to transitionDB in vm
type StateTransition struct {
}

func (st StateTransition) TransitionCSDB(ctx sdk.Context) (*big.Int, sdk.Result) {
	ctx.Logger().Info("TransitionCSDB ...")
	return nil, nil
}