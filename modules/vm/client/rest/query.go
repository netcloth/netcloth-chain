package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/netcloth/netcloth-chain/modules/vm/types"
	"github.com/netcloth/netcloth-chain/types/rest"

	"github.com/netcloth/netcloth-chain/client/context"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/vm/storage/{addr}/{key}",
		getStorageFn(cliCtx),
	).Methods("GET")
}

func queryStorage(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addr := vars["addr"]
		key := vars["key"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/vm/%s/%s/%s", types.QueryStorage, addr, key)
		res, height, err := cliCtx.Query(route)
		if err != nil {
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getStorageFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryStorage(cliCtx)
}
