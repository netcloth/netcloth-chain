package protocol

import (
	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"github.com/netcloth/netcloth-chain/app/v0/cipal"
	"github.com/netcloth/netcloth-chain/app/v0/distribution"
	"github.com/netcloth/netcloth-chain/app/v0/gov"
	"github.com/netcloth/netcloth-chain/app/v0/guardian"
	"github.com/netcloth/netcloth-chain/app/v0/ipal"
	"github.com/netcloth/netcloth-chain/app/v0/mint"
	"github.com/netcloth/netcloth-chain/app/v0/params"
	"github.com/netcloth/netcloth-chain/app/v0/slashing"
	"github.com/netcloth/netcloth-chain/app/v0/staking"
	"github.com/netcloth/netcloth-chain/app/v0/supply"
	"github.com/netcloth/netcloth-chain/app/v0/upgrade"
	"github.com/netcloth/netcloth-chain/app/v0/vm"
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	MainStoreKey     = "main"
	paramsStoreKey   = params.StoreKey
	supplyStoreKey   = supply.StoreKey
	StakingStoreKey  = staking.StoreKey
	mintStoreKey     = mint.StoreKey
	distrStoreKey    = distribution.StoreKey
	slashingStoreKey = slashing.StoreKey
	cipalStoreKey    = cipal.StoreKey
	ipalStoreKey     = ipal.StoreKey
	VMStoreKey       = vm.StoreKey
	VMCodeKey        = vm.CodeKey
	VMStoreDebugKey  = vm.StoreDebugKey
	govStoreKey      = gov.StoreKey
	authStoreKey     = auth.StoreKey
	UpgradeStoreKey  = upgrade.StoreKey
	GuardianStoreKey = guardian.StoreKey

	paramsTStoreKey  = params.TStoreKey
	stakingTStoreKey = staking.TStoreKey
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
		GuardianStoreKey,
	)

	TKeys = sdk.NewTransientStoreKeys(
		paramsTStoreKey,
		stakingTStoreKey,
	)
)
