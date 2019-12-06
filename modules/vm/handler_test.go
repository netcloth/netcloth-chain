package vm

import (
	"strings"
	"testing"

	"github.com/netcloth/netcloth-chain/modules/vm/types"

	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/netcloth/netcloth-chain/types"
)

func newSdkAddress() sdk.AccAddress {
	tmpKey := secp256k1.GenPrivKey().PubKey()
	return sdk.AccAddress(tmpKey.Address().Bytes())
}

func TestInvalidMsg(t *testing.T) {
	k := Keeper{}
	h := NewHandler(k)

	res := h(sdk.NewContext(nil, abci.Header{}, false, nil), sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized Msg type"))
}

func TestMsgContractCreate(t *testing.T) {
	fromAddr := newSdkAddress()
	amount := sdk.NewInt64Coin(sdk.NativeTokenName, 150)
	code := []byte("xxxx")

	msg := types.NewMsgContractCreate(fromAddr, amount, code)
	require.NotNil(t, msg)
	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), types.TypeMsgContractCreate)

}

func TestMsgContractCall(t *testing.T) {
}
