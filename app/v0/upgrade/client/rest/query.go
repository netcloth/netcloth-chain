package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	upgcli "github.com/netcloth/netcloth-chain/app/v0/upgrade/client"
	"github.com/netcloth/netcloth-chain/app/v0/upgrade/types"
	"github.com/netcloth/netcloth-chain/client/context"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/upgrade/info",
		InfoHandlerFn(cliCtx),
	).Methods("GET")
}

// VersionInfo is the struct of version info
type VersionInfo struct {
	Version        string `json:"version"`
	UpgradeVersion int64  `json:"upgrade_version"`
	StartHeight    int64  `json:"start_height"`
	ProposalID     int64  `json:"proposal_id"`
}

// InfoHandlerFn - HTTP request handler to query the upgrade info
func InfoHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cdc := cliCtx.Codec
		resCurrentversion, _, _ := cliCtx.QueryStore(sdk.CurrentVersionKey, sdk.MainStore)
		var currentVersion uint64
		cdc.MustUnmarshalBinaryLengthPrefixed(resCurrentversion, &currentVersion)

		resProposalid, _, _ := cliCtx.QueryStore(types.GetSuccessVersionKey(currentVersion), "upgrade")
		var proposalID uint64
		cdc.MustUnmarshalBinaryLengthPrefixed(resProposalid, &proposalID)

		resCurrentversioninfo, _, err := cliCtx.QueryStore(types.GetProposalIDKey(proposalID), "upgrade")
		var currentVersionInfo types.VersionInfo
		cdc.MustUnmarshalBinaryLengthPrefixed(resCurrentversioninfo, &currentVersionInfo)

		resUpgradeinprogress, _, _ := cliCtx.QueryStore(sdk.UpgradeConfigKey, sdk.MainStore)
		var upgradeInProgress sdk.UpgradeConfig
		if err == nil && len(resUpgradeinprogress) != 0 {
			cdc.MustUnmarshalBinaryLengthPrefixed(resUpgradeinprogress, &upgradeInProgress)
		}

		resLastfailedversion, _, err := cliCtx.QueryStore(sdk.LastFailedVersionKey, sdk.MainStore)
		var lastFailedVersion uint64
		if err == nil && len(resLastfailedversion) != 0 {
			cdc.MustUnmarshalBinaryLengthPrefixed(resLastfailedversion, &lastFailedVersion)
		} else {
			lastFailedVersion = 0
		}

		upgradeInfoOutput := upgcli.NewUpgradeInfoOutput(currentVersionInfo, lastFailedVersion, upgradeInProgress)

		output, err := cdc.MarshalJSONIndent(upgradeInfoOutput, "", "  ")
		if err != nil {
			//utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Write(output)
	}
}
