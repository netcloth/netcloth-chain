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
	GasLimit  sdk.Int
	Recipient sdk.AccAddress
	Amount    sdk.Int
	Payload   []byte
	CSDB      *types.CommitStateDB
}

func (st StateTransition) CanTransfer(acc sdk.AccAddress, amount *big.Int) bool {
	return true
}

func (st StateTransition) Transfer(from, to sdk.AccAddress, amount *big.Int) {
	//st.CSDB.BK.SendCoins(st.CSDB.Ctx, from, to, sdk.NewCoins(sdk.NewCoin("unch", sdk.NewInt(amount.Int64()))))
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

		CoinBase:    ctx.BlockHeader().ProposerAddress,
		GasLimit:    uint64(st.GasLimit.Int64()),
		BlockNumber: sdk.NewInt(ctx.BlockHeader().Height).BigInt(),
	}

	cfg := Config{}

	evm := NewEVM(evmCtx, *st.CSDB, cfg)

	ret, contractAddr, leftOverGas, e := evm.Create(st.Sender, st.Payload, 100000000, sdk.NewInt(0).BigInt())
	fmt.Fprint(os.Stderr, fmt.Sprintf("ret = %v, contractAddr = %v, leftOverGas = %v, err = %v\n", ret, contractAddr, leftOverGas, e))
	if e != nil {
		return nil, sdk.ErrInternal("contract deploy err").Result()
	}

	return nil, sdk.Result{}
}
