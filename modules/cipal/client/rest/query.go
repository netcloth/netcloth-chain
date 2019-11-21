package rest

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"

	"github.com/netcloth/netcloth-chain/client/context"
	"github.com/netcloth/netcloth-chain/modules/cipal/types"
	"github.com/netcloth/netcloth-chain/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/cipal/query/{accAddress}",
		CIPALFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/cipal/batch_query",
		CIPALsFn(cliCtx),
	).Methods("POST")
}

func queryCIPAL(cliCtx context.CLIContext, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addr := vars["accAddress"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryCIPALParams(addr)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(endpoint, bz)
		if err != nil && !strings.Contains(err.Error(), "not found") {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryCIPALs(cliCtx context.CLIContext, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params types.QueryCIPALsParams

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &params) {
			fmt.Fprint(os.Stderr, fmt.Sprintf("params = %v\n", params))
			return
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(endpoint, bz)
		if err != nil && !strings.Contains(err.Error(), "not found") {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func CIPALFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryCIPAL(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCIPAL))
}

func CIPALsFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryCIPALs(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCIPALs))
}
