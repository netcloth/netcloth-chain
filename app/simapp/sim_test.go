package simapp

import (
	v0 "github.com/netcloth/netcloth-chain/app/v0"
	"github.com/netcloth/netcloth-chain/app/v0/simulation"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	app "github.com/netcloth/netcloth-chain/app"
	baseapp "github.com/netcloth/netcloth-chain/baseapp"
)

func TestFullAppSimulation(t *testing.T) {
	config, db, dir, logger, skip, err := SetupSimulation("leveldb-app-sim", "Simulation")
	if skip {
		t.Skip("skipping application simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		db.Close()
		require.NoError(t, os.RemoveAll(dir))
	}()

	app := app.NewNCHApp(logger, db, nil, true, FlagPeriodValue, baseapp.FauxMerkleMode())
	require.Equal(t, "SimApp", app.Name())

	// run randomized simulation
	curProtocol := app.Engine.GetCurrentProtocol()
	cdc := curProtocol.GetCodec()
	sm := curProtocol.GetSimulationManager()
	_, _, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, AppStateFn(cdc, sm),
		SimulationOperations(app, cdc, config),
		v0.ModuleAccountAddrs(), config,
	)

	// export state and simParams before the simulation error is checked
	//err = CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		PrintStats(db)
	}
}
