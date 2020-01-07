package vm

import (
	"fmt"
	"math/big"
	"os"

	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/netcloth/netcloth-chain/modules/vm/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// StateTransition defines data to transitionDB in vm
type StateTransition struct {
	Sender    sdk.AccAddress
	GasLimit  uint64
	Recipient sdk.AccAddress
	Amount    sdk.Int
	Payload   []byte
	StateDB   *types.CommitStateDB
}

func (st StateTransition) CanTransfer(acc sdk.AccAddress, amount *big.Int) bool {
	return st.StateDB.GetBalance(acc).Cmp(amount) >= 0
}

func (st StateTransition) Transfer(from, to sdk.AccAddress, amount *big.Int) {
	st.StateDB.SubBalance(from, amount)
	st.StateDB.AddBalance(to, amount)
}

func (st StateTransition) GetHash(uint64) sdk.Hash {
	return sdk.Hash{}
}

// func returns: (ret []byte, usedGas uint64, failed bool, sdk.Result)
func (st StateTransition) TransitionCSDB(ctx sdk.Context, constGasConfig *[256]uint64, vmCommonGasConfig *types.VMCommonGasParams) (*big.Int, sdk.Result) {
	evmCtx := Context{
		CanTransfer: st.CanTransfer,
		Transfer:    st.Transfer,
		GetHash:     st.GetHash,

		Origin: st.Sender,

		CoinBase:    ctx.BlockHeader().ProposerAddress, // TODO: should be proposer account address
		GasLimit:    st.GasLimit,
		BlockNumber: sdk.NewInt(ctx.BlockHeader().Height).BigInt(),
	}

	cfg := Config{OpConstGasConfig: constGasConfig, CommonGasConfig: vmCommonGasConfig}

	evm := NewEVM(evmCtx, st.StateDB.WithTxHash(tmhash.Sum(ctx.TxBytes())), cfg)

	var (
		ret         []byte
		leftOverGas uint64
		addr        sdk.AccAddress
		vmerr       sdk.Error
	)

	if st.Recipient == nil {
		ret, addr, leftOverGas, vmerr = evm.Create(st.Sender, st.Payload, st.GasLimit, st.Amount.BigInt())
		fmt.Fprint(os.Stderr, "\n\n+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
		fmt.Fprint(os.Stderr, "+                                                                             +\n")
		fmt.Fprint(os.Stderr, fmt.Sprintf("+         contractAddr = %s           +\n", addr))
		fmt.Fprint(os.Stderr, "+                                                                             +\n")
		fmt.Fprint(os.Stderr, "+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n\n")
	} else {
		ret, leftOverGas, vmerr = evm.Call(st.Sender, st.Recipient, st.Payload, st.GasLimit, st.Amount.BigInt())
	}

	fmt.Fprint(os.Stderr, fmt.Sprintf("ret = %x, \nconsumed gas = %v , leftOverGas = %v, err = %v\n", ret, st.GasLimit-leftOverGas, leftOverGas, vmerr))

	ctx.GasMeter().ConsumeGas(st.GasLimit-leftOverGas, "EVM execution consumption")
	if vmerr != nil {
		return nil, vmerr.Result()
	}


	st.StateDB.Finalise(true)
	return nil, sdk.Result{Data: ret, GasUsed: st.GasLimit - leftOverGas}
}
