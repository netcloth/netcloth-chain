package types

import (
	"github.com/netcloth/netcloth-chain/modules/auth/exported"
	sdk "github.com/netcloth/netcloth-chain/types"
)

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) exported.Account
}
