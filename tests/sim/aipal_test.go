package sim

import (
	"fmt"
	"testing"

	itypes "github.com/netcloth/netcloth-chain/modules/ipal/types"
	"github.com/netcloth/netcloth-chain/tests"
	"github.com/netcloth/netcloth-chain/types"
	"github.com/stretchr/testify/require"
)

/*
PS important
test prepare:
1. must set modules/ipal/types/params.go: DefaultUnbondingTime = time.Minute * 1
2. after step1, must recompile and install nchd and nchcli on testting host

otherwise test cases may fail
*/

func aipalClaimCmd(server, moniker, website, details, serviceType, amt, nchcliHome, port string) string {
	return fmt.Sprintf(`nchcli ipal claim -y --from foo --server %s --moniker=%s --website=%s --bond %s --details %s --home=%s --service_type=%s --node localhost:%s -o json`, server, moniker, website, amt, details, nchcliHome, serviceType, port)
}

func aipalParamsCmd(nchcliHome, port string) string {
	return fmt.Sprintf("nchcli q ipal params --home=%s --node localhost:%s -o json", nchcliHome, port)
}

func aipalQueryCmd(nchcliHome, port string) string {
	return fmt.Sprintf("nchcli q ipal servicenodes --home=%s --node localhost:%s -o json", nchcliHome, port)
}

func accountQueryCmd(t *testing.T, acc, nchcliHome, port string) string {
	fooAddr := executeGetAccAddress(t, fmt.Sprintf("nchcli keys show %s -a --home=%s", acc, nchcliHome))
	return fmt.Sprintf("nchcli query account %s --home=%s --node localhost:%s -o json", fooAddr, nchcliHome, port)
}

func executeGetServiceNodes(t *testing.T, cmdStr string) (nodes itypes.ServiceNodes) {
	out, _ := tests.ExecuteT(t, cmdStr, "")

	cdc := MakeCodec()

	err := cdc.UnmarshalJSON([]byte(out), &nodes)
	require.NoError(t, err, "acc %v, err %v", string(out), err)

	return
}

func Test_aipal(t *testing.T) {
	//t.Parallel()

	_, servAddr, port, nchdHome, nchcliHome, p2pAddr := initFixtures(t)

	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("nchd start --home=%s --rpc.laddr=%v --p2p.laddr=%v", nchdHome, servAddr, p2pAddr))
	defer proc.Stop(false)

	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(1, port)

	cmdGetAccount := accountQueryCmd(t, "foo", nchcliHome, port)
	initAccount := executeGetAccount(t, cmdGetAccount)
	//foo Account's init pnch is 10000000pnch
	require.True(t, initAccount.Coins.IsEqual(types.NewCoins(types.NewCoin(types.NativeTokenName, types.NewInt(10000000)))))

	//failed for bond insufficient
	r := executeWrite(t, aipalClaimCmd("sky", "sky", "sky", "sky", "1", "100000pnch", nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	tmpAccount := executeGetAccount(t, cmdGetAccount)
	require.True(t, tmpAccount.Coins.IsEqual(initAccount.Coins))

	//succ to claim ipal
	bondAmount := types.NewCoin(types.NativeTokenName, types.NewInt(2000000))
	r = executeWrite(t, aipalClaimCmd("sky", "sky", "sky.com", "skygood", "1", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	tmpAccount = executeGetAccount(t, cmdGetAccount)
	require.True(t, tmpAccount.Coins.Add(types.NewCoins(bondAmount)).IsEqual(initAccount.Coins))
	nodes := executeGetServiceNodes(t, aipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 1)
	require.True(t, nodes[0].Bond.IsEqual(bondAmount))
	require.True(t, nodes[0].ServiceType == 1)
	require.True(t, nodes[0].Moniker == "sky")
	require.True(t, nodes[0].ServerEndPoint == "sky")
	require.True(t, nodes[0].Details == "skygood")
	require.True(t, nodes[0].Website == "sky.com")

	//update bond amount
	bondAmount = types.NewCoin(types.NativeTokenName, types.NewInt(3000000))
	r = executeWrite(t, aipalClaimCmd("sky", "sky", "sky", "sky", "3", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	tmpAccount = executeGetAccount(t, cmdGetAccount)
	require.True(t, tmpAccount.Coins.Add(types.NewCoins(bondAmount)).IsEqual(initAccount.Coins))
	nodes = executeGetServiceNodes(t, aipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 1)
	require.True(t, nodes[0].Bond.IsEqual(bondAmount))
	require.True(t, nodes[0].ServiceType == 3)
	require.True(t, nodes[0].Moniker == "sky")
	require.True(t, nodes[0].ServerEndPoint == "sky")
	require.True(t, nodes[0].Details == "sky")
	require.True(t, nodes[0].Website == "sky")

	//test for unbond
	bondAmount = types.NewCoin(types.NativeTokenName, types.NewInt(2000000))
	r = executeWrite(t, aipalClaimCmd("sky", "sky", "sky", "sky", "2", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	tmpAccount = executeGetAccount(t, cmdGetAccount)
	require.True(t, tmpAccount.Coins.Add(types.NewCoins(types.NewCoin(types.NativeTokenName, types.NewInt(3000000)))).IsEqual(initAccount.Coins))
	nodes = executeGetServiceNodes(t, aipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 1)
	require.True(t, nodes[0].Bond.IsEqual(bondAmount))

	tests.WaitForNextNBlocksTM(13, port)
	tmpAccount = executeGetAccount(t, cmdGetAccount)
	require.True(t, tmpAccount.Coins.Add(types.NewCoins(bondAmount)).IsEqual(initAccount.Coins), tmpAccount.String())
	nodes = executeGetServiceNodes(t, aipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 1)
	require.True(t, nodes[0].Bond.IsEqual(bondAmount))

	//test for unbond
	bondAmount = types.NewCoin(types.NativeTokenName, types.NewInt(20000))
	r = executeWrite(t, aipalClaimCmd("sky", "sky", "sky", "sky", "2", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	tmpAccount = executeGetAccount(t, cmdGetAccount)
	require.True(t, tmpAccount.Coins.Add(types.NewCoins(types.NewCoin(types.NativeTokenName, types.NewInt(2000000)))).IsEqual(initAccount.Coins))
	nodes = executeGetServiceNodes(t, aipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 0)

	tests.WaitForNextNBlocksTM(13, port)
	tmpAccount = executeGetAccount(t, cmdGetAccount)
	require.True(t, tmpAccount.Coins.IsEqual(initAccount.Coins), tmpAccount.String())
	nodes = executeGetServiceNodes(t, aipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 0)
}

func Test_unbond(t *testing.T) {
	//t.Parallel()

	_, servAddr, port, nchdHome, nchcliHome, p2pAddr := initFixtures(t)

	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("nchd start --home=%s --rpc.laddr=%v --p2p.laddr=%v", nchdHome, servAddr, p2pAddr))
	defer proc.Stop(false)

	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(1, port)

	cmdGetAccount := accountQueryCmd(t, "foo", nchcliHome, port)
	initAccount := executeGetAccount(t, cmdGetAccount)
	//foo Account's init pnch is 10000000pnch
	require.True(t, initAccount.Coins.IsEqual(types.NewCoins(types.NewCoin(types.NativeTokenName, types.NewInt(10000000)))))

	bondAmount := types.NewCoin(types.NativeTokenName, types.NewInt(1000000))
	r := executeWrite(t, aipalClaimCmd("sky", "sky", "sky", "sky", "2", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	tmpAccount := executeGetAccount(t, cmdGetAccount)
	require.True(t, tmpAccount.Coins.Add(types.NewCoins(types.NewCoin(types.NativeTokenName, types.NewInt(1000000)))).IsEqual(initAccount.Coins), tmpAccount.Coins.String())
	nodes := executeGetServiceNodes(t, aipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 1)

	bondAmount = types.NewCoin(types.NativeTokenName, types.NewInt(10000))
	r = executeWrite(t, aipalClaimCmd("sky", "sky", "sky", "sky", "2", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	tmpAccount = executeGetAccount(t, cmdGetAccount)
	require.True(t, tmpAccount.Coins.Add(types.NewCoins(types.NewCoin(types.NativeTokenName, types.NewInt(1000000)))).IsEqual(initAccount.Coins), tmpAccount.Coins.String())
	nodes = executeGetServiceNodes(t, aipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 0)

	bondAmount = types.NewCoin(types.NativeTokenName, types.NewInt(2000000))
	r = executeWrite(t, aipalClaimCmd("sky", "sky", "sky", "sky", "2", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	tmpAccount = executeGetAccount(t, cmdGetAccount)
	require.True(t, tmpAccount.Coins.Add(types.NewCoins(types.NewCoin(types.NativeTokenName, types.NewInt(3000000)))).IsEqual(initAccount.Coins), tmpAccount.Coins.String())
	nodes = executeGetServiceNodes(t, aipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 1)

	bondAmount = types.NewCoin(types.NativeTokenName, types.NewInt(10000))
	r = executeWrite(t, aipalClaimCmd("sky", "sky", "sky", "sky", "2", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	tmpAccount = executeGetAccount(t, cmdGetAccount)
	require.True(t, tmpAccount.Coins.Add(types.NewCoins(types.NewCoin(types.NativeTokenName, types.NewInt(3000000)))).IsEqual(initAccount.Coins), tmpAccount.Coins.String())
	nodes = executeGetServiceNodes(t, aipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 0)

	bondAmount = types.NewCoin(types.NativeTokenName, types.NewInt(3000000))
	r = executeWrite(t, aipalClaimCmd("sky", "sky", "sky", "sky", "2", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	tmpAccount = executeGetAccount(t, cmdGetAccount)
	require.True(t, tmpAccount.Coins.Add(types.NewCoins(types.NewCoin(types.NativeTokenName, types.NewInt(6000000)))).IsEqual(initAccount.Coins), tmpAccount.Coins.String())
	nodes = executeGetServiceNodes(t, aipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 1)

	bondAmount = types.NewCoin(types.NativeTokenName, types.NewInt(10000))
	r = executeWrite(t, aipalClaimCmd("sky", "sky", "sky", "sky", "2", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	tmpAccount = executeGetAccount(t, cmdGetAccount)
	require.True(t, tmpAccount.Coins.Add(types.NewCoins(types.NewCoin(types.NativeTokenName, types.NewInt(6000000)))).IsEqual(initAccount.Coins), tmpAccount.Coins.String())
	nodes = executeGetServiceNodes(t, aipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 0)

	bondAmount = types.NewCoin(types.NativeTokenName, types.NewInt(3000000))
	r = executeWrite(t, aipalClaimCmd("sky", "sky", "sky", "sky", "2", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	tmpAccount = executeGetAccount(t, cmdGetAccount)
	require.True(t, tmpAccount.Coins.Add(types.NewCoins(types.NewCoin(types.NativeTokenName, types.NewInt(9000000)))).IsEqual(initAccount.Coins), tmpAccount.Coins.String())
	nodes = executeGetServiceNodes(t, aipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 1)

	tests.WaitForNextNBlocksTM(13, port)
	tmpAccount = executeGetAccount(t, cmdGetAccount)
	require.True(t, tmpAccount.Coins.Add(types.NewCoins(types.NewCoin(types.NativeTokenName, types.NewInt(3000000)))).IsEqual(initAccount.Coins), tmpAccount.Coins.String())
	nodes = executeGetServiceNodes(t, aipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 1)
}
