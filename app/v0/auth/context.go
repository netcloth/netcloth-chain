package auth

import sdk "github.com/netcloth/netcloth-chain/types"

type contextKey int // local to the auth module

const (
	contextKeyFeePayers contextKey = iota
)

// WithFeePayers - initialize Context with FeePayer
func WithFeePayers(ctx sdk.Context, account Account) sdk.Context {
	return ctx.WithValue(contextKeyFeePayers, account)
}

// GetFeePayers - get FeePayer from Context
func GetFeePayers(ctx sdk.Context) Account {
	v := ctx.Value(contextKeyFeePayers)
	if v == nil {
		return nil
	}
	return v.(Account)
}
