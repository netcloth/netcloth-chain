package sim

import (
    "fmt"
    "github.com/NetCloth/netcloth-chain/tests"
    "github.com/stretchr/testify/require"
    "testing"
)

func aipalClaimCmd(amt, nchcliHome, port string) string {
    return fmt.Sprintf("nchcli aipal claim --from=foo --moniker=foo --website=sky.com --server=sky.com --details=\" sky nb  \" --service_type=\"storage|chatting\" --bond=%s --home=%s --node localhost:%s", amt, nchcliHome, port)
}

func aipalParamsCmd(nchcliHome, port string) string {
    return fmt.Sprintf("nchcli q aipal params --home=%s --node localhost:%s -o json", nchcliHome, port)
}

func aipalQueryCmd(nchcliHome, port string) string {
    return fmt.Sprintf("nchcli q aipal servicenodes --home=%s --node localhost:%s -o json", nchcliHome, port)
}

func accountQueryCmd(t *testing.T, acc, nchcliHome, port string) string {
    fooAddr := executeGetAccAddress(t, fmt.Sprintf("nchcli keys show %s -a --home=%s", acc, nchcliHome))
    return fmt.Sprintf("nchcli query account %s --home=%s --node localhost:%s -o json", fooAddr, nchcliHome, port)
}

func Test_aipal(t *testing.T) {
    t.Parallel()

    _, servAddr, port, nchdHome, nchcliHome, p2pAddr := initFixtures(t)

    proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("nchd start --home=%s --rpc.laddr=%v --p2p.laddr=%v", nchdHome, servAddr, p2pAddr))
    defer proc.Stop(false)

    tests.WaitForTMStart(port)
    tests.WaitForNextNBlocksTM(1, port)


    cmdGetAccount := accountQueryCmd(t, "foo", nchcliHome, port)
    acc := executeGetAccount(t, cmdGetAccount)

    //failed for amount insufficient
    r := executeWrite(t, aipalClaimCmd("100000unch", nchcliHome, port), DefaultKeyPass)
    require.True(t, r)
    tests.WaitForNextNBlocksTM(1, port)
    newAcc := executeGetAccount(t, cmdGetAccount)
    require.True(t, acc.Coins.IsEqual(newAcc.Coins))
}
