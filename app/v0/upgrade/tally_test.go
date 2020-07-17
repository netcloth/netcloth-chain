package upgrade

import (
	"testing"

	"github.com/netcloth/netcloth-chain/app/v0/staking"
	sdk "github.com/netcloth/netcloth-chain/types"

	"github.com/stretchr/testify/require"
)

func TestTallyPassed(t *testing.T) {
	ctx, keeper, stakingKeeper, _ := CreateTestInput(t, 1000)
	description := staking.NewDescription("moniker1", "identity1", "website1", "details1")
	for i := 0; i < 4; i++ {
		validator := staking.NewValidator(sdk.ValAddress(Addrs[i]), PKs[i], description)
		validator.Status = sdk.Bonded
		validator.Tokens = sdk.TokensFromConsensusPower(1)
		stakingKeeper.SetValidator(ctx, validator)
		stakingKeeper.SetValidatorByPowerIndex(ctx, validator)
		keeper.SetSignal(ctx, 1, validator.GetConsAddr().String())
	}

	require.True(t, tally(ctx, 1, keeper, sdk.NewDecWithPrec(75, 2)))
}

func TestTallyNotPassed(t *testing.T) {
	ctx, keeper, stakingKeeper, _ := CreateTestInput(t, 1000)
	description := staking.NewDescription("moniker2", "identity2", "website2", "details2")
	for i := 0; i < 4; i++ {
		validator := staking.NewValidator(sdk.ValAddress(Addrs[i]), PKs[i], description)
		validator.Status = sdk.Bonded
		validator.Tokens = sdk.TokensFromConsensusPower(1)
		stakingKeeper.SetValidator(ctx, validator)
		stakingKeeper.SetValidatorByPowerIndex(ctx, validator)
		if i%2 == 0 {
			keeper.SetSignal(ctx, 1, validator.GetConsAddr().String())
		}
	}
	require.True(t, tally(ctx, 1, keeper, sdk.NewDecWithPrec(5, 2)))
	require.False(t, tally(ctx, 1, keeper, sdk.NewDecWithPrec(75, 2)))
}
