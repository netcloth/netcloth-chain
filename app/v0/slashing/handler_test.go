package slashing

import (
	"testing"

	"github.com/netcloth/netcloth-chain/app/v0/staking"
	sdk "github.com/netcloth/netcloth-chain/types"

	"github.com/stretchr/testify/require"
)

func TestCannotUnjailUnlessJailed(t *testing.T) {
	// initial setup
	ctx, ck, sk, _, keeper := createTestInput(t, DefaultParams())
	slashHandler := NewHandler(keeper)
	amt := sdk.TokensFromConsensusPower(100)
	addr, val := addrs[0], pks[0]
	msg := NewTestMsgCreateValidator(addr, val, amt)
	_, err := staking.NewHandler(sk)(ctx, msg)
	require.Nil(t, err, "%v", err)
	staking.EndBlocker(ctx, sk)

	require.Equal(
		t, ck.GetCoins(ctx, sdk.AccAddress(addr)),
		sdk.Coins{sdk.NewCoin(sk.GetParams(ctx).BondDenom, initTokens.Sub(amt))},
	)
	require.Equal(t, amt, sk.Validator(ctx, addr).GetBondedTokens())

	// assert non-jailed validator can't be unjailed
	_, err = slashHandler(ctx, NewMsgUnjail(addr))
	require.NotNil(t, err, "allowed unjail of non-jailed validator")
}

func TestCannotUnjailWithMaxLever(t *testing.T) {
	// TODO
}
