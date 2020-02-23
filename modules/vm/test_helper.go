package vm

import (
	"os"
	"time"

	"github.com/netcloth/netcloth-chain/modules/auth/exported"
	authtype "github.com/netcloth/netcloth-chain/modules/auth/types"

	"github.com/tendermint/tendermint/crypto"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/netcloth/netcloth-chain/baseapp"
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/auth"
	"github.com/netcloth/netcloth-chain/modules/cipal"
	distr "github.com/netcloth/netcloth-chain/modules/distribution"
	"github.com/netcloth/netcloth-chain/modules/gov"
	"github.com/netcloth/netcloth-chain/modules/ipal"
	"github.com/netcloth/netcloth-chain/modules/mint"
	"github.com/netcloth/netcloth-chain/modules/params"
	"github.com/netcloth/netcloth-chain/modules/slashing"
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
		bam.MainStoreKey,
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
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	paramsKeeper := params.NewKeeper(types.ModuleCdc, keyParams, tkeyParams)
	accountKeeper := auth.NewAccountKeeper(types.ModuleCdc, keyAcc, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)

	keys := sdk.NewKVStoreKeys(auth.StoreKey, StoreKey, CodeKey)

	return NewEVM(Context{}, NewCommitStateDB(accountKeeper, keys[StoreKey], keys[CodeKey]), Config{})
}
