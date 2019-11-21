package types

import "math/big"
import sdk "github.com/netcloth/netcloth-chain/types"

// StateTransition defines data to transitionDB in vm
type StateTransition struct {
	Sender    sdk.AccAddress
	Price     sdk.Int
	GasLimit  sdk.Int
	Recipient sdk.AccAddress
	Amount    sdk.Int
	Payload   []byte
}

func (st StateTransition) TransitionCSDB(ctx sdk.Context) (*big.Int, sdk.Result) {
	ctx.Logger().Info("TransitionCSDB ...")
	return nil, sdk.Result{}
}
