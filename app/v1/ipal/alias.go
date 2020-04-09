package ipal

import (
	"github.com/netcloth/netcloth-chain/app/v1/ipal/keeper"
	"github.com/netcloth/netcloth-chain/app/v1/ipal/types"
)

const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = keeper.DefaultParamspace
)

var (
	NewKeeper              = keeper.NewKeeper
	NewQuerier             = keeper.NewQuerier
	RegisterCodec          = types.RegisterCodec
	NewIPALNodeObject      = types.NewIPALNode
	NewMsgIPALNodeClaim    = types.NewMsgIPALNodeClaim
	ModuleCdc              = types.ModuleCdc
	AttributeValueCategory = types.AttributeValueCategory
	NewEndpoint            = types.NewEndpoint
	ErrEmptyInputs         = types.ErrEmptyInputs
	ErrBadDenom            = types.ErrBadDenom
	ErrBondInsufficient    = types.ErrBondInsufficient
	ErrMonikerExist        = types.ErrMonikerExist
	ErrEndpointsFormat     = types.ErrEndpointsFormat
	ErrEndpointsEmpty      = types.ErrEndpointsEmpty
	ErrEndpointsDuplicate  = types.ErrEndpointsDuplicate
)

type (
	Keeper           = keeper.Keeper
	MsgIPALNodeClaim = types.MsgIPALNodeClaim
	Endpoint         = types.Endpoint
	Endpoints        = types.Endpoints
)
