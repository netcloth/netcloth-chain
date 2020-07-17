package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	mempl "github.com/tendermint/tendermint/mempool"
	"github.com/tendermint/tendermint/proxy"
	tmsm "github.com/tendermint/tendermint/state"
	tmstore "github.com/tendermint/tendermint/store"
	tm "github.com/tendermint/tendermint/types"

	"github.com/netcloth/netcloth-chain/app"
	"github.com/netcloth/netcloth-chain/server"
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	DefaultReplayFromHeight = 1
	DefaultReplayToHeight   = -1

	flagReplayFromHeight = "from"
	flagReplayToHeight   = "to"
)

func replayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replay <root-dir> --from [from_height] --to [to_height]",
		Short: "Replay nchd transactions",

		RunE: func(_ *cobra.Command, args []string) error {
			from := viper.GetInt64(flagReplayFromHeight)
			to := viper.GetInt64(flagReplayToHeight)
			if from <= 0 {
				return fmt.Errorf("from block height must >= 1")
			}

			return replayTxs(args[0], from, to)
		},
		Args: cobra.ExactArgs(1),
	}

	cmd.Flags().Int64(flagReplayFromHeight, DefaultReplayFromHeight, "from block height")
	cmd.Flags().Int64(flagReplayToHeight, DefaultReplayToHeight, "to block height")

	return cmd
}

func replayTxs(rootDir string, from, to int64) error {
	configDir := filepath.Join(rootDir, "config")
	dataDir := filepath.Join(rootDir, "data")
	ctx := server.NewDefaultContext()

	if DefaultReplayFromHeight == from {
		statedbDir := filepath.Join(dataDir, "state.db")
		appdbDir := filepath.Join(dataDir, "application.db")

		fmt.Printf("state database: %s\n", statedbDir)
		fmt.Printf("app database: %s\n", appdbDir)
		err := os.RemoveAll(statedbDir)
		if err != nil {
			return err
		}

		err = os.RemoveAll(appdbDir)
		if err != nil {
			return err
		}
	}

	// App DB
	// appDB := dbm.NewMemDB()
	fmt.Fprintln(os.Stderr, "Opening app database")
	appDB, err := sdk.NewLevelDB("application", dataDir)
	if err != nil {
		return err
	}

	// TM DB
	// tmDB := dbm.NewMemDB()
	fmt.Fprintln(os.Stderr, "Opening state database")
	tmDB, err := sdk.NewLevelDB("state", dataDir)
	if err != nil {
		return err
	}

	// Blockchain DB
	fmt.Fprintln(os.Stderr, "Opening blockstore database")
	bcDB, err := sdk.NewLevelDB("blockstore", dataDir)
	if err != nil {
		return err
	}

	// Application
	fmt.Fprintln(os.Stderr, "Creating application")
	myapp := app.NewNCHApp(ctx.Logger, appDB, nil, true, uint(1))

	// Genesis
	var genDocPath = filepath.Join(configDir, "genesis.json")
	genDoc, err := tm.GenesisDocFromFile(genDocPath)
	if err != nil {
		return err
	}
	genState, err := tmsm.MakeGenesisState(genDoc)
	if err != nil {
		return err
	}
	// tmsm.SaveState(tmDB, genState)

	cc := proxy.NewLocalClientCreator(myapp)
	proxyApp := proxy.NewAppConns(cc)
	err = proxyApp.Start()
	if err != nil {
		return err
	}
	defer func() {
		_ = proxyApp.Stop()
	}()

	state := tmsm.LoadState(tmDB)
	if state.LastBlockHeight == 0 {
		// Send InitChain msg
		fmt.Fprintln(os.Stderr, "Sending InitChain msg")
		validators := tm.TM2PB.ValidatorUpdates(genState.Validators)
		csParams := tm.TM2PB.ConsensusParams(genDoc.ConsensusParams)
		req := abci.RequestInitChain{
			Time:            genDoc.GenesisTime,
			ChainId:         genDoc.ChainID,
			ConsensusParams: csParams,
			Validators:      validators,
			AppStateBytes:   genDoc.AppState,
		}
		res, err := proxyApp.Consensus().InitChainSync(req)
		if err != nil {
			return err
		}
		newValidatorz, err := tm.PB2TM.ValidatorUpdates(res.Validators)
		if err != nil {
			return err
		}
		newValidators := tm.NewValidatorSet(newValidatorz)

		// Take the genesis state.
		state = genState
		state.Validators = newValidators
		state.NextValidators = newValidators

		tmsm.SaveState(tmDB, state)
	}

	// Create executor
	fmt.Fprintln(os.Stderr, "Creating block executor")

	mempoolInstance := mempl.NewCListMempool(cfg.DefaultMempoolConfig(), proxyApp.Mempool(), 1)
	blockExec := tmsm.NewBlockExecutor(tmDB, ctx.Logger, proxyApp.Consensus(), mempoolInstance, tmsm.MockEvidencePool{})

	// Create block store
	fmt.Fprintln(os.Stderr, "Creating block store")
	blockStore := tmstore.NewBlockStore(bcDB)

	tz := []time.Duration{0, 0, 0}
	for i := state.LastBlockHeight + 1; to == -1 || i <= to; i++ {
		fmt.Fprintln(os.Stderr, "Running block ", i)
		t1 := time.Now()

		// Apply block
		fmt.Printf("loading and applying block %d\n", i)
		blockmeta := blockStore.LoadBlockMeta(i)
		if blockmeta == nil {
			fmt.Printf("Couldn't find block meta %d... done?\n", i)
			return nil
		}
		block := blockStore.LoadBlock(int64(i))
		if block == nil {
			return fmt.Errorf("couldn't find block %d", i)
		}

		t2 := time.Now()

		state, err = blockExec.ApplyBlock(state, blockmeta.BlockID, block)
		if err != nil {
			return err
		}

		t3 := time.Now()
		tz[0] += t2.Sub(t1)
		tz[1] += t3.Sub(t2)

		fmt.Fprintf(os.Stderr, "new app hash: %X\n", state.AppHash)
		fmt.Fprintln(os.Stderr, tz)
	}

	return nil
}
