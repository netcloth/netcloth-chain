package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"
	"time"

	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/netcloth/netcloth-chain/app/v0/staking/types"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/types/module"
	"github.com/netcloth/netcloth-chain/types/simulation"
)

// Simulation parameter constants
const (
	unbondingTime     = "unbonding_time"
	maxValidators     = "max_validators"
	historicalEntries = "historical_entries"
)

// GenUnbondingTime randomized UnbondingTime
func GenUnbondingTime(r *rand.Rand) (ubdTime time.Duration) {
	return time.Duration(simulation.RandIntBetween(r, 60, 60*60*24*3*2)) * time.Second
}

// GenMaxValidators randomized MaxValidators
func GenMaxValidators(r *rand.Rand) (maxValidators uint32) {
	return uint32(r.Intn(250) + 1)
}

// GetHistEntries randomized HistoricalEntries between 0-100.
func GetHistEntries(r *rand.Rand) uint32 {
	return uint32(r.Intn(101))
}

// RandomizedGenState generates a random GenesisState for staking
func RandomizedGenState(simState *module.SimulationState) {
	// params
	var (
		unbondTime  time.Duration
		maxVals     uint32
		histEntries uint32
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, unbondingTime, &unbondTime, simState.Rand,
		func(r *rand.Rand) { unbondTime = GenUnbondingTime(r) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, maxValidators, &maxVals, simState.Rand,
		func(r *rand.Rand) { maxVals = GenMaxValidators(r) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, historicalEntries, &histEntries, simState.Rand,
		func(r *rand.Rand) { histEntries = GetHistEntries(r) },
	)

	// NOTE: the slashing module need to be defined after the staking module on the
	// NewSimulationManager constructor for this to work
	simState.UnbondTime = unbondTime
	var MaxValidatorsExtendingInterval = time.Duration(60) * 60 * 8766
	params := types.NewParams(
		simState.UnbondTime,
		uint16(maxVals),
		700,
		1,
		tmtime.Now().Add(time.Second*MaxValidatorsExtendingInterval),
		100,
		sdk.DefaultBondDenom,
		sdk.NewDec(20),
	)

	// validators & delegations
	var (
		validators  []types.Validator
		delegations []types.Delegation
	)

	valAddrs := make([]sdk.ValAddress, simState.NumBonded)

	for i := 0; i < int(simState.NumBonded); i++ {
		valAddr := sdk.ValAddress(simState.Accounts[i].Address)
		valAddrs[i] = valAddr

		maxCommission := sdk.NewDecWithPrec(int64(simulation.RandIntBetween(simState.Rand, 1, 100)), 2)
		commission := types.NewCommission(
			simulation.RandomDecAmount(simState.Rand, maxCommission),
			maxCommission,
			simulation.RandomDecAmount(simState.Rand, maxCommission),
		)

		validator := types.NewValidator(valAddr, simState.Accounts[i].PubKey, types.Description{})
		validator.Tokens = sdk.NewInt(simState.InitialStake)
		validator.DelegatorShares = sdk.NewDec(simState.InitialStake)
		validator.Commission = commission
		validator.Status = sdk.Bonded
		validator.SelfDelegation = sdk.NewDec(simState.InitialStake)

		delegation := types.NewDelegation(simState.Accounts[i].Address, valAddr, sdk.NewDec(simState.InitialStake))

		validators = append(validators, validator)
		delegations = append(delegations, delegation)
	}

	stakingGenesis := types.NewGenesisState(params, validators, delegations)

	fmt.Println(string(simState.GenState[types.ModuleName]))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(stakingGenesis)
}
