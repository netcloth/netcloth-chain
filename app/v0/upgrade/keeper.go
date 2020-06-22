package upgrade

import (
	"github.com/netcloth/netcloth-chain/app/v0/staking"
	"github.com/netcloth/netcloth-chain/app/v0/staking/exported"
	"github.com/netcloth/netcloth-chain/app/v0/upgrade/types"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
)

type Keeper struct {
	storeKey       sdk.StoreKey
	cdc            *codec.Codec
	protocolKeeper sdk.ProtocolKeeper
	sk             staking.Keeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, protocolKeeper sdk.ProtocolKeeper, sk staking.Keeper) Keeper {
	keeper := Keeper{
		key,
		cdc,
		protocolKeeper,
		sk,
	}
	return keeper
}

func (k Keeper) AddNewVersionInfo(ctx sdk.Context, versionInfo types.VersionInfo) {
	kvStore := ctx.KVStore(k.storeKey)

	versionInfoBytes, err := k.cdc.MarshalBinaryLengthPrefixed(versionInfo)
	if err != nil {
		panic(err)
	}
	kvStore.Set(types.GetProposalIDKey(versionInfo.UpgradeInfo.ProposalID), versionInfoBytes)

	proposalIDBytes, err := k.cdc.MarshalBinaryLengthPrefixed(versionInfo.UpgradeInfo.ProposalID)
	if err != nil {
		panic(err)
	}

	if versionInfo.Success {
		kvStore.Set(types.GetSuccessVersionKey(versionInfo.UpgradeInfo.Protocol.Version), proposalIDBytes)
	} else {
		kvStore.Set(types.GetFailedVersionKey(versionInfo.UpgradeInfo.Protocol.Version, versionInfo.UpgradeInfo.ProposalID), proposalIDBytes)
	}
}

func (k Keeper) SetSignal(ctx sdk.Context, protocol uint64, address string) {
	kvStore := ctx.KVStore(k.storeKey)
	cmsgBytes, err := k.cdc.MarshalBinaryLengthPrefixed(true)
	if err != nil {
		panic(err)
	}
	kvStore.Set(types.GetSignalKey(protocol, address), cmsgBytes)
}

func (k Keeper) GetSignal(ctx sdk.Context, protocol uint64, address string) bool {
	kvStore := ctx.KVStore(k.storeKey)
	flagBytes := kvStore.Get(types.GetSignalKey(protocol, address))
	if flagBytes != nil {
		var flag bool
		err := k.cdc.UnmarshalBinaryLengthPrefixed(flagBytes, &flag)
		if err != nil {
			panic(err)
		}
		return true
	}
	return false
}

func (k Keeper) DeleteSignal(ctx sdk.Context, protocol uint64, address string) bool {
	if ok := k.GetSignal(ctx, protocol, address); ok {
		kvStore := ctx.KVStore(k.storeKey)
		kvStore.Delete(types.GetSignalKey(protocol, address))
		return true
	}
	return false
}

// IterateBondedValidatorsByPower iterates bonded validators by power
func (k Keeper) IterateBondedValidatorsByPower(ctx sdk.Context,
	fn func(index int64, validator exported.ValidatorI) (stop bool)) {
	k.sk.IterateBondedValidatorsByPower(ctx, fn)
}
