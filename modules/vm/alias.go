package vm

import (
	"github.com/netcloth/netcloth-chain/modules/vm/common"
	"github.com/netcloth/netcloth-chain/modules/vm/keeper"
	"github.com/netcloth/netcloth-chain/modules/vm/types"
)

const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	TStoreKey         = types.TStoreKey
	QuerierRoute      = types.QuerierRoute
	DefaultCodespace  = types.DefaultCodespace
	DefaultParamspace = keeper.DefaultParamspace
)

type (
	Keeper = keeper.Keeper

	StateTransition = types.StateTransition

	MsgContractCreate = types.MsgContractCreate
	MsgContractCall   = types.MsgContractCall

	CommitStateDB = types.CommitStateDB
)

var (
	NewKeeper = keeper.NewKeeper

	CreateAddress = common.CreateAddress

	ErrOutOfGas                 = types.ErrOutOfGas
	ErrCodeStoreOutOfGas        = types.ErrCodeStoreOutOfGas
	ErrDepth                    = types.ErrDepth
	ErrTraceLimitReached        = types.ErrTraceLimitReached
	ErrInsufficientBalance      = types.ErrInsufficientBalance
	ErrContractAddressCollision = types.ErrContractAddressCollision
	ErrNoCompatibleInterpreter  = types.ErrNoCompatibleInterpreter
)
