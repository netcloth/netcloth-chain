package auth

import sdk "github.com/netcloth/netcloth-chain/types"

type contextKey int // local to the auth module

const (
	contextKeySigners contextKey = iota
)

func WithSigners(ctx sdk.Context, accounts []Account) sdk.Context {
	return ctx.WithValue(contextKeySigners, accounts)
}

func GetSigners(ctx sdk.Context) []Account {
	v := ctx.Value(contextKeySigners)
	if v == nil {
		return []Account{}
	}
	return v.([]Account)
}
