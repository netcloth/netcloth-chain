//+build ledger test_ledger_mock

package keys

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/netcloth/netcloth-chain/client/flags"
	"github.com/netcloth/netcloth-chain/crypto/keys"
	"github.com/netcloth/netcloth-chain/tests"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// TODO: fix this testcase
func Test_runAddCmdLedger(t *testing.T) {
	cmd := addKeyCommand()
	assert.NotNil(t, cmd)

	// Prepare a keybase
	kbHome, kbCleanUp := tests.NewTestCaseDir(t)
	assert.NotNil(t, kbHome)
	defer kbCleanUp()
	viper.Set(flags.FlagHome, kbHome)
	viper.Set(flags.FlagUseLedger, true)

	/// Test Text
	viper.Set(cli.OutputFlag, OutputFormatText)
	// Now enter password
	mockIn, _, _ := tests.ApplyMockIO(cmd)
	mockIn.Reset("test1234\ntest1234\n")
	assert.NoError(t, runAddCmd(cmd, []string{"keyname1"}))

	// Now check that it has been stored properly
	kb, err := NewKeyBaseFromHomeFlag()
	assert.NoError(t, err)
	assert.NotNil(t, kb)
	key1, err := kb.Get("keyname1")
	assert.NoError(t, err)
	assert.NotNil(t, key1)

	assert.Equal(t, "keyname1", key1.GetName())
	assert.Equal(t, keys.TypeLedger, key1.GetType())
	assert.Equal(t,
		"nchpub1addwnpepqd87l8xhcnrrtzxnkql7k55ph8fr9jarf4hn6udwukfprlalu8lgw0urza0",
		sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, key1.GetPubKey()))
}
