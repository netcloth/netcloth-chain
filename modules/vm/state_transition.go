package vm

import (
	"fmt"
	"math/big"

	abci "github.com/tendermint/tendermint/abci/types"
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

func (st StateTransition) GetHashFn(header abci.Header) func() sdk.Hash {
	return func() sdk.Hash {
		var res = sdk.Hash{}
		blockID := header.GetLastBlockId()
		res.SetBytes(blockID.GetHash())
		return res
	}
}

func (st StateTransition) TransitionCSDB(ctx sdk.Context, vmParams *types.Params) (*big.Int, *sdk.Result, error) {
	evmCtx := Context{
		CanTransfer: st.CanTransfer,
		Transfer:    st.Transfer,
		GetHash:     st.GetHashFn(ctx.BlockHeader()),

		Origin: st.Sender,

		CoinBase:    ctx.BlockHeader().ProposerAddress, // validator consensus address, not account address
		Time:        sdk.NewInt(int64(ctx.BlockHeader().Time.Unix())).BigInt(),
		GasLimit:    st.GasLimit,
		BlockNumber: sdk.NewInt(ctx.BlockHeader().Height).BigInt(),
	}

	// This gas meter is set up to consume gas from gaskv during evm execution and be ignored
	currentGasMeter := ctx.GasMeter()
	csdb := st.StateDB.WithContext(ctx.WithGasMeter(sdk.NewInfiniteGasMeter())).WithTxHash(tmhash.Sum(ctx.TxBytes()))
	// Clear cache of accounts to handle changes outside of the EVM
	csdb.UpdateAccounts()

	cfg := Config{OpConstGasConfig: &vmParams.VMOpGasParams, CommonGasConfig: &vmParams.VMCommonGasParams}
	evm := NewEVM(evmCtx, csdb, cfg)

	var (
		ret         []byte
		leftOverGas uint64
		addr        sdk.AccAddress
		vmerr       error
	)

	if st.Recipient.Empty() {
		ret, addr, leftOverGas, vmerr = evm.Create(st.Sender, st.Payload, st.GasLimit, st.Amount.BigInt())
		ctx.Logger().Info(fmt.Sprintf("create contract, consumed gas = %v , leftOverGas = %v, err = %v\n", st.GasLimit-leftOverGas, leftOverGas, vmerr))
	} else {
		ret, leftOverGas, vmerr = evm.Call(st.Sender, st.Recipient, st.Payload, st.GasLimit, st.Amount.BigInt())

		if vmerr == ErrExecutionReverted {
			ctx.Logger().Info(fmt.Sprintf("VM error: revert. Reason provided by the contract: %v", string(ret[4:])))
		}

		ctx.Logger().Info(fmt.Sprintf("call contract, ret = %v, consumed gas = %v , leftOverGas = %v, err = %v\n", ret, st.GasLimit-leftOverGas, leftOverGas, vmerr))
	}

	ctx.WithGasMeter(currentGasMeter).GasMeter().ConsumeGas(st.GasLimit-leftOverGas, "EVM execution consumption")
	if vmerr != nil {
		return nil, &sdk.Result{Data: ret, GasUsed: st.GasLimit - leftOverGas}, vmerr
	}

	st.StateDB.Finalise(true)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeNewContract,
			sdk.NewAttribute(types.AttributeKeyAddress, addr.String()),
		),
	})

	return nil, &sdk.Result{Data: ret, GasUsed: st.GasLimit - leftOverGas}, nil
}

func DoStateTransition(ctx sdk.Context, msg types.MsgContract, k Keeper, gasLimit uint64, readonly bool) (*big.Int, *sdk.Result, error) {
	st := StateTransition{
		Sender:    msg.From,
		Recipient: msg.To,
		GasLimit:  gasLimit,
		Payload:   msg.Payload,
		Amount:    msg.Amount.Amount,
		StateDB:   k.StateDB.WithContext(ctx),
	}

	if readonly {
		st.StateDB = types.NewStateDB(k.StateDB).WithContext(ctx)
	}

	params := k.GetParams(ctx)
	return st.TransitionCSDB(ctx, &params)
}
