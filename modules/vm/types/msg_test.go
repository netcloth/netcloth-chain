package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/netcloth/netcloth-chain/types"
)

func TestMsgContractCreate(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	coin123 := sdk.NewInt64Coin(sdk.NativeTokenName, 123)
	coin0 := sdk.NewInt64Coin(sdk.NativeTokenName, 0)
	coin123eth := sdk.NewInt64Coin("eth", 123)
	coin0eth := sdk.NewInt64Coin("eth", 0)

	payload := []byte("payload")
	payloadEmpty := []byte("")

	var emptyAddr sdk.AccAddress

	cases := []struct {
		valid bool
		tx    MsgContractCreate
	}{
		{true, NewMsgContractCreate(addr1, coin123, payload)},
		{true, NewMsgContractCreate(addr1, coin123eth, payload)},
		{true, NewMsgContractCreate(addr1, coin0, payload)},
		{true, NewMsgContractCreate(addr1, coin0eth, payload)},

		{false, NewMsgContractCreate(addr1, coin123, nil)},
		{false, NewMsgContractCreate(addr1, coin123, payloadEmpty)},
		{false, NewMsgContractCreate(emptyAddr, coin123, payload)},
	}

	for _, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err)
		} else {
			require.NotNil(t, err)
		}
	}
}

func TestMsgContractCreateGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	coin123 := sdk.NewInt64Coin(sdk.NativeTokenName, 123)
	payload := []byte("payload")
	var msg = NewMsgContractCreate(addr1, coin123, payload)
	res := msg.GetSignBytes()

	expected := `{"type":"nch/MsgContractCreate","value":{"amount":{"amount":"123","denom":"pnch"},"code":"cGF5bG9hZA==","from":"nch1veex7mg3k0xqr"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgContractCallGetSigners(t *testing.T) {
	var msg = NewMsgContractCreate(sdk.AccAddress([]byte("from")), sdk.NewInt64Coin(sdk.NativeTokenName, 123), []byte("payload"))
	res := msg.GetSigners()

	require.Equal(t, fmt.Sprintf("%v", res), "[66726F6D]")
}

func TestMsgContractCreateRoute(t *testing.T) {
	// construct a MsgContractCreate
	addr1 := sdk.AccAddress([]byte("from"))
	coin := sdk.NewInt64Coin(sdk.NativeTokenName, 10)
	code := []byte("contract code")
	msg := NewMsgContractCreate(addr1, coin, code)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), TypeMsgContractCreate)
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
