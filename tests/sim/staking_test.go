package sim

import (
	"fmt"
	"github.com/NetCloth/netcloth-chain/tests"
	"testing"
)

func TestMock(t *testing.T) {
	t.Parallel()

	_, servAddr, port, nchdHome, nchcliHome, p2pAddr := initFixtures(t)

	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("nchd start --home=%s --rpc.laddr=%v --p2p.laddr=%v", nchdHome, servAddr, p2pAddr))
	defer proc.Stop(false)

	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(1, port)

	fooAddr := executeGetAccAddress(t, fmt.Sprintf("nchcli keys show foo -a --home=%s", nchcliHome))

	fooAcc := executeGetAccount(t, fmt.Sprintf("nchcli query account %s --home=%s --node localhost:%s -o json", fooAddr, nchcliHome, port))
	fmt.Println(fooAcc.Coins.String())
}