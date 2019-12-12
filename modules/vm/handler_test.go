package vm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/modules/genaccounts"
	"github.com/netcloth/netcloth-chain/modules/vm/common"
	"github.com/netcloth/netcloth-chain/modules/vm/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func TestInvalidMsg(t *testing.T) {
	k := Keeper{}
	h := NewHandler(k)

	res := h(sdk.NewContext(nil, abci.Header{}, false, nil), sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized Msg type"))
}

func TestMsgContractCreateAndCall(t *testing.T) {
	fromAddr := newSdkAddress()
	coins := sdk.NewCoins(sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(10000000000)))
	genAcc := genaccounts.NewGenesisAccountRaw(fromAddr, coins, coins, 0, 1, "", "")
	code := sdk.FromHex("60806040526010600055600860015534801561001a57600080fd5b506101748061002a6000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80633033413b146100515780635d33a27f1461006f578063a17a9e661461008d578063dac0eb07146100cf575b600080fd5b6100596100fd565b6040518082815260200191505060405180910390f35b610077610103565b6040518082815260200191505060405180910390f35b6100b9600480360360208110156100a357600080fd5b8101908080359060200190929190505050610109565b6040518082815260200191505060405180910390f35b6100fb600480360360208110156100e557600080fd5b8101908080359060200190929190505050610124565b005b60005481565b60015481565b60008160008082825401925050819055506000549050919050565b61012d81610109565b6001600082825401925050819055505056fea265627a7a72315820a78769f5b6dbc72a06c9d07bed3d8ff5830a8cd9963a84ec8821de089e9b1e5e64736f6c634300050d0032")

	fmt.Println(fmt.Sprintf("addr: %s, nonce: %d", fromAddr.String(), genAcc.Sequence))
	contractAddr := CreateAddress(fromAddr, genAcc.Sequence)
	fmt.Println(fmt.Sprintf("contract addr: %s", contractAddr.String()))

	k, ctx := setupTest()
	handler := NewHandler(k)

	// test MsgContractCreate
	msgCreate := types.NewMsgContractCreate(fromAddr, sdk.NewInt64Coin(sdk.NativeTokenName, 1000), code)
	require.NotNil(t, msgCreate)
	require.Equal(t, msgCreate.Route(), RouterKey)
	require.Equal(t, msgCreate.Type(), types.TypeMsgContractCreate)

	resCreate := handler(ctx, msgCreate)
	require.True(t, resCreate.IsOK())
	if len(resCreate.Log) > 0 {
		fmt.Println("logs: ", resCreate.Log)
	}
	require.NotNil(t, k.CSDB.GetCode(contractAddr))

	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	contractAcc := k.GetAccount(ctx, contractAddr)
	if contractAcc != nil {
		fmt.Println(contractAcc.GetCoins())
	}
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++")

	// test MsgContractCall
	msgCall := types.NewMsgContractCall(fromAddr, contractAddr, sdk.NewInt64Coin(sdk.NativeTokenName, 0), common.FromHex("3033413b"))
	require.NotNil(t, msgCall)
	require.Equal(t, msgCall.Route(), RouterKey)
	require.Equal(t, msgCall.Type(), types.TypeMsgContractCall)

	resCall := handler(ctx, msgCall)
	require.True(t, resCall.IsOK())
	if len(resCall.Log) > 0 {
		fmt.Println("logs: ", resCall.Log)
	}
}
