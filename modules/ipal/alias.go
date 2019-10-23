package ipal

import (
	"github.com/NetCloth/netcloth-chain/modules/ipal/keeper"
	"github.com/NetCloth/netcloth-chain/modules/ipal/types"
)

const (
	StoreKey            = types.StoreKey
	RouterKey           = types.RouterKey
	QuerierRoute        = types.QuerierRoute
	DefaultCodespace    = types.DefaultCodespace
	DefaultParamspace   = keeper.DefaultParamspace
)

var (
	// functions aliases
	RegisterCodec                  = types.RegisterCodec
	ErrEmptyInputs                 = types.ErrEmptyInputs
	ErrStringTooLong               = types.ErrStringTooLong
	ErrInvalidSignature            = types.ErrInvalidSignature
	ErrIPALClaimUserRequestExpired = types.ErrIPALClaimUserRequestExpired

	NewIPALObject       = types.NewIPALObject
	NewServerNodeObject = types.NewServerNodeObject
	NewQuerier = keeper.NewQuerier

	NewMsgIPALClaim        = types.NewMsgIPALClaim
	NewMsgServiceNodeClaim = types.NewMsgServiceNodeClaim

	ModuleName = types.ModuleName

	// variable aliases
	ModuleCdc = types.ModuleCdc

	NewKeeper = keeper.NewKeeper

	AttributeValueCategory = types.AttributeValueCategory
)

type (
	Keeper = keeper.Keeper

	MsgIPALClaim        = types.MsgIPALClaim
	MsgServiceNodeClaim = types.MsgServiceNodeClaim
)
