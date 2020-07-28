package subspace

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/store"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	sdk "github.com/netcloth/netcloth-chain/types"
)

var (
	keyUnbondingTime = []byte("UnbondingTime")
	keyMaxValidators = []byte("MaxValidators")
	keyBondDenom     = []byte("BondDenom")

	key  = sdk.NewKVStoreKey("storekey")
	tkey = sdk.NewTransientStoreKey("transientstorekey")
)

// Keys for parameter access
const (
	TestParamStore = "ParamsTest"
)

type params struct {
	UnbondingTime time.Duration `json:"unbonding_time" yaml:"unbonding_time"`
	MaxValidators uint16        `json:"max_validators" yaml:"max_validators"`
	BondDenom     string        `json:"bond_denom" yaml:"bond_denom"`
}

func validateUnbondingTime(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < (24 * time.Hour) {
		return fmt.Errorf("unbonding time must be at least one day")
	}

	return nil
}

func validateMaxValidators(i interface{}) error {
	_, ok := i.(uint16)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateBondDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(v) == 0 {
		return errors.New("denom cannot be empty")
	}

	return nil
}

func (p *params) ParamSetPairs() ParamSetPairs {
	return ParamSetPairs{
		{keyUnbondingTime, &p.UnbondingTime, validateUnbondingTime},
		{keyMaxValidators, &p.MaxValidators, validateMaxValidators},
		{keyBondDenom, &p.BondDenom, validateBondDenom},
	}
}

func paramKeyTable() KeyTable {
	return NewKeyTable().RegisterParamSet(&params{})
}

// Returns components for testing
func DefaultTestComponents(t *testing.T) (sdk.Context, Subspace, func() sdk.CommitID) {
	cdc := codec.New()
	key := sdk.NewKVStoreKey(StoreKey)
	tkey := sdk.NewTransientStoreKey(TStoreKey)
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.SetTracer(os.Stdout)
	ms.SetTracingContext(sdk.TraceContext{})
	ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewTMLogger(os.Stdout))
	subspace := NewSubspace(cdc, key, tkey, TestParamStore)

	return ctx, subspace, ms.Commit
}
