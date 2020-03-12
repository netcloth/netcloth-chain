package cipal

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/netcloth/netcloth-chain/types"
)

func TestInvalidMsg(t *testing.T) {
	k := Keeper{}
	h := NewHandler(k)

	res, err := h(sdk.NewContext(nil, abci.Header{}, false, nil), sdk.NewTestMsg())
	require.Nil(t, res)
	require.NotNil(t, err)
	require.True(t, strings.Contains(err.Error(), "unrecognized cipal message type"))
}
