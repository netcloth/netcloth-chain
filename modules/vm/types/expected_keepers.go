package types

import (
	"github.com/netcloth/netcloth-chain/modules/auth/exported"
	sdk "github.com/netcloth/netcloth-chain/types"
)

type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) exported.Account
	GetAccount(ctx sdk.Context, addr sdk.Address) exported.Account
	SetAccount(ctx sdk.Context, acc exported.Account)
}

type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
}
