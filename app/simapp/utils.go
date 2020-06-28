package simapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"math/rand"

	"github.com/tendermint/tendermint/libs/log"
	//"github.com/tendermint/tendermint/types/time"
	dbm "github.com/tendermint/tm-db"

	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/simapp/helpers"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/types/module"
	simtypes "github.com/netcloth/netcloth-chain/types/simulation"
)

// SetupSimulation creates the config, db (levelDB), temporary directory and logger for
// the simulation tests. If `FlagEnabledValue` is false it skips the current test.
// Returns error on an invalid db intantiation or temp dir creation.
func SetupSimulation(dirPrefix, dbName string) (simtypes.Config, dbm.DB, string, log.Logger, bool, error) {
	if !FlagEnabledValue {
		return simtypes.Config{}, nil, "", nil, true, nil
	}

	config := NewConfigFromFlags()
	//r := rand.New(rand.NewSource(time.Now().Unix()))
	//config.Seed = r.Int63()
	config.ChainID = helpers.SimAppChainID

	var logger log.Logger
	if FlagVerboseValue {
		logger = log.TestingLogger()
	} else {
		logger = log.NewNopLogger()
	}

	dir, err := ioutil.TempDir("", dirPrefix)
	if err != nil {
		return simtypes.Config{}, nil, "", nil, false, err
	}

	db, err := sdk.NewLevelDB(dbName, dir)
	if err != nil {
		return simtypes.Config{}, nil, "", nil, false, err
	}

	return config, db, dir, logger, false, nil
}

// SimulationOperations retrieves the simulation params from the provided file path
// and returns all the modules weighted operations
func SimulationOperations(app App, cdc *codec.Codec, config simtypes.Config) []simtypes.WeightedOperation {
	simState := module.SimulationState{
		AppParams: make(simtypes.AppParams),
		Cdc:       cdc,
	}

	if config.ParamsFile != "" {
		bz, err := ioutil.ReadFile(config.ParamsFile)
		if err != nil {
			panic(err)
		}

		app.Codec().MustUnmarshalJSON(bz, &simState.AppParams)
	}

	simState.ParamChanges = app.SimulationManager().GenerateParamChanges(config.Seed)
	simState.Contents = app.SimulationManager().GetProposalContents(simState)
	return app.SimulationManager().WeightedOperations(simState)
}

// CheckExportSimulation exports the app state and simulation parameters to JSON
// if the export paths are defined.
func CheckExportSimulation(
	app App, config simtypes.Config, params simtypes.Params,
) error {
	//if config.ExportStatePath != "" {
	//	fmt.Println("exporting app state...")
	//	appState, _, _, err := app.ExportAppStateAndValidators(false, nil)
	//	if err != nil {
	//		return err
	//	}
	//
	//	if err := ioutil.WriteFile(config.ExportStatePath, []byte(appState), 0600); err != nil {
	//		return err
	//	}
	//}

	if config.ExportParamsPath != "" {
		fmt.Println("exporting simulation params...")
		paramsBz, err := json.MarshalIndent(params, "", " ")
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(config.ExportParamsPath, paramsBz, 0600); err != nil {
			return err
		}
	}
	return nil
}

// PrintStats prints the corresponding statistics from the app DB.
func PrintStats(db dbm.DB) {
	fmt.Println("\nLevelDB Stats")
	fmt.Println(db.Stats()["leveldb.stats"])
	fmt.Println("LevelDB cached block size", db.Stats()["leveldb.cachedblock"])
}
