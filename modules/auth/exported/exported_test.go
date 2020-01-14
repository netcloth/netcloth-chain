package exported

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/netcloth/netcloth-chain/modules/auth/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func TestGenesisAccountsContains(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	acc := types.NewBaseAccount(addr, nil, secp256k1.GenPrivKey().PubKey(), 0, 0)

	genAccounts := GenesisAccounts{}
	require.False(t, genAccounts.Contains(acc.GetAddress()))

	genAccounts = append(genAccounts, acc)
	require.True(t, genAccounts.Contains(acc.GetAddress()))
}
