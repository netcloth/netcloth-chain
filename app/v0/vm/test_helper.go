package vm

import (
	"os"
	"time"

	"github.com/netcloth/netcloth-chain/app/protocol"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"github.com/netcloth/netcloth-chain/app/v0/auth/exported"
	authtype "github.com/netcloth/netcloth-chain/app/v0/auth/types"
	"github.com/netcloth/netcloth-chain/app/v0/cipal"
	distr "github.com/netcloth/netcloth-chain/app/v0/distribution"
	"github.com/netcloth/netcloth-chain/app/v0/gov"
	"github.com/netcloth/netcloth-chain/app/v0/ipal"
	"github.com/netcloth/netcloth-chain/app/v0/mint"
	"github.com/netcloth/netcloth-chain/app/v0/params"
	"github.com/netcloth/netcloth-chain/app/v0/slashing"
	"github.com/netcloth/netcloth-chain/app/v0/staking"
	"github.com/netcloth/netcloth-chain/app/v0/supply"
	"github.com/netcloth/netcloth-chain/app/v0/vm/types"
	"github.com/netcloth/netcloth-chain/codec"
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

func ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}
	return modAccAddrs
}

func KeyTestPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func setupTest() (vmKeeper Keeper, ctx sdk.Context) {
	cdc := codec.New()
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	cdc.RegisterConcrete(&authtype.BaseAccount{}, "nch/Account", nil)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	keys := sdk.NewKVStoreKeys(
		protocol.MainStoreKey,
		auth.StoreKey,
		auth.RefundKey,
		staking.StoreKey,
		supply.StoreKey,
		mint.StoreKey,
		distr.StoreKey,
		slashing.StoreKey,
		gov.StoreKey,
		params.StoreKey,
		cipal.StoreKey,
		ipal.StoreKey,
		StoreKey,
		CodeKey,
		StoreDebugKey,
	)
	tkeys := sdk.NewTransientStoreKeys(staking.TStoreKey, staking.TStoreKey, params.TStoreKey)

	paramsKeeper := params.NewKeeper(cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
	authSubspace := paramsKeeper.Subspace(auth.DefaultParamspace)

	vmSubspace := paramsKeeper.Subspace(DefaultParamspace)

	// add keepers
	accountKeeper := auth.NewAccountKeeper(cdc, keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)

	vmKeeper = NewKeeper(
		cdc,
		keys[StoreKey],
		keys[CodeKey],
		keys[StoreDebugKey],
		vmSubspace,
		accountKeeper)

	for _, key := range keys {
		ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, nil) // db nil
	}
	for _, key := range tkeys {
		ms.MountStoreWithDB(key, sdk.StoreTypeTransient, nil) // db nil
	}
	ms.LoadLatestVersion()

	ctx = sdk.NewContext(ms, abci.Header{Time: time.Unix(0, 0)}, false, log.NewTMLogger(os.Stdout))
	vmKeeper.SetParams(ctx, types.DefaultParams())

	return
}

func GetTestAccount() auth.BaseAccount {
	_, pubKey, addr := KeyTestPubAddr()
	acc := auth.NewBaseAccountWithAddress(addr)
	acc.SetPubKey(pubKey)
	acc.SetSequence(0)
	acc.SetCoins(sdk.NewCoins(sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(10000000000))))

	return acc
}

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

	return NewEVM(Context{}, NewCommitStateDB(accountKeeper, authKey, authKey, sdk.NewKVStoreKey(StoreDebugKey)).WithContext(sdk.NewContext(ms, abci.Header{}, false, logger)), Config{})
}
