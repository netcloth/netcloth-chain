package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/netcloth/netcloth-chain/app/v0/auth/client/utils"
	"github.com/netcloth/netcloth-chain/app/v0/auth/types"
	genutilrest "github.com/netcloth/netcloth-chain/app/v0/genutil/client/rest"
	"github.com/netcloth/netcloth-chain/client/context"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/types/rest"
)

// query accountREST Handler
func QueryAccountRequestHandlerFn(storeName string, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		accGetter := types.NewAccountRetriever(cliCtx)

		account, height, err := accGetter.GetAccountWithHeight(addr)
		if err != nil {
			// TODO: Handle more appropriately based on the error type.
			if err := accGetter.EnsureExists(addr); err != nil {
				cliCtx = cliCtx.WithHeight(height)
				rest.PostProcessResponse(w, cliCtx, types.BaseAccount{})
				return
			}

			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, account)
	}
}

// QueryTxsHandlerFn implements a REST handler that searches for transactions.
// Genesis transactions are returned if the height parameter is set to zero,
// otherwise the transactions are searched for by events.
func QueryTxsRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			rest.WriteErrorResponse(
				w, http.StatusBadRequest,
				fmt.Sprintf("failed to parse query parameters: %s", err),
			)
			return
		}

		// if the height query param is set to zero, query for genesis transactions
		heightStr := r.FormValue("height")
		if heightStr != "" {
			if height, err := strconv.ParseInt(heightStr, 10, 64); err == nil && height == 0 {
				genutilrest.QueryGenesisTxs(cliCtx, w)
				return
			}
		}

		var (
			events      []string
			txs         []sdk.TxResponse
			page, limit int
		)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		if len(r.Form) == 0 {
			rest.PostProcessResponseBare(w, cliCtx, txs)
			return
		}

		events, page, limit, err = rest.ParseHTTPArgs(r)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		searchResult, err := utils.QueryTxsByEvents(cliCtx, events, page, limit)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, searchResult)
	}
}

// QueryTxRequestHandlerFn implements a REST handler that queries a transaction
// by hash in a committed block.
func QueryTxRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hashHexStr := vars["hash"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		output, err := utils.QueryTx(cliCtx, hashHexStr)
		if err != nil {
			if strings.Contains(err.Error(), hashHexStr) {
				rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if output.Empty() {
			rest.WriteErrorResponse(w, http.StatusNotFound, fmt.Sprintf("no transaction found with hash %s", hashHexStr))
		}

		rest.PostProcessResponseBare(w, cliCtx, output)
	}
}

func queryParamsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParams)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func EstimateGas(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BroadcastReq

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(req.Tx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData("/app/simulate", txBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var gas uint64
		cliCtx.Codec.MustUnmarshalBinaryLengthPrefixed(res, &gas)

		type GasRes struct {
			Gas uint64
			Res string
		}
		type GasResult struct {
			Height int64  `json:"height"`
			Result GasRes `json:"result"`
		}

		gr := GasResult{Height: height}
		gr.Result.Gas = gas
		resp, _ := json.MarshalIndent(gr, "", "  ")
		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}
