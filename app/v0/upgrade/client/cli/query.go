package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/netcloth/netcloth-chain/app/v0/staking"
	"github.com/netcloth/netcloth-chain/app/v0/staking/types"
	upgcli "github.com/netcloth/netcloth-chain/app/v0/upgrade/client"
	upgtypes "github.com/netcloth/netcloth-chain/app/v0/upgrade/types"
	"github.com/netcloth/netcloth-chain/client"
	"github.com/netcloth/netcloth-chain/client/context"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	flagDetail = "detail"
)

func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	ipalQueryCmd := &cobra.Command{
		Use:                        upgtypes.ModuleName,
		Short:                      "Querying commands for upgrade",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ipalQueryCmd.AddCommand(client.GetCommands(
		GetInfoCmd(queryRoute, cdc),
		GetCmdQuerySignals(queryRoute, cdc),
	)...)

	return ipalQueryCmd

}

func GetInfoCmd(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "info",
		Short:   "Query the information of upgrade module",
		Example: "nchcli query upgrade info",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			resCurrentversion, _, _ := cliCtx.QueryStore(sdk.CurrentVersionKey, sdk.MainStore)
			var currentVersion uint64
			cdc.MustUnmarshalBinaryLengthPrefixed(resCurrentversion, &currentVersion)

			resProposalid, _, _ := cliCtx.QueryStore(upgtypes.GetSuccessVersionKey(currentVersion), storeName)
			var proposalID uint64
			cdc.MustUnmarshalBinaryLengthPrefixed(resProposalid, &proposalID)

			resCurrentversioninfo, _, err := cliCtx.QueryStore(upgtypes.GetProposalIDKey(proposalID), storeName)
			var currentVersionInfo upgtypes.VersionInfo
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

			return cliCtx.PrintOutput(upgradeInfoOutput)
		},
	}
	return cmd
}

func GetCmdQuerySignals(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-signals",
		Short:   "Query the information of signals",
		Example: "nchcli query upgrade query-signals",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res_upgradeConfig, _, err := cliCtx.QueryStore(sdk.UpgradeConfigKey, sdk.MainStore)
			if err != nil {
				return err
			}
			if len(res_upgradeConfig) == 0 {
				fmt.Println("No Software Upgrade Switch Period is in process.")
				return err
			}

			var upgradeConfig sdk.UpgradeConfig
			if err = cdc.UnmarshalBinaryLengthPrefixed(res_upgradeConfig, &upgradeConfig); err != nil {
				return err
			}

			validatorConsAddrs := make(map[string]bool)
			res, _, err := cliCtx.QuerySubspace(upgtypes.GetSignalPrefixKey(upgradeConfig.Protocol.Version), storeName)
			if err != nil {
				return err
			}

			for _, kv := range res {
				validatorConsAddrs[upgtypes.GetAddressFromSignalKey(kv.Key)] = true
			}

			if len(validatorConsAddrs) == 0 {
				fmt.Println("No validator has started the new version.")
				return nil
			}

			key := staking.ValidatorsKey
			resKVs, _, err := cliCtx.QuerySubspace(key, "staking")
			if err != nil {
				return err
			}

			isDetail := viper.GetBool(flagDetail)
			totalVotingPower := sdk.ZeroDec()
			signalsVotingPower := sdk.ZeroDec()

			for _, kv := range resKVs {
				validator := types.MustUnmarshalValidator(cdc, kv.Value)
				power := sdk.NewDec(validator.GetConsensusPower())
				totalVotingPower = totalVotingPower.Add(power)
				if _, ok := validatorConsAddrs[validator.GetConsAddr().String()]; ok {
					signalsVotingPower = signalsVotingPower.Add(power)
					if isDetail {
						fmt.Println(validator.GetOperator().String(), " ", power)
					}
				}
			}
			fmt.Println("signalsVotingPower/totalVotingPower = " + signalsVotingPower.Quo(totalVotingPower).String())
			return nil
		},
	}
	cmd.Flags().Bool(flagDetail, false, "details of siganls")
	return cmd
}
