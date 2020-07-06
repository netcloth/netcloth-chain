package vm

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/app/v0/vm/common"
	keep "github.com/netcloth/netcloth-chain/app/v0/vm/keeper"
	"github.com/netcloth/netcloth-chain/app/v0/vm/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func TestInvalidMsg(t *testing.T) {
	k := Keeper{}
	h := NewHandler(k)

	res, err := h(sdk.NewContext(nil, abci.Header{}, false, nil), sdk.NewTestMsg())
	require.NotNil(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "unrecognized vm message type"))
}

// contract code from: https://docs.netcloth.org/contracts/contract.html
func TestMsgContractCreateAndCall(t *testing.T) {
	initPower := int64(1000000)
	ctx, accountKeeper, vmKeeper, _ := keep.CreateTestInput(t, false, initPower)

	acc := accountKeeper.GetAccount(ctx, keep.Addrs[0])
	code := sdk.FromHex("608060405234801561001057600080fd5b506509184e72a0006000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610344806100696000396000f300608060405260043610610057576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806327e235e31461005c57806370a08231146100b3578063a9059cbb1461010a575b600080fd5b34801561006857600080fd5b5061009d600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610162565b6040518082815260200191505060405180910390f35b3480156100bf57600080fd5b506100f4600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061017a565b6040518082815260200191505060405180910390f35b610148600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506101c2565b604051808215151515815260200191505060405180910390f35b60006020528060005260406000206000915090505481565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b6000816000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541015151561021157600080fd5b816000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282540392505081905550816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a360019050929150505600a165627a7a7230582015481e18f5439ee76271037928d88d33cc7d7d4bf1e5e801b78db9e902f255560029")

	fmt.Printf("addr: %s, nonce: %d\n", acc.GetAddress().String(), acc.GetSequence())
	contractAddr := CreateAddress(acc.GetAddress(), acc.GetSequence())
	fmt.Printf("contract addr: %s\n", contractAddr.String())

	handler := NewHandler(vmKeeper)

	// test ContractCreate
	msgCreate := types.NewMsgContract(acc.GetAddress(), nil, code, sdk.NewInt64Coin(sdk.NativeTokenName, 0))
	require.NotNil(t, msgCreate)
	require.Equal(t, msgCreate.Route(), RouterKey)
	require.Equal(t, msgCreate.Type(), types.TypeMsgContractCreate)

	resCreate, err := handler(ctx, msgCreate)
	require.Nil(t, err)
	if len(resCreate.Log) > 0 {
		fmt.Println("logs: ", resCreate.Log)
	}
	require.NotNil(t, vmKeeper.StateDB.GetCode(contractAddr))

	// end blocker
	EndBlocker(ctx, vmKeeper)

	// test ContractCall
	msgCall := types.NewMsgContract(acc.GetAddress(), contractAddr, common.FromHex("a9059cbb0000000000000000000000005376329591cde25497d29de88ec553229ad10a610000000000000000000000000000000000000000000000000000000000000064"), sdk.NewInt64Coin(sdk.NativeTokenName, 10))
	require.NotNil(t, msgCall)
	require.Equal(t, msgCall.Route(), RouterKey)
	require.Equal(t, msgCall.Type(), types.TypeMsgContractCall)

	resCall, err := handler(ctx, msgCall)
	require.Nil(t, err)
	if len(resCall.Log) > 0 {
		ctx.Logger().Debug(fmt.Sprintf("event logs: %v", resCall.Log))
	}

	// get event logs
	logs := vmKeeper.GetLogs(ctx, sdk.BytesToHash(tmhash.Sum(ctx.TxBytes())))
	d, err := json.Marshal(logs)
	require.Nil(t, err)
	ctx.Logger().Debug(fmt.Sprintf("get event logs: %s", string(d)))
}
