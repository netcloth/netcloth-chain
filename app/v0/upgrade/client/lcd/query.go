package lcd

import (
	"github.com/netcloth/netcloth-chain/app/v0/upgrade/types"
	"net/http"

	upgcli "github.com/netcloth/netcloth-chain/app/v0/upgrade/client"
	"github.com/netcloth/netcloth-chain/client/context"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
)

type VersionInfo struct {
	Version        string `json:"version"`
	UpgradeVersion int64  `json:"upgrade_version"`
	StartHeight    int64  `json:"start_height"`
	ProposalId     int64  `json:"proposal_id"`
}

func InfoHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		resCurrentversion, _, _ := cliCtx.QueryStore(sdk.CurrentVersionKey, sdk.MainStore)
		var currentVersion uint64
		cdc.MustUnmarshalBinaryLengthPrefixed(resCurrentversion, &currentVersion)

		resProposalid, _, _ := cliCtx.QueryStore(types.GetSuccessVersionKey(currentVersion), storeName)
		var proposalID uint64
		cdc.MustUnmarshalBinaryLengthPrefixed(resProposalid, &proposalID)

		resCurrentversioninfo, _, err := cliCtx.QueryStore(types.GetProposalIDKey(proposalID), storeName)
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
