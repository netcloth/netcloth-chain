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
	code := []byte("0x60806040526000805534801561001457600080fd5b5060ce806100236000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80631003e2d214602d575b600080fd5b605660048036036020811015604157600080fd5b81019080803590602001909291905050506077565b60405180831515151581526020018281526020019250505060405180910390f35b600080826000808282540192505081905550826000541015609457fe5b91509156fea265627a7a7231582028f571a0c96eb3df56b211520f012b6a45280a1b0ce349f80991da8d1a443dd364736f6c634300050d0032")

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
