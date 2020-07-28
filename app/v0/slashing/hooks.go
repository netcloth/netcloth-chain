// nolint
package slashing

import (
	"time"

	"github.com/tendermint/tendermint/crypto"

	"github.com/netcloth/netcloth-chain/app/v0/slashing/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// AfterValidatorBonded - call hook if registered
func (k Keeper) AfterValidatorBonded(ctx sdk.Context, address sdk.ConsAddress, _ sdk.ValAddress) {
	// Update the signing info start height or create a new signing info
	_, found := k.GetValidatorSigningInfo(ctx, address)
	if !found {
		signingInfo := types.NewValidatorSigningInfo(
			address,
			ctx.BlockHeight(),
			0,
			time.Unix(0, 0),
			false,
			0,
		)
		k.SetValidatorSigningInfo(ctx, address, signingInfo)
	}
}

// AfterValidatorCreated - When a validator is created, add the address-pubkey relation.
func (k Keeper) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	validator := k.sk.Validator(ctx, valAddr)
	k.addPubkey(ctx, validator.GetConsPubKey())
}

// AfterValidatorRemoved - When a validator is removed, delete the address-pubkey relation.
func (k Keeper) AfterValidatorRemoved(ctx sdk.Context, address sdk.ConsAddress) {
	k.deleteAddrPubkeyRelation(ctx, crypto.Address(address))
}

//_________________________________________________________________________________________

// Hooks wrapper struct for slashing keeper
type Hooks struct {
	k Keeper
}

var _ types.StakingHooks = Hooks{}

// Hooks return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// AfterValidatorBonded - Implements sdk.ValidatorHooks
func (h Hooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.k.AfterValidatorBonded(ctx, consAddr, valAddr)
}

// AfterValidatorRemoved - Implements sdk.ValidatorHooks
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, _ sdk.ValAddress) {
	h.k.AfterValidatorRemoved(ctx, consAddr)
}

// AfterValidatorCreated - Implements sdk.ValidatorHooks
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	h.k.AfterValidatorCreated(ctx, valAddr)
}

// AfterValidatorBeginUnbonding - unused hooks
func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) {}

// BeforeValidatorModified - unused hooks
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress) {}

// BeforeDelegationCreated - unused hooks
func (h Hooks) BeforeDelegationCreated(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {}

// BeforeDelegationSharesModified - unused hooks
func (h Hooks) BeforeDelegationSharesModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {}

// BeforeDelegationRemoved - unused hooks
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {}

// AfterDelegationModified - unused hooks
func (h Hooks) AfterDelegationModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {}

// BeforeValidatorSlashed - unused hooks
func (h Hooks) BeforeValidatorSlashed(_ sdk.Context, _ sdk.ValAddress, _ sdk.Dec) {}
