package client

import (
	"github.com/gorilla/mux"

	"github.com/netcloth/netcloth-chain/client/context"
)

// Register routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	RegisterRPCRoutes(cliCtx, r)
}
