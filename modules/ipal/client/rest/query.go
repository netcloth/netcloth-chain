package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/NetCloth/netcloth-chain/client/context"
	"github.com/NetCloth/netcloth-chain/modules/ipal/types"
	"github.com/NetCloth/netcloth-chain/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/ipal/ipal/{accAddress}",
		IPALFn(cliCtx),
	).Methods("GET")
}

func queryIPAL(cliCtx context.CLIContext, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addr := vars["accAddress"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryIPALParams(addr)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
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

func IPALFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryIPAL(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryIPAL))
}
