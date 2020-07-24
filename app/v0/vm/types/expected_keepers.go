package types

import (
	"github.com/netcloth/netcloth-chain/app/v0/auth/exported"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// AccountKeeper defines the expected account keeper used for vm
type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) exported.Account
	RemoveAccount(ctx sdk.Context, acc exported.Account)
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) exported.Account
	SetAccount(ctx sdk.Context, acc exported.Account)
}

// BankKeeper defines the expected bank keeper used for vm
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}
