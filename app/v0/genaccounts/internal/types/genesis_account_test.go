package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/netcloth/netcloth-chain/types"
)

func TestGenesisAccountsContains(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	acc := GenesisAccount{Address: addr}

	genAccounts := GenesisAccounts{}
	require.False(t, genAccounts.Contains(acc.Address))

	genAccounts = append(genAccounts, acc)
	require.True(t, genAccounts.Contains(acc.Address))
}
