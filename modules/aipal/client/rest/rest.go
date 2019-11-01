package rest

import (
    "github.com/NetCloth/netcloth-chain/client/context"
    "github.com/gorilla/mux"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
    registerQueryRoutes(cliCtx, r)
}
