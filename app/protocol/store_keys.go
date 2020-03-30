package protocol

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	paramsStoreKey   = "params"
	supplyStoreKey   = "supply"
	stakingStoreKey  = "staking"
	mintStoreKey     = "mint"
	distrStoreKey    = "distribution"
	slashingStoreKey = "slashing"
	cipalStoreKey    = "cipal"
	ipalStoreKey     = "ipal"
	vmStoreKey       = "vm"
	vmCodeKey        = "vm_code"
	vmStoreDebugKey  = "vm_decode"
	govStoreKey      = "gov"

	paramsTStoreKey  = "transient_params"
	stakingTStoreKey = "transient_staking"
)

var (
	Keys = sdk.NewKVStoreKeys(
		paramsStoreKey,
		supplyStoreKey,
		stakingStoreKey,
		mintStoreKey,
		distrStoreKey,
		slashingStoreKey,
		cipalStoreKey,
		ipalStoreKey,
		vmStoreKey,
		vmCodeKey,
		vmStoreDebugKey,
		govStoreKey,
	)

	TKeys = sdk.NewTransientStoreKeys(
		paramsTStoreKey,
		stakingTStoreKey,
	)
)
