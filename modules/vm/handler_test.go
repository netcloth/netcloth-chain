package vm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/auth"
	"github.com/netcloth/netcloth-chain/modules/bank"
	"github.com/netcloth/netcloth-chain/modules/params"
	"github.com/netcloth/netcloth-chain/modules/vm/types"
	"github.com/netcloth/netcloth-chain/store"
	sdk "github.com/netcloth/netcloth-chain/types"
)

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
	var (
		storeKey      = sdk.NewKVStoreKey("store")
		tStoreKey     = sdk.NewTransientStoreKey("transient_store")
		keyAcc        = sdk.NewKVStoreKey(auth.StoreKey)
		keyParams     = sdk.NewKVStoreKey(params.StoreKey)
		tkeyParams    = sdk.NewTransientStoreKey(params.TStoreKey)
		paramsKeeper  = params.NewKeeper(types.ModuleCdc, keyParams, tkeyParams, params.DefaultCodespace)
		accountKeeper = auth.NewAccountKeeper(types.ModuleCdc, keyAcc, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
		bankKeeper    = bank.NewBaseKeeper(accountKeeper, paramsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, nil)

		db = dbm.NewMemDB()
		ms = store.NewCommitMultiStore(db)
	)

	fromAddr := newSdkAddress()
	amount := sdk.NewInt64Coin(sdk.NativeTokenName, 150)
	code := sdk.FromHex("608060405260008055600060015534801561001957600080fd5b50610174806100296000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80633033413b146100515780635d33a27f1461006f578063a17a9e661461008d578063dac0eb07146100cf575b600080fd5b6100596100fd565b6040518082815260200191505060405180910390f35b610077610103565b6040518082815260200191505060405180910390f35b6100b9600480360360208110156100a357600080fd5b8101908080359060200190929190505050610109565b6040518082815260200191505060405180910390f35b6100fb600480360360208110156100e557600080fd5b8101908080359060200190929190505050610124565b005b60005481565b60015481565b60008160008082825401925050819055506000549050919050565b61012d81610109565b6001600082825401925050819055505056fea265627a7a72315820a86d07b87e00fb978bd8446525d467cb857ac5f65dcef76cd39123c8d11a717864736f6c634300050d0032")
	fmt.Println(code)

	msg := types.NewMsgContractCreate(fromAddr, amount, code)
	require.NotNil(t, msg)
	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), types.TypeMsgContractCreate)

	codec := codec.New()
	k := NewKeeper(codec, storeKey, tStoreKey, types.DefaultCodespace, params.NewSubspace(codec, keyParams, tkeyParams, "param_subspace"), accountKeeper, bankKeeper, NewCommitStateDB(accountKeeper, bankKeeper, storageKey, codeKey))
	h := NewHandler(k)

	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tStoreKey, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	res := h(sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger()), msg)

	require.False(t, res.IsOK())
	fmt.Println("logs: ", res.Log)
}

func TestMsgContractCall(t *testing.T) {
}
