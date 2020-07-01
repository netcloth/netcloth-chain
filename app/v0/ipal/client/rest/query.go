package rest

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"

	"github.com/netcloth/netcloth-chain/app/v0/ipal/types"
	"github.com/netcloth/netcloth-chain/client/context"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/ipal/list",
		listHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/ipal/node/{accAddr}",
		nodeHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/ipal/nodes",
		nodesHandlerFn(cliCtx),
	).Methods("POST") //TODO should be GET, but go-sdk use lib[github.com/parnurzeal/gorequest] which can not use GET method to send body
}

func listHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		resKVs, height, err := cliCtx.QuerySubspace(types.IPALNodeByBondKey, types.StoreKey)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var ipalNodes types.IPALNodes
		if len(resKVs) > 0 {
			for i := len(resKVs) - 1; i >= 0; i-- {
				ipalNodes = append(ipalNodes, types.MustUnmarshalIPALNode(cliCtx.Codec, resKVs[i].Value))
			}
		}

		res, err := cliCtx.Codec.MarshalJSONIndent(ipalNodes, "", "  ")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryNode(cliCtx context.CLIContext, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32accAddr := vars["accAddr"]

		accAddr, err := sdk.AccAddressFromBech32(bech32accAddr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryIPALNodeParams(accAddr)

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

func queryNodes(cliCtx context.CLIContext, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var params types.QueryIPALNodesParams

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &params) {
			fmt.Fprint(os.Stderr, fmt.Sprintf("params = %v\n", params))
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
		if err != nil && !strings.Contains(err.Error(), "not found") {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func nodeHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryNode(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryIPALNode))
}

func nodesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryNodes(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryIPALNodes))
}
