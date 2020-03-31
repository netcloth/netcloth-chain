package protocol

import (
	"fmt"
	sdk "github.com/netcloth/netcloth-chain/types"
	"regexp"
)

var isAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

type Router struct {
	routes map[string]sdk.Handler
}

var _ sdk.Router = NewRouter()

// NewRouter returns a reference to a new router.
//
// TODO: Either make the function private or make return type (router) public.
func NewRouter() *Router { // nolint: golint
	return &Router{
		routes: make(map[string]sdk.Handler),
	}
}

// AddRoute adds a route path to the router with a given handler. The route must
// be alphanumeric.
func (rtr *Router) AddRoute(path string, h sdk.Handler) sdk.Router {
	if !isAlphaNumeric(path) {
		panic("route expressions can only contain alphanumeric characters")
	}
	if rtr.routes[path] != nil {
		panic(fmt.Sprintf("route %s has already been initialized", path))
	}

	rtr.routes[path] = h
	return rtr
}

// Route returns a handler for a given route path.
//
// TODO: Handle expressive matches.
func (rtr *Router) Route(_ sdk.Context, path string) sdk.Handler {
	return rtr.routes[path]
}
