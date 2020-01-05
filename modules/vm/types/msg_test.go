package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/netcloth/netcloth-chain/types"
)

func TestMsgContractCreate(t *testing.T) {
	// construct a MsgContractCreate
	addr1 := sdk.AccAddress([]byte("from"))
	coin := sdk.NewInt64Coin(sdk.NativeTokenName, 10)
	code := []byte("contract code")
	msg := NewMsgContractCreate(addr1, coin, code)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), TypeMsgContractCreate)
}

func TestMsgContractCreateRoute(t *testing.T) {

}

func TestMsgContractCall(t *testing.T) {

}

func TestMsgContractCallRoute(t *testing.T) {
	// construct a MsgContractCall
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))
	coin := sdk.NewInt64Coin(sdk.NativeTokenName, 10)
	payload := []byte("payload")

	msg := NewMsgContractCall(addr1, addr2, coin, payload)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), TypeMsgContractCall)
}
