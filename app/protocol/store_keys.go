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
	VMStoreKey       = "vm" //TODO merge keys here and modules key
	VMCodeKey        = "vm_code"
	VMStoreDebugKey  = "vm_decode"
	govStoreKey      = "gov"
	authStoreKey     = "acc"
	UpgradeStoreKey  = "upgrade"

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
		VMStoreKey,
		VMCodeKey,
		VMStoreDebugKey,
		govStoreKey,
		authStoreKey,
		UpgradeStoreKey,
	)

	TKeys = sdk.NewTransientStoreKeys(
		paramsTStoreKey,
		stakingTStoreKey,
	)
)
