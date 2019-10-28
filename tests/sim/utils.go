package sim

import (
	"encoding/json"
	"fmt"
	"github.com/NetCloth/netcloth-chain/codec"
	"github.com/NetCloth/netcloth-chain/modules/aipal"
	"github.com/NetCloth/netcloth-chain/modules/auth"
	"github.com/NetCloth/netcloth-chain/modules/bank"
	"github.com/NetCloth/netcloth-chain/modules/crisis"
	"github.com/NetCloth/netcloth-chain/modules/distribution"
	"github.com/NetCloth/netcloth-chain/modules/gov"
	"github.com/NetCloth/netcloth-chain/modules/ipal"
	"github.com/NetCloth/netcloth-chain/modules/params"
	"github.com/NetCloth/netcloth-chain/modules/slashing"
	"github.com/NetCloth/netcloth-chain/modules/staking"
	"github.com/NetCloth/netcloth-chain/modules/supply"
	"github.com/NetCloth/netcloth-chain/server"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/types"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/NetCloth/netcloth-chain/client/keys"
	"github.com/NetCloth/netcloth-chain/tests"
	sdk "github.com/NetCloth/netcloth-chain/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

const (
	DefaultKeyPass = "12345678"
)

type KeyOutput struct {
	Name	string `json:"name"`
	Type	string `json:"type"`
	Address string `json:"address"`
	PubKey  string `json:"pubkey"`
	Seed    string `json:"seed,omitempty"`
}

type GenesisFileAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         []string       `json:"coins"`
	Sequence      uint64         `json:"sequence_number"`
	AccountNumber uint64         `json:"account_number"`
}

func getTestingHomeDirs(name string) (string, string) {
	tmpDir := os.TempDir()
	nchdHome := fmt.Sprintf("%s%s%s%s.test_nchd", tmpDir, string(os.PathSeparator), name, string(os.PathSeparator))
	nchcliHome := fmt.Sprintf("%s%s%s%s.test_nchcli", tmpDir, string(os.PathSeparator), name, string(os.PathSeparator))
	return nchdHome, nchcliHome
}

func initFixtures(t *testing.T) (chainID, servAddr, port, nchdHome, nchcliHome, p2p2Addr string) {
	nchdHome, nchcliHome = getTestingHomeDirs(t.Name())
	tests.ExecuteT(t, fmt.Sprintf("rm -rf %s ", nchdHome), "")

	executeWrite(t, fmt.Sprintf("nchcli keys delete --home=%s foo", nchcliHome), DefaultKeyPass)
	executeWrite(t, fmt.Sprintf("nchcli keys delete --home=%s bar", nchcliHome), DefaultKeyPass)
	executeWriteCheckErr(t, fmt.Sprintf("nchcli keys add --home=%s foo", nchcliHome), DefaultKeyPass)
	executeWriteCheckErr(t, fmt.Sprintf("nchcli keys add --home=%s bar", nchcliHome), DefaultKeyPass)

	chainID = executeInit(t, fmt.Sprintf("nchd init nch-foo -o --home=%s", nchdHome))
	tests.ExecuteT(t, fmt.Sprintf("nchcli config chain-id %s --home=%s", chainID, nchcliHome), "")
	tests.ExecuteT(t, fmt.Sprintf("nchcli config trust-node true --home=%s", nchcliHome), "")

	fooAccAddress := executeGetAccAddress(t, fmt.Sprintf("nchcli keys show foo -a --home=%s", nchcliHome))
	executeWrite(t, fmt.Sprintf("nchd add-genesis-account %s 11000000unch --home=%s", fooAccAddress, nchdHome), DefaultKeyPass)

	fooPubkey := executeGetAccAddress(t, fmt.Sprintf("nchd tendermint show-validator --home=%s", nchdHome)) //TODO refact executeGetAccAddress
	executeWrite(t, fmt.Sprintf("nchd gentx --amount 1000000unch --commission-rate 0.10 --commission-max-rate 0.20 --commission-max-change-rate 0.10 --pubkey %s --name foo --home=%s --home-client=%s", fooPubkey, nchdHome, nchcliHome), DefaultKeyPass)
	tests.ExecuteT(t, fmt.Sprintf("nchd collect-gentxs --home=%s", nchdHome), "")
	//genFile := filepath.Join(nchdHome, "config", "genesis.json")
	//genDoc := readGenesisFile(t, genFile)

	servAddr, port, err := server.FreeTCPAddr()
	require.NoError(t, err)

	p2p2Addr, _, err = server.FreeTCPAddr()
	require.NoError(t, err)

	return
}

func executeWrite(t *testing.T, cmdStr string, writes ...string) (exitSuccess bool) {
	if strings.Contains(cmdStr, "--from") && strings.Contains(cmdStr, "--fee") {
		cmdStr = cmdStr + " --commit"
	}

	exitSuccess, _, _ = executeWriteRetStreams(t, cmdStr, writes...)
	return
}

func executeWriteRetStreams(t *testing.T, cmdStr string, writes ...string) (bool, string, string) {
	proc := tests.GoExecuteT(t, cmdStr)

	for _, write := range writes {
		_, err := proc.StdinPipe.Write([]byte(write + "\n"))
		require.NoError(t, err)
	}

	stdout, stderr, err := proc.ReadAll()
	if err != nil {
		fmt.Println("Err on proc.ReadAll()", err, cmdStr)
	}

	if len(stdout) > 0 {
		t.Log("Stdout:", string(stdout))
		//t.Log("Stdout:", cmn.Green(string(stdout)))
	}

	if len(stderr) > 0 {
		t.Log("Stderr:", string(stderr))
	}

	proc.Wait()
	return proc.ExitState.Success(), string(stdout), string(stderr)
}

func executeWriteCheckErr(t *testing.T, cmdStr string, writes ...string) {
	require.True(t, executeWrite(t, cmdStr, writes...))
}

func executeInit(t *testing.T, cmdStr string) (chainID string) {
	_, stderr := tests.ExecuteT(t, cmdStr, DefaultKeyPass)

	var initRes map[string]json.RawMessage
	err := json.Unmarshal([]byte(stderr), &initRes)
	require.NoError(t, err)

	err = json.Unmarshal(initRes["chain_id"], &chainID)
	require.NoError(t, err)

	return
}

func executeGetAccAddress(t *testing.T, cmdStr string) (accAddress string) {
	stdout, _ := tests.ExecuteT(t, cmdStr, "")

	accAddress = string([]byte(stdout))
	return
}

func executeGetAddrPK(t *testing.T, cmdStr string) (sdk.AccAddress, crypto.PubKey) {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var ko KeyOutput
	keys.UnmarshalJSON([]byte(out), &ko)

	pk, err := sdk.GetAccPubKeyBech32(ko.PubKey)
	require.NoError(t, err)

	accAddr, err := sdk.AccAddressFromBech32(ko.Address)
	require.NoError(t, err)

	return accAddr, pk
}

func executeGetAccount(t *testing.T, cmdStr string) (acc auth.BaseAccount) {
	out, _ := tests.ExecuteT(t, cmdStr, "")

	var res map[string]json.RawMessage
	err := json.Unmarshal([]byte(out), &res)
	require.NoError(t, err, "out %v, err %v", out, err)

	cdc := MakeCodec()

	err = cdc.UnmarshalJSON([]byte(out), &acc)
	require.NoError(t, err, "acc %v, err %v", string(out), err)

	return
}

func readGenesisFile(t *testing.T, genFile string) types.GenesisDoc {
	var genDoc types.GenesisDoc
	fp, err := os.Open(genFile)
	require.NoError(t, err)
	fileContents, err := ioutil.ReadAll(fp)
	require.NoError(t, err)
	defer fp.Close()
	err = codec.Cdc.UnmarshalJSON(fileContents, &genDoc)
	require.NoError(t, err)
	return genDoc
}

func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	params.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	crisis.RegisterCodec(cdc)
	distribution.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	ipal.RegisterCodec(cdc)
	aipal.RegisterCodec(cdc)
	slashing.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{},
		"tendermint/PubKeySecp256k1", nil)
	return cdc
}
