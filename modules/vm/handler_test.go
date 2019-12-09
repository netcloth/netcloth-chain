package vm

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/auth"
	"github.com/netcloth/netcloth-chain/modules/bank"
	distr "github.com/netcloth/netcloth-chain/modules/distribution"
	"github.com/netcloth/netcloth-chain/modules/gov"
	"github.com/netcloth/netcloth-chain/modules/mint"
	"github.com/netcloth/netcloth-chain/modules/params"
	"github.com/netcloth/netcloth-chain/modules/staking"
	"github.com/netcloth/netcloth-chain/modules/supply"
	"github.com/netcloth/netcloth-chain/modules/vm/types"
	"github.com/netcloth/netcloth-chain/store"
	sdk "github.com/netcloth/netcloth-chain/types"
)

var (
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
	}
)

func moduleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}
	return modAccAddrs
}

func setupTest() (keeper Keeper, ctx sdk.Context) {
	cdc := codec.New()

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	keys := sdk.NewKVStoreKeys(params.StoreKey, auth.StoreKey, supply.StoreKey, staking.StoreKey)
	tkeys := sdk.NewTransientStoreKeys(params.TStoreKey, staking.TStoreKey)

	storageKey = sdk.NewKVStoreKey("store")
	codeKey = sdk.NewKVStoreKey("code")
	tStoreKey := sdk.NewTransientStoreKey("transient_store")

	paramsKeeper := params.NewKeeper(cdc, keys[params.StoreKey], tkeys[params.TStoreKey], params.DefaultCodespace)

	authSubspace := paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := paramsKeeper.Subspace(bank.DefaultParamspace)

	ms.MountStoreWithDB(keys[auth.StoreKey], sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeys[staking.TStoreKey], sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keys[staking.StoreKey], sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys[supply.StoreKey], sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys[params.StoreKey], sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeys[params.TStoreKey], sdk.StoreTypeTransient, db)

	ms.LoadLatestVersion()

	accountKeeper := auth.NewAccountKeeper(cdc, keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, bankSubspace, bank.DefaultCodespace, moduleAccountAddrs())

	keeper = NewKeeper(cdc, storageKey, tStoreKey, types.DefaultCodespace, params.NewSubspace(cdc, keyParams, tkeyParams, "param_subspace"), accountKeeper, bankKeeper, NewCommitStateDB(accountKeeper, bankKeeper, storageKey, codeKey))
	ctx = sdk.NewContext(ms, abci.Header{Time: time.Unix(0, 0)}, false, log.NewTMLogger(os.Stdout))

	return
}

func newSdkAddress() sdk.AccAddress {
	tmpKey := secp256k1.GenPrivKey().PubKey()
	return sdk.AccAddress(tmpKey.Address().Bytes())
}

func TestInvalidMsg(t *testing.T) {
	k := Keeper{}
	h := NewHandler(k)

	res := h(sdk.NewContext(nil, abci.Header{}, false, nil), sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized Msg type"))
}

func TestMsgContractCreate(t *testing.T) {
	fromAddr := newSdkAddress()
	amount := sdk.NewInt64Coin(sdk.NativeTokenName, 0)
	code := sdk.FromHex("608060405234801561001057600080fd5b506102d3806100206000396000f30060806040526004361061004c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630dbe671f14610051578063abd2178d146100e1575b600080fd5b34801561005d57600080fd5b5061006661014a565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100a657808201518184015260208101905061008b565b50505050905090810190601f1680156100d35780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156100ed57600080fd5b50610148600480360381019080803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506101e8565b005b60008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156101e05780601f106101b5576101008083540402835291602001916101e0565b820191906000526020600020905b8154815290600101906020018083116101c357829003601f168201915b505050505081565b80600090805190602001906101fe929190610202565b5050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061024357805160ff1916838001178555610271565b82800160010185558215610271579182015b82811115610270578251825591602001919060010190610255565b5b50905061027e9190610282565b5090565b6102a491905b808211156102a0576000816000905550600101610288565b5090565b905600a165627a7a723058203fbad36965a5b8ae9c581ead091230bc57c9232bbede655626ff299b08ff97d50029")
	fmt.Println(code)

	msg := types.NewMsgContractCreate(fromAddr, amount, code)
	require.NotNil(t, msg)
	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), types.TypeMsgContractCreate)

	k, ctx := setupTest()
	h := NewHandler(k)

	res := h(ctx, msg)

	require.True(t, res.IsOK())
	fmt.Println("logs: ", res.Log)
}

func TestMsgContractCall(t *testing.T) {
}
