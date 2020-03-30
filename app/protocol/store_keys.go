package protocol

import (
	"github.com/netcloth/netcloth-chain/modules/cipal"
	distr "github.com/netcloth/netcloth-chain/modules/distribution"
	"github.com/netcloth/netcloth-chain/modules/gov"
	"github.com/netcloth/netcloth-chain/modules/ipal"
	"github.com/netcloth/netcloth-chain/modules/mint"
	"github.com/netcloth/netcloth-chain/modules/params"
	"github.com/netcloth/netcloth-chain/modules/slashing"
	"github.com/netcloth/netcloth-chain/modules/staking"
	"github.com/netcloth/netcloth-chain/modules/supply"
	"github.com/netcloth/netcloth-chain/modules/vm"
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	paramsStoreKey   = params.StoreKey
	supplyStoreKey   = supply.StoreKey
	stakingStoreKey  = staking.StoreKey
	mintStoreKey     = mint.StoreKey
	distrStoreKey    = distr.StoreKey
	slashingStoreKey = slashing.StoreKey
	cipalStoreKey    = cipal.StoreKey
	ipalStoreKey     = ipal.StoreKey
	vmStoreKey       = vm.StoreKey
	vmCodeKey        = vm.CodeKey
	vmStoreDebugKey  = vm.StoreDebugKey
	govStoreKey      = gov.StoreKey

	paramsTStoreKey  = params.TStoreKey
	stakingTStoreKey = staking.TStoreKey
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
