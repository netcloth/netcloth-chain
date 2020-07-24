package slashing

import (
	"github.com/netcloth/netcloth-chain/app/v0/auth/exported"
	sdk "github.com/netcloth/netcloth-chain/types"
)

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) exported.Account
}
