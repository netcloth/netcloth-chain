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

	r.HandleFunc(
		fmt.Sprintf("/vm/%s", types.EstimateGas),
		estimateGasFn(cliCtx),
	).Methods("POST")
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

func estimateGas(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var params types.MsgContract
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &params) {
			return
		}

		if params.From == nil || params.Payload == nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request"))
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		d, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			return
		}

		route := fmt.Sprintf("custom/vm/%s", types.EstimateGas)
		res, height, err := cliCtx.QueryWithData(route, d)
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

func estimateGasFn(cliCtx context.CLIContext) http.HandlerFunc {
	return estimateGas(cliCtx)
}
