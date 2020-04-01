package types

import "regexp"

var IsAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

type Router interface {
	AddRoute(r string, h Handler) Router
	Route(ctx Context, path string) Handler
}

type QueryRouter interface {
	AddRoute(r string, h Querier) QueryRouter
	Route(path string) Querier
}
