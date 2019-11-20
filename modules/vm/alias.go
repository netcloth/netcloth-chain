package vm

import (
	"github.com/netcloth/netcloth-chain/modules/vm/keeper"
	"github.com/netcloth/netcloth-chain/modules/vm/types"
)

type (
	Keeper = keeper.Keeper

	StateTransition = types.StateTransition

	MsgContractCreate = types.MsgContractCreate
	MsgContractCall   = types.MsgContractCall
)
