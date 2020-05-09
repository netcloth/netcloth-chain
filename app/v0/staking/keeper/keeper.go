package keeper

import (
	"container/list"
	"fmt"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/netcloth/netcloth-chain/app/v0/params"
	"github.com/netcloth/netcloth-chain/app/v0/staking/types"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
)

const aminoCacheSize = 500

// Implements ValidatorSet interface
var _ types.ValidatorSet = Keeper{}

// Implements DelegationSet interface
var _ types.DelegationSet = Keeper{}

// keeper of the staking store
type Keeper struct {
	storeKey           sdk.StoreKey
	storeTKey          sdk.StoreKey
	cdc                *codec.Codec
	supplyKeeper       types.SupplyKeeper
	hooks              types.StakingHooks
	paramstore         params.Subspace
	validatorCache     map[string]cachedValidator
	validatorCacheList *list.List
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(cdc *codec.Codec, key, tkey sdk.StoreKey, supplyKeeper types.SupplyKeeper,
	paramstore params.Subspace) Keeper {

	// ensure bonded and not bonded module accounts are set
	if addr := supplyKeeper.GetModuleAddress(types.BondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	if addr := supplyKeeper.GetModuleAddress(types.NotBondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}

	return Keeper{
		storeKey:           key,
		storeTKey:          tkey,
		cdc:                cdc,
		supplyKeeper:       supplyKeeper,
		paramstore:         paramstore.WithKeyTable(ParamKeyTable()),
		hooks:              nil,
		validatorCache:     make(map[string]cachedValidator, aminoCacheSize),
		validatorCacheList: list.New(),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("modules/%s", types.ModuleName))
}

// Set the validator hooks
func (k *Keeper) SetHooks(sh types.StakingHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set validator hooks twice")
	}
	k.hooks = sh
	return k
}

// Load the last total validator power.
func (k Keeper) GetLastTotalPower(ctx sdk.Context) (power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.LastTotalPowerKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &power)
	return
}

// Set the last total validator power.
func (k Keeper) SetLastTotalPower(ctx sdk.Context, power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(power)
	store.Set(types.LastTotalPowerKey, b)
}

// Calculate the ValidatorUpdates for the current block
// Called in each EndBlock
func (k Keeper) BlockValidatorUpdates(ctx sdk.Context) []abci.ValidatorUpdate {
	// Calculate validator set changes.
	//
	// NOTE: ApplyAndReturnValidatorSetUpdates has to come before
	// UnbondAllMatureValidatorQueue.
	// This fixes a bug when the unbonding period is instant (is the case in
	// some of the tests). The test expected the validator to be completely
	// unbonded after the Endblocker (go from Bonded -> Unbonding during
	// ApplyAndReturnValidatorSetUpdates and then Unbonding -> Unbonded during
	// UnbondAllMatureValidatorQueue).
	validatorUpdates := k.ApplyAndReturnValidatorSetUpdates(ctx)

	// Unbond all mature validators from the unbonding queue.
	k.UnbondAllMatureValidatorQueue(ctx)

	// Remove all mature unbonding delegations from the ubd queue.
	matureUnbonds := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, dvPair := range matureUnbonds {
		err := k.CompleteUnbonding(ctx, dvPair.DelegatorAddress, dvPair.ValidatorAddress)
		if err != nil {
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnbonding,
				sdk.NewAttribute(types.AttributeKeyValidator, dvPair.ValidatorAddress.String()),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvPair.DelegatorAddress.String()),
			),
		)
	}

	// Remove all mature redelegations from the red queue.
	matureRedelegations := k.DequeueAllMatureRedelegationQueue(ctx, ctx.BlockHeader().Time)
	for _, dvvTriplet := range matureRedelegations {
		err := k.CompleteRedelegation(ctx, dvvTriplet.DelegatorAddress,
			dvvTriplet.ValidatorSrcAddress, dvvTriplet.ValidatorDstAddress)
		if err != nil {
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteRedelegation,
				sdk.NewAttribute(types.AttributeKeyDelegator, dvvTriplet.DelegatorAddress.String()),
				sdk.NewAttribute(types.AttributeKeySrcValidator, dvvTriplet.ValidatorSrcAddress.String()),
				sdk.NewAttribute(types.AttributeKeyDstValidator, dvvTriplet.ValidatorDstAddress.String()),
			),
		)
	}

	// update max validators
	k.UpdateMaxValidators(ctx)

	return validatorUpdates
}

func (k Keeper) UpdateMaxValidators(ctx sdk.Context) {
	p := k.GetParams(ctx)
	logger := k.Logger(ctx)

	// time not up
	if ctx.BlockTime().Before(p.NextExtendingTime) {
		return
	}

	// time up, update max validators
	if p.MaxValidatorsExtendingLimit > p.MaxValidators {
		logger.Info(fmt.Sprintf("Time up! Update maxValidators at blockTime: %v", ctx.BlockTime()))

		p.MaxValidators += p.MaxValidatorsExtendingSpeed
		if p.MaxValidators > p.MaxValidatorsExtendingLimit {
			p.MaxValidators = p.MaxValidatorsExtendingLimit
		}

		// update nextExtending time
		p.NextExtendingTime = p.NextExtendingTime.Add(time.Second * types.MaxValidatorsExtendingInterval)
		k.SetParams(ctx, p)
	} else {
		logger.Info(fmt.Sprintf("MaxValidators has reached upper limit %v at blockTime: %v", p.MaxValidatorsExtendingLimit, ctx.BlockTime()))
	}
}
