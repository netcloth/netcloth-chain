package vm

import (
	"math/big"
	"time"

	sdk "github.com/netcloth/netcloth-chain/types"
)

// Tracer is used to collect execution traces from an EVM transaction
// execution. CaptureState is called for each step of the VM with the
// current VM state.
// Note that reference types are actual VM data structures; make copies
// if you need to retain them beyond the current call.
type Tracer interface {
	CaptureStart(from sdk.AccAddress, to sdk.AccAddress, call bool, input []byte, gas uint64, value *big.Int) error
	CaptureState(env *VM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, contract *Contract, depth int, err error) error
	CaptureFault(env *VM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, contract *Contract, depth int, err error) error
	CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) error
}
