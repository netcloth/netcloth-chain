package rest

import (
	"github.com/NetCloth/netcloth-chain/client/context"
	"github.com/gorilla/mux"
	"net/http"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/ipal/{accAddress}",
		mockHandler(cliCtx),
	).Methods("GET")
}

func mockHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}