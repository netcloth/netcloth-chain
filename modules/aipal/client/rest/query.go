package rest

import (
	"net/http"

	"github.com/NetCloth/netcloth-chain/client/context"
	"github.com/NetCloth/netcloth-chain/modules/aipal/types"
	"github.com/NetCloth/netcloth-chain/types/rest"
	"github.com/gorilla/mux"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/aipal/list",
		listHandlerFn(cliCtx),
	).Methods("GET")
}

func listHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		resKVs, height, err := cliCtx.QuerySubspace(types.ServiceNodeByBondKey, types.StoreKey)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var serverNodes types.ServiceNodes
		if len(resKVs) > 0 {
			for i := len(resKVs) - 1; i >= 0; i-- {
				serverNodes = append(serverNodes, types.MustUnmarshalServerNodeObject(cliCtx.Codec, resKVs[i].Value))
			}
		}

		res, err := cliCtx.Codec.MarshalJSONIndent(serverNodes, "", "  ")

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
