package types

import (
	authexported "github.com/netcloth/netcloth-chain/app/v0/auth/exported"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	NewAccount(sdk.Context, authexported.Account) authexported.Account
	SetAccount(sdk.Context, authexported.Account)
	IterateAccounts(ctx sdk.Context, process func(authexported.Account) (stop bool))
}
