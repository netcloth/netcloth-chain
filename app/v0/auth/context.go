package auth

import sdk "github.com/netcloth/netcloth-chain/types"

type contextKey int // local to the auth module

const (
	contextKeyFeePayers contextKey = iota
)

func WithFeePayers(ctx sdk.Context, account Account) sdk.Context {
	return ctx.WithValue(contextKeyFeePayers, account)
}

func GetFeePayers(ctx sdk.Context) Account {
	v := ctx.Value(contextKeyFeePayers)
	if v == nil {
		return nil
	}
	return v.(Account)
}
