package vm

import (
	"fmt"
	"math/big"
	"os"

	"github.com/netcloth/netcloth-chain/modules/vm/types"

	sdk "github.com/netcloth/netcloth-chain/types"
)

// StateTransition defines data to transitionDB in vm
type StateTransition struct {
	Sender    sdk.AccAddress
	Price     sdk.Int
	GasLimit  uint64
	Recipient sdk.AccAddress
	Amount    sdk.Int
	Payload   []byte
	CSDB      *types.CommitStateDB
}

func (st StateTransition) CanTransfer(acc sdk.AccAddress, amount *big.Int) bool {
	return st.CSDB.GetBalance(acc).Cmp(amount) >= 0
}

func (st StateTransition) Transfer(from, to sdk.AccAddress, amount *big.Int) sdk.Error {
	return st.CSDB.BK.SendCoins(st.CSDB.Ctx, from, to, sdk.NewCoins(sdk.NewCoin("unch", sdk.NewInt(amount.Int64()))))
}

func (st StateTransition) GetHash(uint64) sdk.Hash {
	return sdk.Hash{}
}

func (st StateTransition) TransitionCSDB(ctx sdk.Context) (*big.Int, sdk.Result) {

	evmCtx := Context{
		CanTransfer: st.CanTransfer,
		Transfer:    st.Transfer,
		GetHash:     st.GetHash,

		Origin:   st.Sender,
		GasPrice: st.Price.BigInt(),

		CoinBase:    ctx.BlockHeader().ProposerAddress, // TODO: should be proposer account address
		GasLimit:    st.GasLimit,
		BlockNumber: sdk.NewInt(ctx.BlockHeader().Height).BigInt(),
	}

	cfg := Config{}

	evm := NewEVM(evmCtx, *st.CSDB, cfg)

	var (
		ret         []byte
		leftOverGas uint64
		addr        sdk.AccAddress
		err         sdk.Error
	)

	if st.Recipient == nil {
		ret, addr, leftOverGas, err = evm.Create(st.Sender, st.Payload, 100000000, st.Amount.BigInt())
		fmt.Fprint(os.Stderr, fmt.Sprintf("contractAddr = %s, leftOverGas = %v, err = %v\n", addr, leftOverGas, err))

		if err != nil {
			return nil, sdk.ErrInternal("contract deploy err").Result()
		}
	} else {
		ret, leftOverGas, err = evm.Call(st.Sender, st.Recipient, st.Payload, 1000000000, st.Amount.BigInt())
		fmt.Fprint(os.Stderr, fmt.Sprintf("ret = %x, leftOverGas = %v, err = %v\n", ret, leftOverGas, err))
	}

	st.CSDB.Finalise(true)

	return nil, sdk.Result{Data: ret, GasUsed: st.GasLimit - leftOverGas}
}
