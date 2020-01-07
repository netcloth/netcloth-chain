package vm

import (
	"fmt"
	"hash"
	"os"
	"sync/atomic"

	"github.com/netcloth/netcloth-chain/modules/vm/common/math"

	sdk "github.com/netcloth/netcloth-chain/types"
)

type Config struct {
	Debug                   bool   // Enables debugging
	Tracer                  Tracer // Opcode logger
	NoRecursion             bool   // Disables call, callcode, delegate call and create
	EnablePreimageRecording bool   // Enables recording of SHA3/keccak preimages

	JumpTable      [256]operation // EVM instruction table, automatically populated if unset
	ConstGasConfig *[256]uint64

	EWASMInterpreter string // External EWASM interpreter options
	EVMInterpreter   string // External EVM interpreter options
}

type hashState interface {
	hash.Hash
	Read([]byte) (int, error)
}

// Interpreter is used to run Ethereum based contracts and will utilise the
// passed environment to query external sources for state information.
// The Interpreter will run the byte code VM based on the passed
// configuration.
type Interpreter interface {
	// Run loops and evaluates the contract's code with the given input data and returns
	// the return byte-slice and an error if one occurred.
	Run(contract *Contract, input []byte, static bool) ([]byte, sdk.Error)
	// CanRun tells if the contract, passed as an argument, can be
	// run by the current interpreter. This is meant so that the
	// caller can do something like:
	//
	// ```golang
	// for _, interpreter := range interpreters {
	//   if interpreter.CanRun(contract.code) {
	//     interpreter.Run(contract.code, input)
	//   }
	// }
	// ```
	CanRun([]byte) bool
}

// EVMInterpreter represents an EVM interpreter
type EVMInterpreter struct {
	evm *EVM
	cfg Config

	intPool *intPool

	hasher     hashState
	hasherBuf  sdk.Hash
	readOnly   bool   // whether to throw on stateful modifications
	returnData []byte // last CALL's return data for subsequent reuse
}

// Run loops and evaluates the contract's code with the given input data and returns
// the return byte-slice and an error if one occurred.
//
// It's important to note that any errors returned by the interpreter should be
// considered a revert-and-consume-all-gas operation except for
// errExecutionReverted which means revert-and-keep-gas-left.
func (in *EVMInterpreter) Run(contract *Contract, input []byte, readOnly bool) (ret []byte, err sdk.Error) {
	if in.intPool == nil {
		in.intPool = poolOfIntPools.get()
		defer func() {
			poolOfIntPools.put(in.intPool)
			in.intPool = nil
		}()
	}

	// Increment the call depth which is restricted to 1024
	in.evm.depth++
	defer func() { in.evm.depth-- }()

	// Make sure the readOnly is only set if we aren't in readOnly yet.
	// This makes also sure that the readOnly flag isn't removed for child calls.
	if readOnly && !in.readOnly {
		in.readOnly = true
		defer func() { in.readOnly = false }()
	}

	// Reset the previous call's return data. It's unimportant to preserve the old buffer
	// as every returning call will return new data anyway.
	in.returnData = nil

	// Don't bother with the execution if there's no code.
	if len(contract.Code) == 0 {
		return nil, nil
	}

	var (
		op    OpCode        // current opcode
		mem   = NewMemory() // bound memory
		stack = newstack()  // local stack
		// For optimisation reason we're using uint64 as the program counter.
		// It's theoretically possible to go above 2^64. The YP defines the PC
		// to be uint256. Practically much less so feasible.
		pc   = uint64(0) // program counter
		cost uint64
		// copies used by tracer
		pcCopy  uint64 // needed for the deferred Tracer
		gasCopy uint64 // for Tracer to log gas remaining before execution
		logged  bool   // deferred Tracer should ignore already logged steps
		res     []byte // result of the opcode execution function
	)
	contract.Input = input

	// Reclaim the stack as an int pool when the execution stops
	defer func() { in.intPool.put(stack.data...) }()

	if in.cfg.Debug {
		defer func() {
			if err != nil {
				if !logged {
					in.cfg.Tracer.CaptureState(in.evm, pcCopy, op, gasCopy, cost, mem, stack, contract, in.evm.depth, err)
				} else {
					in.cfg.Tracer.CaptureFault(in.evm, pcCopy, op, gasCopy, cost, mem, stack, contract, in.evm.depth, err)
				}
			}
		}()
	}
	// The Interpreter main run loop (contextual). This loop runs until either an
	// explicit STOP, RETURN or SELFDESTRUCT is executed, an error occurred during
	// the execution of one of the operations or until the done flag is set by the
	// parent context.
	for atomic.LoadInt32(&in.evm.abort) == 0 {
		if in.cfg.Debug {
			// Capture pre-execution values for tracing.
			logged, pcCopy, gasCopy = false, pc, contract.Gas
		}

		// Get the operation from the jump table and validate the stack to ensure there are
		// enough stack items available to perform the operation.
		op = contract.GetOp(pc)
		fmt.Fprintf(os.Stderr, fmt.Sprintf("op code = %d, op = %s\n", op, op))
		operation := in.cfg.JumpTable[op]
		operation.constantGas = in.cfg.ConstGasConfig[op]
		if !operation.valid {
			//return nil, fmt.Errorf("invalid opcode 0x%x", int(op))
			return nil, sdk.ErrInternal(fmt.Sprintf("invalid opcode 0x%x", int(op)))
		}
		// Validate stack
		if sLen := stack.len(); sLen < operation.minStack {
			//return nil, fmt.Errorf("stack underflow (%d <=> %d)", sLen, operation.minStack)
			return nil, sdk.ErrInternal(fmt.Sprintf("stack underflow (%d <=> %d)", sLen, operation.minStack))
		} else if sLen > operation.maxStack {
			//return nil, fmt.Errorf("stack limit reached %d (%d)", sLen, operation.maxStack)
			return nil, sdk.ErrInternal(fmt.Sprintf("stack limit reached %d (%d)", sLen, operation.maxStack))
		}
		// If the operation is valid, enforce and write restrictions
		if in.readOnly {
			// If the interpreter is operating in readonly mode, make sure no
			// state-modifying operation is performed. The 3rd stack item
			// for a call operation is the value. Transferring value from one
			// account to the others means the state is modified and should also
			// return with an error.
			if operation.writes || (op == CALL && stack.Back(2).Sign() != 0) {
				//return nil, errWriteProtection
				//return nil, fmt.Errorf("errWriteProtection")
				return nil, sdk.ErrInternal("errWriteProtection")
			}
		}
		// Static portion of gas
		cost = operation.constantGas // For tracing
		if !contract.UseGas(operation.constantGas) {
			//return nil, fmt.Errorf("ErrOutOfGas")

			return nil, sdk.ErrInternal("ErrOutOfGas")
		}

		var memorySize uint64
		// calculate the new memory size and expand the memory to fit
		// the operation
		// Memory check needs to be done prior to evaluating the dynamic gas portion,
		// to detect calculation overflows
		if operation.memorySize != nil {
			memSize, overflow := operation.memorySize(stack)
			if overflow {
				//return nil, errGasUintOverflow
				return nil, sdk.ErrInternal("errGasUintOverflow")
			}
			// memory is expanded in words of 32 bytes. Gas
			// is also calculated in words.
			if memorySize, overflow = math.SafeMul(toWordSize(memSize), 32); overflow {
				//return nil, errGasUintOverflow
				//return nil, fmt.Errorf("errGasUintOverflow")
				return nil, sdk.ErrInternal("errGasUintOverflow")
			}
		}
		// Dynamic portion of gas
		// consume the gas and return an error if not enough gas is available.
		// cost is explicitly set so that the capture state defer method can get the proper cost
		if operation.dynamicGas != nil {
			var dynamicCost uint64
			dynamicCost, err = operation.dynamicGas(in.evm, contract, stack, mem, memorySize)
			cost += dynamicCost // total cost, for debug tracing
			if err != nil || !contract.UseGas(dynamicCost) {
				//return nil, ErrOutOfGas
				//return nil, fmt.Errorf("ErrOutOfGas")
				return nil, sdk.ErrInternal("ErrOutOfGas")

			}
		}
		if memorySize > 0 {
			mem.Resize(memorySize)
		}

		if in.cfg.Debug {
			in.cfg.Tracer.CaptureState(in.evm, pc, op, gasCopy, cost, mem, stack, contract, in.evm.depth, err)
			logged = true
		}

		stack.Print()
		// execute the operation
		res, err = operation.execute(&pc, in, contract, mem, stack)
		// verifyPool is a build flag. Pool verification makes sure the integrity
		// of the integer pool by comparing values to a default value.
		if verifyPool {
			verifyIntegerPool(in.intPool)
		}
		// if the operation clears the return data (e.g. it has returning data)
		// set the last return to the result of the operation.
		if operation.returns {
			in.returnData = res
		}

		switch {
		case err != nil:
			return nil, err
		case operation.reverts:
			return res, ErrExecutionReverted()
		case operation.halts:
			return res, nil
		case !operation.jumps:
			pc++
		}
	}
	return nil, nil
}

// CanRun tells if the contract, passed as an argument, can be
// run by the current interpreter.
func NewEVMInterpreter(evm *EVM, cfg Config) *EVMInterpreter {
	if !cfg.JumpTable[STOP].valid {
		jt := istanbulInstructionSet
		cfg.JumpTable = jt
	}

	return &EVMInterpreter{
		evm: evm,
		cfg: cfg,
	}
}

func (in *EVMInterpreter) CanRun(code []byte) bool {
	return true
}
