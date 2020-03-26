package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/netcloth/netcloth-chain/types"
)

var (
	addr1        = sdk.AccAddress([]byte("from"))
	moniker      = "moniker"
	website      = "website"
	details      = "details"
	extension    = ""
	endpoints, _ = EndpointsFromString("1|http://1.1.1.1,3|http://2.2.2.2", ",", "|")
	bond         = sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(sdk.NativeTokenFraction))
)

func TestMsgIPALNodeClaimRoute(t *testing.T) {
	var msg = NewMsgIPALNodeClaim(addr1, moniker, website, details, extension, endpoints, bond)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), TypeMsgIPALNodeClaim)
}

func TestMsgIPALNodeClaimValidation(t *testing.T) {
	var emptyAddr sdk.AccAddress

	var negativeCoin = sdk.Coin{Denom: sdk.NativeTokenName, Amount: sdk.NewInt(int64(-1))}
	var xnchCoin = sdk.NewCoin("xnch", sdk.NewInt(sdk.NativeTokenFraction))

	// duplicate endpoints
	ep := Endpoint{1, "http://1.1.1.1"}
	var dupEndPoints = Endpoints{}
	dupEndPoints = append(dupEndPoints, ep)
	dupEndPoints = append(dupEndPoints, ep)

	cases := []struct {
		valid bool
		tx    MsgIPALNodeClaim
	}{
		{true, NewMsgIPALNodeClaim(addr1, moniker, website, details, extension, endpoints, bond)}, // valid

		{false, NewMsgIPALNodeClaim(emptyAddr, moniker, website, details, extension, endpoints, bond)}, // empty from addr
		{false, NewMsgIPALNodeClaim(addr1, "", website, details, extension, endpoints, bond)},          // empty moniker
		{false, NewMsgIPALNodeClaim(addr1, "", website, details, extension, endpoints, xnchCoin)},      //  other bond coins
		{false, NewMsgIPALNodeClaim(addr1, "", website, details, extension, endpoints, negativeCoin)},  //  negative coins
		{false, NewMsgIPALNodeClaim(addr1, moniker, website, details, extension, Endpoints{}, bond)},   // empty endpoints
		{false, NewMsgIPALNodeClaim(addr1, moniker, website, details, extension, dupEndPoints, bond)},  // duplicate endpoints
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

func TestMsgIPALNodeClaimGetSignBytes(t *testing.T) {
	var msg = NewMsgIPALNodeClaim(addr1, moniker, website, details, extension, endpoints, bond)
	res := msg.GetSignBytes()

	expected := `{"type":"nch/IPALClaim","value":{"bond":{"amount":"1000000000000","denom":"pnch"},"details":"details","endpoints":[{"endpoint":"http://1.1.1.1","type":"1"},{"endpoint":"http://2.2.2.2","type":"3"}],"extension":"","moniker":"moniker","operator_address":"nch1veex7mg3k0xqr","website":"website"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgIPALNodeClaimGetSigners(t *testing.T) {
	var msg = NewMsgIPALNodeClaim(sdk.AccAddress([]byte("input1")), moniker, website, details, extension, endpoints, bond)
	res := msg.GetSigners()

	require.Equal(t, fmt.Sprintf("%v", res), "[696E70757431]")
}
