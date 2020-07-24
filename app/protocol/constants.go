package protocol

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

// all modules
const (
	ParamsModuleName       = "params"
	SupplyModuleName       = "supply"
	StakingModuleName      = "staking"
	MintModuleName         = "mint"
	DistributionModuleName = "distribution"
	SlashingModuleName     = "slashing"
	GovModuleName          = "gov"
	AuthModuleName         = "auth"
	UpgradeModuleName      = "upgrade"
	GuardianModuleName     = "guardian"
	IpalModuleName         = "ipal"
	CIpalModuleName        = "cipal"
	VMModuleName           = "vm"
)

// all store keys name
const (
	MainStoreKey = "main"

	ParamsStoreKey       = ParamsModuleName
	SupplyStoreKey       = SupplyModuleName
	StakingStoreKey      = StakingModuleName
	MintStoreKey         = MintModuleName
	DistributionStoreKey = DistributionModuleName
	SlashingStoreKey     = SlashingModuleName
	GovStoreKey          = GovModuleName
	AuthStoreKey         = AuthModuleName
	UpgradeStoreKey      = UpgradeModuleName
	GuardianStoreKey     = GuardianModuleName
	IpalStoreKey         = IpalModuleName
	CIpalStoreKey        = CIpalModuleName
	VMStoreKey           = VMModuleName

	ParamsTStoreKey  = "transient_" + ParamsStoreKey
	StakingTStoreKey = "transient_" + StakingStoreKey
)

// all store keys
var (
	Keys = sdk.NewKVStoreKeys(
		MainStoreKey,
		ParamsStoreKey,
		SupplyStoreKey,
		StakingStoreKey,
		MintStoreKey,
		DistributionStoreKey,
		SlashingStoreKey,
		CIpalStoreKey,
		IpalStoreKey,
		VMStoreKey,
		GovStoreKey,
		AuthStoreKey,
		UpgradeStoreKey,
		GuardianStoreKey,
	)

	TKeys = sdk.NewTransientStoreKeys(
		ParamsTStoreKey,
		StakingTStoreKey,
	)
)
