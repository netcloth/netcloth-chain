package rest

import (
    "net/http"

    "github.com/gorilla/mux"

    "github.com/NetCloth/netcloth-chain/client/context"
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
