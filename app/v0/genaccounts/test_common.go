package genaccounts

import (
	"github.com/netcloth/netcloth-chain/app/v0/auth"
	authtypes "github.com/netcloth/netcloth-chain/app/v0/auth/types"
	"github.com/netcloth/netcloth-chain/app/v0/genaccounts/internal/types"
	"github.com/netcloth/netcloth-chain/app/v0/params/subspace"
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/store"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/tendermint/tendermint/libs/log"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	dbm "github.com/tendermint/tm-db"
)

type testInput struct {
	cdc *codec.Codec
	ctx sdk.Context
	ak  types.AccountKeeper
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	auth.RegisterCodec(cdc)
	authCapKey := sdk.NewKVStoreKey("authCapKey")
	keyParams := sdk.NewKVStoreKey("subspace")
	tkeyParams := sdk.NewTransientStoreKey("transient_subspace")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	if err := ms.LoadLatestVersion(); err != nil {
		panic(err)
	}

	ps := subspace.NewSubspace(cdc, keyParams, tkeyParams, "genaccounts")
	ak := auth.NewAccountKeeper(cdc, authCapKey, ps, authtypes.ProtoBaseAccount)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	return testInput{cdc: cdc, ctx: ctx, ak: ak}
}
func KeyTestPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}
