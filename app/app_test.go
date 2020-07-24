package app

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	sdk "github.com/netcloth/netcloth-chain/types"
	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmsm "github.com/tendermint/tendermint/state"
	tm "github.com/tendermint/tendermint/types"
	db "github.com/tendermint/tm-db"
)

var (
	// the genesis file in unittest/ should be modified with this
	totalModuleNum = 16
)

func TestExport(t *testing.T) {
	db := db.NewMemDB()
	app := NewNCHApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)

	// make  simulation of abci.RequestInitChain
	genDoc, err := tm.GenesisDocFromFile("./genesis/genesis.json")
	require.NoError(t, err)
	genState, err := tmsm.MakeGenesisState(genDoc)
	require.NoError(t, err)

	validators := tm.TM2PB.ValidatorUpdates(genState.Validators)
	csParams := tm.TM2PB.ConsensusParams(genDoc.ConsensusParams)
	initChainRequest := abci.RequestInitChain{
		Time:            genDoc.GenesisTime,
		ChainId:         genDoc.ChainID,
		ConsensusParams: csParams,
		Validators:      validators,
		AppStateBytes:   genDoc.AppState,
	}

	// init chain
	require.NotPanics(t, func() {
		app.InitChain(initChainRequest)
	})

	// abci begin block
	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 1}})

	// abci end block
	app.EndBlock(abci.RequestEndBlock{Height: 1})

	// abci commit
	app.Commit()

	// block height should turn to 1
	require.Equal(t, int64(1), app.LastBlockHeight())

	// export the state of the latest height
	// situation 1: without jail white list
	appStateBytes, vals, err := app.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")

	var appState map[string]json.RawMessage
	app.Engine.GetCurrentProtocol().GetCodec().MustUnmarshalJSON(appStateBytes, &appState)
	require.Equal(t, totalModuleNum, len(appState))
	require.Equal(t, 1, len(vals))

	// situation 2: with jail white list
	jailWhiteList := []string{"nchvaloper1surzkw6dedqa29ntgdtfwar73y0u3v3k5fglqp"}
	_, _, err = app.ExportAppStateAndValidators(true, jailWhiteList)
	require.NoError(t, err)

	require.Equal(t, totalModuleNum, len(appState))

	// situation 3: with wrong format jail white list
	jailWhiteList = []string{"10q0rk5qnyag7wfvvt7rtphlw589m7frs863s3m"}

	require.Panics(t, func() {
		_, _, _ = app.ExportAppStateAndValidators(true, jailWhiteList)
	})

	// situation 4 : validator in the jail white list is not existed in the skakingKeeper
	jailWhiteList = []string{"nchvaloper1cmq8t3cgkwusdc236pjqn58cs60er83ua6dqh3"}
	require.Panics(t, func() {
		_, _, _ = app.ExportAppStateAndValidators(true, jailWhiteList)
	})

	///////////////////// test postEndBloker /////////////////////
	// situation 1
	testInput := &abci.ResponseEndBlock{}
	event1 := abci.Event{
		Type: "test",
		Attributes: []cmn.KVPair{
			{Key: []byte("key1"), Value: []byte("value1")},
		},
	}
	event2 := abci.Event{
		Type: sdk.AppVersionEvent,
		Attributes: []cmn.KVPair{
			{Key: []byte(sdk.AppVersionEvent), Value: []byte(strconv.FormatUint(1024, 10))},
		},
	}
	testInput.Events = append(testInput.Events, event1, event2)
	require.NotPanics(t, func() {
		app.postEndBlocker(testInput)
	})

	// situation 2
	testInput.Events = testInput.Events[:1]
	require.NotPanics(t, func() {
		app.postEndBlocker(testInput)
	})

	// situation 3
	testInput.Events = []abci.Event{
		{
			Type: sdk.AppVersionEvent,
			Attributes: []cmn.KVPair{
				{Key: []byte(sdk.AppVersionEvent), Value: []byte("parse error")},
			},
		},
	}
	require.NotPanics(t, func() {
		app.postEndBlocker(testInput)
	})

	// situation 4
	testInput.Events = []abci.Event{
		{
			Type: sdk.AppVersionEvent,
			Attributes: []cmn.KVPair{
				{Key: []byte(sdk.AppVersionEvent), Value: []byte(strconv.FormatUint(0, 10))},
			},
		},
	}

	require.NotPanics(t, func() {
		app.postEndBlocker(testInput)
	})
}
