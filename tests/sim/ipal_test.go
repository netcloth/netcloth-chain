package sim

import (
	"fmt"
	"testing"

	itypes "github.com/netcloth/netcloth-chain/app/v0/ipal/types"
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

func ipalClaimCmd(moniker, website, details, endpoints, amt, nchcliHome, port string) string {
	return fmt.Sprintf(`nchcli ipal claim -y --from foo --moniker=%s --website=%s --bond %s --details %s --home=%s  --endpoints=%s --node tcp://localhost:%s -o json`, moniker, website, amt, details, nchcliHome, endpoints, port)
}

func ipalQueryCmd(nchcliHome, port string) string {
	return fmt.Sprintf("nchcli q ipal list --home=%s --node tcp://localhost:%s -o json", nchcliHome, port)
}

func accountQueryCmd(t *testing.T, acc, nchcliHome, port string) string {
	fooAddr := executeGetAccAddress(t, fmt.Sprintf("nchcli keys show %s -a --home=%s", acc, nchcliHome))
	return fmt.Sprintf("nchcli query account %s --home=%s --node tcp://localhost:%s -o json", fooAddr, nchcliHome, port)
}

func executeGetIPALNodes(t *testing.T, cmdStr string) (nodes itypes.IPALNodes) {
	out, _ := tests.ExecuteT(t, cmdStr, "")

	cdc := MakeCodec()

	err := cdc.UnmarshalJSON([]byte(out), &nodes)
	require.NoError(t, err, "acc %v, err %v", string(out), err)

	return
}

func Test_ipal(t *testing.T) {
	t.Parallel()
	_, servAddr, port, nchdHome, nchcliHome, p2pAddr := initFixtures(t)

	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("nchd start --home=%s --rpc.laddr=%v --p2p.laddr=%v", nchdHome, servAddr, p2pAddr))
	defer proc.Stop(false)

	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(1, port)

	cmdGetAccount := accountQueryCmd(t, "foo", nchcliHome, port)
	initAccount := executeGetAccount(t, cmdGetAccount)
	//check foo account init coins
	require.Equal(t, initAccount.Coins.AmountOf(types.NativeTokenName), types.NewInt(DefaultGenAccountAmount-1000000000000))

	//failed for bond insufficient
	r := executeWrite(t, ipalClaimCmd("sky", "sky", "sky", "3|http://47.104.248.183", "100pnch", nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	nodes := executeGetIPALNodes(t, ipalQueryCmd(nchcliHome, port))
	require.Equal(t, len(nodes), 0)

	//succ to claim ipal
	bondAmount := types.NewCoin(types.NativeTokenName, types.NewInt(1000000000000000))
	r = executeWrite(t, ipalClaimCmd("sky", "sky.com", "skygood", "1|http://47.104.248.183", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	nodes = executeGetIPALNodes(t, ipalQueryCmd(nchcliHome, port))
	require.Equal(t, len(nodes), 1)
	require.Equal(t, nodes[0].Bond, bondAmount)
	require.Equal(t, nodes[0].Endpoints[0].Type, uint64(1))
	require.Equal(t, nodes[0].Moniker, "sky")
	require.Equal(t, nodes[0].Endpoints[0].Endpoint, "http://47.104.248.183")
	require.Equal(t, nodes[0].Details, "skygood")
	require.Equal(t, nodes[0].Website, "sky.com")

	//update bond amount
	bondAmount = types.NewCoin(types.NativeTokenName, types.NewInt(2000000000000000))
	r = executeWrite(t, ipalClaimCmd("sky", "sky", "sky", "3|http://47.104.248.183", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	nodes = executeGetIPALNodes(t, ipalQueryCmd(nchcliHome, port))
	require.Equal(t, len(nodes), 1)
	require.Equal(t, nodes[0].Bond, bondAmount)
	require.Equal(t, nodes[0].Endpoints[0].Type, uint64(3))
	require.Equal(t, nodes[0].Moniker, "sky")
	require.Equal(t, nodes[0].Endpoints[0].Endpoint, "http://47.104.248.183")
	require.Equal(t, nodes[0].Details, "sky")
	require.Equal(t, nodes[0].Website, "sky")

	//test for unbond
	bondAmount = types.NewCoin(types.NativeTokenName, types.NewInt(1000000000000000))
	r = executeWrite(t, ipalClaimCmd("sky", "sky", "sky", "1|http://47.104.248.183,3|http://47.104.248.183", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	nodes = executeGetIPALNodes(t, ipalQueryCmd(nchcliHome, port))
	require.Equal(t, len(nodes), 1)
	require.Equal(t, len(nodes[0].Endpoints), 2)
	require.Equal(t, nodes[0].Bond, bondAmount)

	tests.WaitForNextNBlocksTM(13, port)
	nodes = executeGetIPALNodes(t, ipalQueryCmd(nchcliHome, port))
	require.Equal(t, len(nodes), 1)
	require.Equal(t, nodes[0].Bond, bondAmount)

	//test for unbond
	bondAmount = types.NewCoin(types.NativeTokenName, types.NewInt(20000))
	r = executeWrite(t, ipalClaimCmd("sky", "sky", "sky", "1|http://47.104.248.183,3|http://47.104.248.183", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	nodes = executeGetIPALNodes(t, ipalQueryCmd(nchcliHome, port))
	require.Equal(t, len(nodes), 0)

	tests.WaitForNextNBlocksTM(13, port)
	nodes = executeGetIPALNodes(t, ipalQueryCmd(nchcliHome, port))
	require.Equal(t, len(nodes), 0)
}

func Test_unbond(t *testing.T) {
	//t.Parallel()

	_, servAddr, port, nchdHome, nchcliHome, p2pAddr := initFixtures(t)

	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("nchd start --home=%s --rpc.laddr=%v --p2p.laddr=%v", nchdHome, servAddr, p2pAddr))
	defer proc.Stop(false)

	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(1, port)

	bondAmount := types.NewCoin(types.NativeTokenName, types.NewInt(1000000000000))
	r := executeWrite(t, ipalClaimCmd("sky", "sky", "sky", "1|http://47.104.248.183", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	nodes := executeGetIPALNodes(t, ipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 1)

	bondAmount = types.NewCoin(types.NativeTokenName, types.NewInt(10000000000))
	r = executeWrite(t, ipalClaimCmd("sky", "sky", "sky", "1|http://47.104.248.183", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	nodes = executeGetIPALNodes(t, ipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 0)

	bondAmount = types.NewCoin(types.NativeTokenName, types.NewInt(2000000000000))
	r = executeWrite(t, ipalClaimCmd("sky", "sky", "sky", "1|http://47.104.248.183", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	nodes = executeGetIPALNodes(t, ipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 1)

	bondAmount = types.NewCoin(types.NativeTokenName, types.NewInt(10000000000))
	r = executeWrite(t, ipalClaimCmd("sky", "sky", "sky", "1|http://47.104.248.183", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	nodes = executeGetIPALNodes(t, ipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 0)

	bondAmount = types.NewCoin(types.NativeTokenName, types.NewInt(3000000000000))
	r = executeWrite(t, ipalClaimCmd("sky", "sky", "sky", "1|http://47.104.248.183", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	nodes = executeGetIPALNodes(t, ipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 1)

	bondAmount = types.NewCoin(types.NativeTokenName, types.NewInt(10000000000))
	r = executeWrite(t, ipalClaimCmd("sky", "sky", "sky", "1|http://47.104.248.183", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	nodes = executeGetIPALNodes(t, ipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 0)

	bondAmount = types.NewCoin(types.NativeTokenName, types.NewInt(3000000000000))
	r = executeWrite(t, ipalClaimCmd("sky", "sky", "sky", "1|http://47.104.248.183", bondAmount.String(), nchcliHome, port), DefaultKeyPass)
	require.True(t, r)
	tests.WaitForNextNBlocksTM(1, port)
	nodes = executeGetIPALNodes(t, ipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 1)

	tests.WaitForNextNBlocksTM(13, port)
	nodes = executeGetIPALNodes(t, ipalQueryCmd(nchcliHome, port))
	require.True(t, len(nodes) == 1)
}
