package protocol

import (
	"fmt"
	sdk "github.com/netcloth/netcloth-chain/types"
)

type queryRouter struct {
	routes map[string]sdk.Querier
}

var _ sdk.QueryRouter = NewQueryRouter()

func NewQueryRouter() *queryRouter { // nolint: golint
	return &queryRouter{
		routes: map[string]sdk.Querier{},
	}
}

func (qrt *queryRouter) AddRoute(path string, q sdk.Querier) sdk.QueryRouter {
	if !isAlphaNumeric(path) {
		panic("route expressions can only contain alphanumeric characters")
	}
	if qrt.routes[path] != nil {
		panic(fmt.Sprintf("route %s has already been initialized", path))
	}

	qrt.routes[path] = q
	return qrt
}

func (qrt *queryRouter) Route(path string) sdk.Querier {
	return qrt.routes[path]
}
