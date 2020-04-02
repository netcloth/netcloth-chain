package upgrade

import (
	"github.com/netcloth/netcloth-chain/app/v0/staking/exported"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func tally(ctx sdk.Context, versionProtocol uint64, k Keeper, threshold sdk.Dec) (passes bool) {
	totalVotingPower := sdk.ZeroDec()
	signalsVotingPower := sdk.ZeroDec()

	k.sk.IterateBondedValidatorsByPower(ctx, func(index int64, validator exported.ValidatorI) (stop bool) {
		power := sdk.NewDec(validator.GetConsensusPower())
		totalVotingPower = totalVotingPower.Add(power)
		valAcc := validator.GetConsAddr().String()
		if ok := k.GetSignal(ctx, versionProtocol, valAcc); ok {
			signalsVotingPower = signalsVotingPower.Add(power)
		}
		return false
	})

	ctx.Logger().Info("Tally Start", "SiganlsVotingPower", signalsVotingPower.String(),
		"TotalVotingPower", totalVotingPower.String(),
		"SiganlsVotingPower/TotalVotingPower", signalsVotingPower.Quo(totalVotingPower).String(),
		"Threshold", threshold.String())

	if signalsVotingPower.Quo(totalVotingPower).GT(threshold) {
		return true
	}
	return false
}
