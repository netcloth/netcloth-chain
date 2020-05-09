package sim

import (
	"fmt"
	"github.com/netcloth/netcloth-chain/tests"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Upgrade(t *testing.T) {
	t.Parallel()
	_, servAddr, port, nchdHome, nchcliHome, p2pAddr := initFixtures(t)

	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("nchd start --home=%s --rpc.laddr=%v --p2p.laddr=%v", nchdHome, servAddr, p2pAddr))
	defer proc.Stop(false)

	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(1, port)

	cmdGetAccount := accountQueryCmd(t, "foo", nchcliHome, port)
	fooAccount := executeGetAccount(t, cmdGetAccount)
	//check foo account init coins
	require.Equal(t, fooAccount.Coins.AmountOf(sdk.NativeTokenName), sdk.NewInt(DefaultGenAccountAmount-1000000000000))
}
