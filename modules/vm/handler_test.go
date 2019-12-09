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
	amount := sdk.NewInt64Coin(sdk.NativeTokenName, 150)
	code := sdk.FromHex("608060405260008055600060015534801561001957600080fd5b50610174806100296000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80633033413b146100515780635d33a27f1461006f578063a17a9e661461008d578063dac0eb07146100cf575b600080fd5b6100596100fd565b6040518082815260200191505060405180910390f35b610077610103565b6040518082815260200191505060405180910390f35b6100b9600480360360208110156100a357600080fd5b8101908080359060200190929190505050610109565b6040518082815260200191505060405180910390f35b6100fb600480360360208110156100e557600080fd5b8101908080359060200190929190505050610124565b005b60005481565b60015481565b60008160008082825401925050819055506000549050919050565b61012d81610109565b6001600082825401925050819055505056fea265627a7a72315820a86d07b87e00fb978bd8446525d467cb857ac5f65dcef76cd39123c8d11a717864736f6c634300050d0032")
	fmt.Println(code)

	msg := types.NewMsgContractCreate(fromAddr, amount, code)
	require.NotNil(t, msg)
	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), types.TypeMsgContractCreate)

	k, ctx := setupTest()
	h := NewHandler(k)

	res := h(ctx, msg)

	require.False(t, res.IsOK())
	fmt.Println("logs: ", res.Log)
}

func TestMsgContractCall(t *testing.T) {
}
