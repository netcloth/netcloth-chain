package vm

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/netcloth/netcloth-chain/modules/vm/common"

	"github.com/netcloth/netcloth-chain/modules/vm/types"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/netcloth/netcloth-chain/types"

	tmcrypto "github.com/tendermint/tendermint/crypto"
)

var emptyCodeHash = crypto.Sha256(nil)

func run(evm *EVM, contract *Contract, input []byte, readOnly bool) ([]byte, error) {
	//if contract.CodeAddr != nil {
	//	precompiles := PrecompiledContractsHomestead
	//	if evm.chainRules.IsByzantium {
	//		precompiles = PrecompiledContractsByzantium
	//	}
	//	if evm.chainRules.IsIstanbul {
	//		precompiles = PrecompiledContractsIstanbul
	//	}
	//	if p := precompiles[*contract.CodeAddr]; p != nil {
	//		return RunPrecompiledContract(p, input, contract)
	//	}
	//}
	for _, interpreter := range evm.interpreters {
		if interpreter.CanRun(contract.Code) {
			if evm.interpreter != interpreter {
				// Ensure that the interpreter pointer is set back
				// to its current value upon return.
				defer func(i Interpreter) {
					evm.interpreter = i
				}(evm.interpreter)
				evm.interpreter = interpreter
			}
			return interpreter.Run(contract, input, readOnly)
		}
	}
	return nil, ErrNoCompatibleInterpreter
}

type codeAndHash struct {
	code []byte
	hash sdk.Hash
}

func (c *codeAndHash) Hash() sdk.Hash {
	if c.hash == (sdk.Hash{}) {
		copy(c.hash[:], crypto.Sha256(c.code))
	}
	return c.hash
}

// Context provides the VM with auxiliary information.
// Once provided it shouldn't be modified
type Context struct {
	ctx sdk.Context
	// Msg information
	Origin   sdk.AccAddress
	GasPrice *big.Int

	// Block information
	CoinBase    sdk.AccAddress
	GasLimit    uint64
	BlockNumber *big.Int
	Time        *big.Int
}

func (ctx *EVM) CanTransfer(sdk.AccAddress, *big.Int) bool {
	//TODO
	//balanceEnough := false
	//coins := acc.GetCoins()
	//for _, coin := range coins {
	//	if coin.IsGTE(msg.Amount) {
	//		balanceEnough = true
	//	}
	//}
	return true
}

func (ctx *EVM) Transfer(from, to sdk.AccAddress, value *big.Int) {
	//TODO
}

func (ctx *EVM) GetHash(uint64) sdk.Hash {
	return sdk.Hash{} //TODO
}

func NewEVMContext(ctx sdk.Context, from sdk.AccAddress) Context {
	return Context{
		ctx:         ctx,
		Origin:      from,
		BlockNumber: new(big.Int).SetInt64(ctx.BlockHeader().Height),
		Time:        new(big.Int).SetInt64(ctx.BlockHeader().Time.Unix()),
		GasLimit:    100000000,                      //TODO fixme
		GasPrice:    new(big.Int).SetInt64(1000000), //TODO fix
	}
}

type EVM struct {
	Context

	// StateDB gives access to the underlying state
	StateDB *CommitStateDB

	// depth is the current call stack
	depth int

	chainConfig *ChainConfig

	// virtual machine configuration options used to initialise the vm
	vmConfig Config

	interpreters []Interpreter
	interpreter  Interpreter

	// abort is used to abort the EVM calling operations
	// NOTE: must be set atomically
	abort int32

	// callGasTemp holds the gas available for the current call. This is needed because the
	// available gas is calculated in gasCall* according to the 63/64 rule and later
	// applied in opCall*.
	callGasTemp uint64
}

func NewEVMInterpreter(evm *EVM, cfg Config) *EVMInterpreter {
	if !cfg.JumpTable[STOP].valid {
		cfg.JumpTable = istanbulInstructionSet
	}

	return &EVMInterpreter{
		evm: evm,
		cfg: cfg,
	}
}

func NewEVM(ctx Context, statedb CommitStateDB, vmConfig Config) *EVM {
	evm := &EVM{
		Context:      ctx,
		StateDB:      &statedb,
		vmConfig:     vmConfig,
		interpreters: make([]Interpreter, 0, 1),
	}

	evm.interpreters = append(evm.interpreters, NewEVMInterpreter(evm, vmConfig))
	evm.interpreter = evm.interpreters[0]
	return evm
}

// Interpreter returns the current interpreter
func (evm *EVM) Interpreter() Interpreter {
	return evm.interpreter
}

// Create creates a new contract using code as deployment code
func (evm *EVM) Create(caller sdk.AccAddress, code []byte, gas uint64, value *big.Int) (ret []byte, contractAddr sdk.AccAddress, leftOverGas uint64, err error) {
	return evm.create(caller, &codeAndHash{code: code}, gas, value)
}

func (evm *EVM) create(caller sdk.AccAddress, codeAndHash *codeAndHash, gas uint64, value *big.Int) ([]byte, sdk.AccAddress, uint64, error) {
	// Depth check execution. Fail if we're trying to execute above the limit
	if evm.depth > int(types.CallCreateDepth) {
		return nil, sdk.AccAddress{}, gas, ErrDepth
	}

	acc := evm.StateDB.AK.GetAccount(evm.Context.ctx, caller)
	if acc == nil {
		return nil, nil, 0, errors.New(fmt.Sprintf("account %s does not exist", caller.String()))
	}

	contractAddr := common.CreateAddress(caller, acc.GetSequence())
	fmt.Fprintf(os.Stderr, fmt.Sprintf("contractAddr = %v\n", contractAddr.String()))
	contractAcc := evm.StateDB.AK.GetAccount(evm.Context.ctx, contractAddr)
	if contractAcc != nil {
		return nil, nil, 0, errors.New(fmt.Sprintf("contract %s existed", contractAddr.String()))
	}

	if !evm.CanTransfer(caller, value) {
		return nil, nil, 0, errors.New(fmt.Sprintf("balace not enouth"))
	}

	codeHash := tmcrypto.Sha256(codeAndHash.code)

	// create account
	contractAcc = evm.StateDB.AK.NewAccountWithAddress(evm.Context.ctx, contractAddr.Bytes())
	contractAcc.SetCodeHash(codeHash)
	evm.StateDB.AK.SetAccount(evm.Context.ctx, contractAcc)

	// transfer
	evm.StateDB.BK.SendCoins(evm.Context.ctx, caller, contractAddr.Bytes(), sdk.NewCoins(sdk.NewCoin("unch", sdk.NewInt(value.Int64()))))

	// store code
	_, found := evm.StateDB.VK.GetContractCode(evm.ctx, codeHash)
	if !found {
		evm.StateDB.VK.SetContractCode(evm.ctx, codeHash, codeAndHash.code)
	}

	//if !evm.CanTransfer(caller, value) {
	//	return nil, sdk.AccAddress{}, gas, ErrInsufficientBalance
	//}

	//nonce := evm.StateDB.GetNonce(caller)
	//evm.StateDB.SetNonce(caller, nonce+1)

	// Ensure there's no existing contract already at the designated address
	//contractHash := evm.StateDB.GetCodeHash(caller.Address())
	//if evm.StateDB.GetNonce(address) != 0 || (contractHash != (sdk.Hash{})) {
	//	return nil, sdk.AccAddress{}, 0, ErrContractAddressCollision
	//}

	// Create a new account on the state
	//snapshot := evm.StateDB.Snapshot()
	//evm.StateDB.CreateAccount(address)
	//evm.StateDB.SetNonce(address, 1)
	//evm.Transfer(caller, address, value)

	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := NewContract(caller, contractAddr, value, gas)
	contract.SetCodeOptionalHash(&contractAddr, codeAndHash)

	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, contractAddr, gas, nil
	}

	//start := time.Now()
	fmt.Fprintf(os.Stderr, fmt.Sprintf("11111111111111111111\n"))
	ret, err := run(evm, contract, nil, false)
	fmt.Fprintf(os.Stderr, fmt.Sprintf("22222222222222, err = %v\n", err))

	maxCodeSizeExceeded := len(ret) > MaxCodeSize
	if err == nil && !maxCodeSizeExceeded {
		createGas := uint64(len(ret)) * types.CreateAccountGas
		if contract.UseGas(createGas) {
			evm.StateDB.SetCode(contractAddr, ret)
		} else {
			err = ErrCodeStoreOutOfGas
		}
	}

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if maxCodeSizeExceeded || (err != nil && (err != ErrCodeStoreOutOfGas)) {
		//evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	// Assign err if contract code size exceeds the max while the err is still empty.
	if maxCodeSizeExceeded && err == nil {
		err = errMaxCodeSizeExceeded
	}
	if evm.vmConfig.Debug && evm.depth == 0 {
		//evm.vmConfig.Tracer.CaptureEnd(ret, gas-contract.Gas, time.Since(start), err)
	}
	return ret, contractAddr, contract.Gas, err
}

// Create2 creates a new contract using code as deployment code.
//
// The different between Create2 with Create is Create2 uses sha3(0xff ++ msg.sender ++ salt ++ sha3(init_code))[12:]
// instead of the usual sender-and-nonce-hash as the address where the contract is initialized at.
func (evm *EVM) Create2(caller ContractRef, code []byte, gas uint64, endowment *big.Int, salt *big.Int) (ret []byte, contractAddr sdk.AccAddress, leftOverGas uint64, err error) {
	// TODO
	return
}

func (evm *EVM) Call(caller ContractRef, addr sdk.AccAddress, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}

	if evm.depth > int(CallCreateDepth) {
		return nil, gas, ErrDepth
	}

	if !evm.CanTransfer(caller.Address(), value) {
		return nil, gas, ErrInsufficientBalance
	}

	var (
		to       = AccountRef(addr)
		snapshot = evm.StateDB.Snapshot()
	)
	if !evm.StateDB.Exist(addr) {
		precompiles := PrecompiledContractsIstanbul

		if precompiles[addr.String()] == nil && value.Sign() == 0 {
			// Calling a non existing account, don't do anything, but ping the tracer
			if evm.vmConfig.Debug && evm.depth == 0 {
				evm.vmConfig.Tracer.CaptureStart(caller.Address(), addr, false, input, gas, value)
				evm.vmConfig.Tracer.CaptureEnd(ret, 0, 0, nil)
			}
			return nil, gas, nil
		}
		evm.StateDB.CreateAccount(addr)
	}
	evm.Transfer(caller.Address(), to.Address(), value)

	contract := NewContract(caller, to, value, gas)
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	start := time.Now()
	if evm.vmConfig.Debug && evm.depth == 0 {
		evm.vmConfig.Tracer.CaptureStart(caller.Address(), addr, false, input, gas, value)

		defer func() {
			evm.vmConfig.Tracer.CaptureEnd(ret, gas-contract.Gas, time.Since(start), err)
		}()
	}
	ret, err = run(evm, contract, input, false)

	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

func (evm *EVM) CallCode(caller ContractRef, addr sdk.AccAddress, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	// TODO
	return
}

func (evm *EVM) DelegateCall(caller ContractRef, addr sdk.AccAddress, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	// TODO
	return
}

func (evm *EVM) StaticCall(caller ContractRef, addr sdk.AccAddress, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	// TODO
	return
}
