package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/netcloth/netcloth-chain/types"
)

func TestMsgUnjailGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress("abcd")
	msg := NewMsgUnjail(sdk.ValAddress(addr))
	bytes := msg.GetSignBytes()
	require.Equal(
		t,
		`{"type":"nch/MsgUnjail","value":{"address":"nchvaloper1v93xxeqkj7uhg"}}`,
		string(bytes),
	)
}
