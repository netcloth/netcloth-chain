package keeper

import (
	"fmt"

	"github.com/netcloth/netcloth-chain/app/v0/supply/internal/types"
	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"
)

// SendCoinsFromModuleToAccount transfers coins from a ModuleAccount to an AccAddress
func (k Keeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string,
	recipientAddr sdk.AccAddress, amt sdk.Coins) error {

	senderAddr := k.GetModuleAddress(senderModule)
	if senderAddr == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	return k.bk.SendCoins(ctx, senderAddr, recipientAddr, amt)
}

// SendCoinsFromModuleToModule transfers coins from a ModuleAccount to another
func (k Keeper) SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error {

	senderAddr := k.GetModuleAddress(senderModule)
	if senderAddr == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	// create the account if it doesn't yet exist
	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	return k.bk.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

// SendCoinsFromAccountToModule transfers coins from an AccAddress to a ModuleAccount
func (k Keeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress,
	recipientModule string, amt sdk.Coins) error {

	// create the account if it doesn't yet exist
	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	return k.bk.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

// DelegateCoinsFromAccountToModule delegates coins and transfers
// them from a delegator account to a module account
func (k Keeper) DelegateCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress,
	recipientModule string, amt sdk.Coins) error {

	// create the account if it doesn't yet exist
	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	if !recipientAcc.HasPermission(types.Staking) {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to receive delegated coins", recipientModule))
	}

	return k.bk.DelegateCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

// UndelegateCoinsFromModuleToAccount undelegates the unbonding coins and transfers
// them from a module account to the delegator account
func (k Keeper) UndelegateCoinsFromModuleToAccount(ctx sdk.Context, senderModule string,
	recipientAddr sdk.AccAddress, amt sdk.Coins) error {

	acc := k.GetModuleAccount(ctx, senderModule)
	if acc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	if !acc.HasPermission(types.Staking) {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to undelegate coins", senderModule))
	}

	return k.bk.UndelegateCoins(ctx, acc.GetAddress(), recipientAddr, amt)
}

// MintCoins creates new coins from thin air and adds it to the module account.
// Panics if the name maps to a non-minter module account or if the amount is invalid.
func (k Keeper) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {

	// create the account if it doesn't yet exist
	acc := k.GetModuleAccount(ctx, moduleName)
	if acc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", moduleName))
	}

	if !acc.HasPermission(types.Minter) {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to mint tokens", moduleName))
	}

	_, err := k.bk.AddCoins(ctx, acc.GetAddress(), amt)
	if err != nil {
		panic(err)
	}

	// update total supply
	supply := k.GetSupply(ctx)
	supply = supply.Inflate(amt)

	k.SetSupply(ctx, supply)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("minted %s from %s module account", amt.String(), moduleName))

	return nil
}

// BurnCoins burns coins deletes coins from the balance of the module account.
// Panics if the name maps to a non-burner module account or if the amount is invalid.
func (k Keeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {

	// create the account if it doesn't yet exist
	acc := k.GetModuleAccount(ctx, moduleName)
	if acc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", moduleName))
	}

	if !acc.HasPermission(types.Burner) {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to burn tokens", moduleName))
	}

	_, err := k.bk.SubtractCoins(ctx, acc.GetAddress(), amt)
	if err != nil {
		panic(err)
	}

	// update total supply
	supply := k.GetSupply(ctx)
	supply = supply.Deflate(amt)
	k.SetSupply(ctx, supply)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("burned %s from %s module account", amt.String(), moduleName))

	return nil
}
