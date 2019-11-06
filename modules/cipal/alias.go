package cipal

import (
	"github.com/NetCloth/netcloth-chain/modules/cipal/keeper"
	"github.com/NetCloth/netcloth-chain/modules/cipal/types"
)

const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultCodespace  = types.DefaultCodespace
	DefaultParamspace = keeper.DefaultParamspace
)

var (
	RegisterCodec                  = types.RegisterCodec
	NewIPALObject                  = types.NewCIPALObject
	NewQuerier                     = keeper.NewQuerier
	NewADParam                     = types.NewADParam
	NewIPALUserRequest             = types.NewCIPALUserRequest
	NewMsgIPALClaim                = types.NewMsgCIPALClaim
	NewKeeper                      = keeper.NewKeeper
	ErrEmptyInputs                 = types.ErrEmptyInputs
	ErrStringTooLong               = types.ErrStringTooLong
	ErrInvalidSignature            = types.ErrInvalidSignature
	ErrIPALClaimUserRequestExpired = types.ErrIPALClaimUserRequestExpired
	ModuleCdc                      = types.ModuleCdc
	AttributeValueCategory         = types.AttributeValueCategory
)

type (
	Keeper          = keeper.Keeper
	MsgIPALClaim    = types.MsgCIPALClaim
	IPALUserRequest = types.CIPALUserRequest
	ADParam         = types.ADParam
)
