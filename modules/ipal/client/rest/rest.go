package rest

import (
    "github.com/gorilla/mux"

    "github.com/NetCloth/netcloth-chain/client/context"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
    registerQueryRoutes(cliCtx, r)
}
