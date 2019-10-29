package rest

import (
    "net/http"
    "github.com/gorilla/mux"
    "github.com/NetCloth/netcloth-chain/client/context"
    "github.com/NetCloth/netcloth-chain/modules/aipal/types"
    "github.com/NetCloth/netcloth-chain/types/rest"
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
        for _, kv := range resKVs {
            serverNodes = append(serverNodes, types.MustUnmarshalServerNodeObject(cliCtx.Codec, kv.Value))
        }

        res, err := cliCtx.Codec.MarshalJSONIndent(serverNodes, "", "  ")

        cliCtx = cliCtx.WithHeight(height)
        rest.PostProcessResponse(w, cliCtx, res)
    }
}
