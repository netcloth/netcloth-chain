package mock

import (
	"github.com/netcloth/netcloth-chain/app"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"os"
)

func SetupMockApp() (log.Logger, dbm.DB, *app.NCHApp) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "app/mock")
	db := dbm.NewMemDB()

	return logger, db, app.NewNCHApp(logger, db, os.Stdout, true, 0)
}
