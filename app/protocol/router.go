package protocol

import (
	"fmt"

	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/utils"
)

// Router provides handlers for each transaction type
type Router struct {
	routes map[string]sdk.Handler
}

var _ sdk.Router = NewRouter()

// NewQueryRouter returns a reference to a new QueryRouter
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]sdk.Handler),
	}
}

// AddRoute adds a query path to the router with a given Querier
func (rtr *Router) AddRoute(path string, h sdk.Handler) sdk.Router {
	if !utils.IsAlphaNumeric(path) {
		panic("route expressions can only contain alphanumeric characters")
	}
	if rtr.routes[path] != nil {
		panic(fmt.Sprintf("route %s has already been initialized", path))
	}

	rtr.routes[path] = h
	return rtr
}

// Route returns the Querier for a given query route path
func (rtr *Router) Route(_ sdk.Context, path string) sdk.Handler {
	return rtr.routes[path]
}
