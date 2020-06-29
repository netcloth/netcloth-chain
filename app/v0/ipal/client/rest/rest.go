package rest

import (
	"github.com/gorilla/mux"

	"github.com/netcloth/netcloth-chain/client/context"
)

// RegisterRoutes registers the routes from the different modules for the LCD.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
}
