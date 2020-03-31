package protocol

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	MainStoreKey     = "main"
	paramsStoreKey   = "params"
	supplyStoreKey   = "supply"
	StakingStoreKey  = "staking"
	mintStoreKey     = "mint"
	distrStoreKey    = "distribution"
	slashingStoreKey = "slashing"
	cipalStoreKey    = "cipal"
	ipalStoreKey     = "ipal"
	vmStoreKey       = "vm"
	vmCodeKey        = "vm_code"
	vmStoreDebugKey  = "vm_decode"
	govStoreKey      = "gov"
	authStoreKey     = "acc"

	paramsTStoreKey  = "transient_params"
	stakingTStoreKey = "transient_staking"
)

var (
	MainKVStoreKey = sdk.NewKVStoreKey(MainStoreKey)

	Keys = sdk.NewKVStoreKeys(
		paramsStoreKey,
		supplyStoreKey,
		StakingStoreKey,
		mintStoreKey,
		distrStoreKey,
		slashingStoreKey,
		cipalStoreKey,
		ipalStoreKey,
		vmStoreKey,
		vmCodeKey,
		vmStoreDebugKey,
		govStoreKey,
		authStoreKey,
	)

	TKeys = sdk.NewTransientStoreKeys(
		paramsTStoreKey,
		stakingTStoreKey,
	)
)
