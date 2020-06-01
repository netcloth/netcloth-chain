package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/netcloth/netcloth-chain/types"
)

func TestMsgContract(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))
	coin123 := sdk.NewInt64Coin(sdk.NativeTokenName, 123)
	coin0 := sdk.NewInt64Coin(sdk.NativeTokenName, 0)
	coin123eth := sdk.NewInt64Coin("eth", 123)
	coin0eth := sdk.NewInt64Coin("eth", 0)
	coinNegative := sdk.Coin{sdk.NativeTokenName, sdk.NewInt(-123)}

	payload := []byte("payload")
	payloadEmpty := []byte("")

	var emptyAddr sdk.AccAddress

	cases := []struct {
		valid bool
		tx    MsgContract
	}{
		// create
		{true, NewMsgContract(addr1, nil, payload, coin123)},

		{false, NewMsgContract(addr1, nil, payload, coin123eth)},
		{true, NewMsgContract(addr1, nil, payload, coin0)},
		{false, NewMsgContract(addr1, nil, payload, coin0eth)},

		{false, NewMsgContract(addr1, nil, payload, coinNegative)},
		{false, NewMsgContract(addr1, nil, nil, coin123)},
		{false, NewMsgContract(addr1, nil, payloadEmpty, coin123)},
		{false, NewMsgContract(emptyAddr, nil, payload, coin123)},

		// call
		{true, NewMsgContract(addr1, addr2, payload, coin123)},
		{false, NewMsgContract(addr1, addr2, payload, coin123eth)},
		{true, NewMsgContract(addr1, addr2, payload, coin0)},
		{false, NewMsgContract(addr1, addr2, payload, coin0eth)},

		{false, NewMsgContract(addr1, addr2, payload, coinNegative)},
		{false, NewMsgContract(addr1, addr2, nil, coin123)},
		{false, NewMsgContract(addr1, addr2, payloadEmpty, coin123)},
		{false, NewMsgContract(emptyAddr, addr2, payload, coin123)},
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

func TestMsgContractGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	coin123 := sdk.NewInt64Coin(sdk.NativeTokenName, 123)
	payload := []byte("payload")
	msg := NewMsgContract(addr1, nil, payload, coin123)

	// with "to" nil
	res := msg.GetSignBytes()
	expected := `{"type":"nch/MsgContract","value":{"amount":{"amount":"123","denom":"pnch"},"from":"nch1veex7mg3k0xqr","payload":"7061796c6f6164","to":""}}`
	require.Equal(t, expected, string(res))

	// with "to" empty
	var addrEmpty sdk.AccAddress
	msg = NewMsgContract(addr1, addrEmpty, payload, coin123)
	res = msg.GetSignBytes()
	expected = `{"type":"nch/MsgContract","value":{"amount":{"amount":"123","denom":"pnch"},"from":"nch1veex7mg3k0xqr","payload":"7061796c6f6164","to":""}}`
	require.Equal(t, expected, string(res))

	// with "to"
	addr2 := sdk.AccAddress([]byte("to"))
	msg = NewMsgContract(addr1, addr2, payload, coin123)
	res = msg.GetSignBytes()
	expected = `{"type":"nch/MsgContract","value":{"amount":{"amount":"123","denom":"pnch"},"from":"nch1veex7mg3k0xqr","payload":"7061796c6f6164","to":"nch1w3hsls558e"}}`
	require.Equal(t, expected, string(res))

}

func TestMsgContractGetSigners(t *testing.T) {
	var msg = NewMsgContract(sdk.AccAddress([]byte("from")), nil, []byte("payload"), sdk.NewInt64Coin(sdk.NativeTokenName, 123))
	res := msg.GetSigners()

	require.Equal(t, fmt.Sprintf("%v", res), "[66726F6D]")
}

func TestMsgContractRoute(t *testing.T) {
	// construct a MsgContract
	addr1 := sdk.AccAddress([]byte("from"))
	coin := sdk.NewInt64Coin(sdk.NativeTokenName, 10)
	payload := []byte("payload")
	msg := NewMsgContract(addr1, nil, payload, coin)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), TypeMsgContract)

	// construct a MsgContract
	addr2 := sdk.AccAddress([]byte("to"))
	msg = NewMsgContract(addr1, addr2, payload, coin)
	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), TypeMsgContract)

}
