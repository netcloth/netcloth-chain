package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/netcloth/netcloth-chain/app/v0/cipal/types"
	"github.com/netcloth/netcloth-chain/client/context"
	"github.com/netcloth/netcloth-chain/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/cipal/query/{accAddress}",
		CIPALFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/cipal/count",
		CIPALCountFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/cipal/batch_query",
		CIPALsFn(cliCtx),
	).Methods("POST")
}

func queryCIPAL(cliCtx context.CLIContext, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := types.NewQueryCIPALParams(mux.Vars(r)["accAddress"])
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
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

func queryCIPALCount(cliCtx context.CLIContext, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params types.QueryCIPALsParams
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &params) {
			return
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(endpoint, bz)
		if err != nil {
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
			return
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(endpoint, bz)
		if err != nil {
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

func CIPALCountFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryCIPAL(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCIPALCount))
}

func CIPALsFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryCIPALs(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCIPALs))
}
