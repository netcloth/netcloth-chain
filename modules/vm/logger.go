package vm

import (
	"math/big"

	sdk "github.com/netcloth/netcloth-chain/types"
)

// Storage represents a contract's storage
type Storage map[sdk.Hash]sdk.Hash

// Copy duplicates the current storage
func (s Storage) Copy() Storage {
	cpy := make(Storage)
	for key, value := range s {
		cpy[key] = value
	}

	return cpy
}

// LogConfig are the configuration options for structured logger the VM
type LogConfig struct {
	DisableMemory  bool // disable memory capture
	DisableStack   bool // disable stack capture
	DisableStorage bool // disable storage capture
	Debug          bool // print output during capture end
	Limit          int  // maximum length of output, but zero means unlimited
}

// StructLog is emitted to the VM each cycle and lists information about the current internal state
// prior to the execution of the statement.
type StructLog struct {
	Pc            uint64                `json:"pc"`
	Op            OpCode                `json:"op"`
	Gas           uint64                `json:"gas"`
	GasCost       uint64                `json:"gasCost"`
	Memory        []byte                `json:"memory"`
	MemorySize    int                   `json:"memSize"`
	Stack         []*big.Int            `json:"stack"`
	Storage       map[sdk.Hash]sdk.Hash `json:"-"`
	Depth         int                   `json:"depth"`
	RefundCounter uint64                `json:"refund"`
	Err           error                 `json:"-"`
}

// OpName formats the operand name in a human-readable format
func (s *StructLog) OpName() string {
	return s.Op.String()
}

// ErrorString formats the log's error as a string
func (s *StructLog) ErrorString() string {
	if s.Err != nil {
		return s.Err.Error()
	}

	return ""
}

// Tracer is used to collect execution traces from an VM transaction
// execution. CaptureState is called for each step of the VM with the
// current VM state.
// Note that reference types are actual VM data structures; make copies
// if you need to retain them beyond the current call.
type Tracer interface {
	//CaptureStart(from sdk.AccAddress, to sdk.AccAddress, call bool, input []byte, gas uint64, value *big.Int) error
	//CaptureState(env *VM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, contract *Contract, depth int, err error) error
	//CaptureFault(env *VM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, contract *Contract, depth int, err error) error
	//CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) error
}

// StructLogger is an VM state logger and implements Tracer.
//
// StructLogger can capture state based on the given Log configuration and also keeps
// a track record of modified storage which is used in reporting snapshots of the
// contract their storage.
type StructLogger struct {
	cfg LogConfig

	logs          []StructLog
	changedValues map[sdk.Address]Storage
	output        []byte
	err           error
}

// NewStructLogger returns a new logger
func NewStructLogger(cfg *LogConfig) *StructLogger {
	logger := &StructLogger{
		changedValues: make(map[sdk.Address]Storage),
	}

	if cfg != nil {
		logger.cfg = *cfg
	}

	return logger
}
