package genaccounts

import (
	"testing"

	"github.com/stretchr/testify/require"

	authexported "github.com/netcloth/netcloth-chain/app/v0/auth/exported"
	authtypes "github.com/netcloth/netcloth-chain/app/v0/auth/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func TestGenesis(t *testing.T) {
	config := setupTestInput()
	genCoin := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000)))

	_, pubKey, addr := KeyTestPubAddr()

	genesisAccount := NewGenesisAccount(
		authtypes.NewBaseAccount(addr, genCoin, pubKey, 0, 1))
	gaccs := make(GenesisState, 0)
	gaccs = append(gaccs, genesisAccount)
	InitGenesis(config.ctx, ModuleCdc, config.ak, gaccs)

	genesisState := ExportGenesis(config.ctx, config.ak)

	require.Equal(t, gaccs, genesisState)

	require.NotNil(t, genesisState)
	require.Len(t, genesisState, 1)

	config.ak.IterateAccounts(config.ctx,
		func(account authexported.Account) (stop bool) {
			require.Equal(t, account, genesisAccount.ToAccount())
			return false
		})
}
