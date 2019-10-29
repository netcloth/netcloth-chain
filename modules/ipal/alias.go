package ipal

import (
	"github.com/NetCloth/netcloth-chain/modules/ipal/keeper"
	"github.com/NetCloth/netcloth-chain/modules/ipal/types"
)

const (
	ModuleName		    = types.ModuleName
	StoreKey            = types.StoreKey
	RouterKey           = types.RouterKey
	QuerierRoute        = types.QuerierRoute
	DefaultCodespace    = types.DefaultCodespace
	DefaultParamspace   = keeper.DefaultParamspace
)

var (
	RegisterCodec					= types.RegisterCodec
	NewIPALObject       			= types.NewIPALObject
	NewQuerier						= keeper.NewQuerier
	NewMsgIPALClaim					= types.NewMsgIPALClaim
	NewKeeper 						= keeper.NewKeeper
	ErrEmptyInputs					= types.ErrEmptyInputs
	ErrStringTooLong				= types.ErrStringTooLong
	ErrInvalidSignature				= types.ErrInvalidSignature
	ErrIPALClaimUserRequestExpired	= types.ErrIPALClaimUserRequestExpired
	ModuleCdc 						= types.ModuleCdc
	AttributeValueCategory 			= types.AttributeValueCategory
)

type (
	Keeper = keeper.Keeper
	MsgIPALClaim = types.MsgIPALClaim
)
