package ante

import (
	abci "github.com/tendermint/tendermint/abci/types"

	authtypes "github.com/netcloth/netcloth-chain/modules/auth/types"
	"github.com/netcloth/netcloth-chain/simapp"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())

	return app, ctx
}
