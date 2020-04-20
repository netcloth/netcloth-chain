package tests

import (
	"github.com/netcloth/netcloth-chain/app/v0/gov"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/netcloth/netcloth-chain/types"
	sdkerrors "github.com/netcloth/netcloth-chain/types/errors"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestInvalidMsg(t *testing.T) {
	k := gov.Keeper{}
	h := gov.NewHandler(k)

	res, err := h(sdk.NewContext(nil, abci.Header{}, false, nil), sdk.NewTestMsg())
	require.Error(t, err)
	require.Nil(t, res)

	_, _, log := sdkerrors.ABCIInfo(err, false)
	require.True(t, strings.Contains(res.Log, "unrecognized gov message type"))
}
