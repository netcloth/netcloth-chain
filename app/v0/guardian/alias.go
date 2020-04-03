package guardian

import (
	"github.com/netcloth/netcloth-chain/app/v0/guardian/types"
)

const (
	Genesis  = types.Genesis
	Ordinary = types.Ordinary

	ModuleName     = types.ModuleName
	RouterKey      = types.RouterKey
	QuerierRoute   = types.QuerierRoute
	QueryProfilers = types.QueryProfilers
)

type (
	MsgAddProfiler    = types.MsgAddProfiler
	MsgDeleteProfiler = types.MsgDeleteProfiler
	Guardian          = types.Guardian
	Profilers         = types.Profilers
)

var (
	NewMsgAddProfiler       = types.NewMsgAddProfiler
	NewMsgDeleteProfiler    = types.NewMsgDeleteProfiler
	NewGuardian             = types.NewGuardian
	GetProfilerKey          = types.GetProfilerKey
	GetProfilersSubspaceKey = types.GetProfilersSubspaceKey

	ErrInvalidOperator       = types.ErrInvalidOperator
	ErrProfilerNotExists     = types.ErrProfilerNotExists
	ErrDeleteGenesisProfiler = types.ErrDeleteGenesisProfiler
	ErrProfilerExists        = types.ErrProfilerExists
	ErrInvalidDescription    = types.ErrInvalidDescription
	ErrAddressEmpty          = types.ErrAddressEmpty
	ErrAddedByEmpty          = types.ErrAddedByEmpty
	ErrDeletedByEmpty        = types.ErrDeletedByEmpty
)
