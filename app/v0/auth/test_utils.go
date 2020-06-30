package auth

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	authtypes "github.com/netcloth/netcloth-chain/app/v0/auth/types"
	"github.com/netcloth/netcloth-chain/app/v0/params"
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/store"
	sdk "github.com/netcloth/netcloth-chain/types"
)

type testInput struct {
	cdc *codec.Codec
	ctx sdk.Context
	ak  AccountKeeper
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()
	cdc := authtypes.ModuleCdc

	authCapKey := sdk.NewKVStoreKey("auth")
	keyParams := sdk.NewKVStoreKey("params")
	tKeyParams := sdk.NewTransientStoreKey("transient_params")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	pk := params.NewKeeper(cdc, keyParams, tKeyParams)
	ak := NewAccountKeeper(cdc, authCapKey, pk.Subspace(DefaultParamspace), ProtoBaseAccount)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	return testInput{cdc: cdc, ctx: ctx, ak: ak}
}
