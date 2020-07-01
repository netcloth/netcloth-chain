package vm

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"github.com/netcloth/netcloth-chain/app/v0/params"
	"github.com/netcloth/netcloth-chain/app/v0/vm/types"
	"github.com/netcloth/netcloth-chain/store"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func newEVM() *EVM {
	authKey := sdk.NewKVStoreKey(auth.StoreKey)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	tParamsKey := sdk.NewTransientStoreKey(params.TStoreKey)

	paramsKeeper := params.NewKeeper(types.ModuleCdc, paramsKey, tParamsKey)
	accountKeeper := auth.NewAccountKeeper(types.ModuleCdc, authKey, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)

	logger := log.NewNopLogger()
	db := dbm.NewMemDB()

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeDB, db)
	ms.LoadLatestVersion()

	return NewEVM(Context{}, NewCommitStateDB(accountKeeper, authKey, authKey, authKey, sdk.NewKVStoreKey(StoreDebugKey)).WithContext(sdk.NewContext(ms, abci.Header{}, false, logger)), Config{})
}
