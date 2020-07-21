package vm

import (
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/netcloth/netcloth-chain/app/v0/vm/common/math"
	"github.com/netcloth/netcloth-chain/app/v0/vm/types"
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

// LogConfig are the configuration options for structured logger the EVM
type LogConfig struct {
	DisableMemory  bool // disable memory capture
	DisableStack   bool // disable stack capture
	DisableStorage bool // disable storage capture
	Debug          bool // print output during capture end
	Limit          int  // maximum length of output, but zero means unlimited
}

// StructLog is emitted to the EVM each cycle and lists information about the current internal state
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
	Depth         uint64                `json:"depth"`
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
	CaptureStart(from sdk.AccAddress, to sdk.AccAddress, call bool, input []byte, gas uint64, value *big.Int) error
	CaptureState(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, contract *Contract, depth uint64, err error) error
	CaptureFault(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, contract *Contract, depth uint64, err error) error
	CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) error
}

// StructLogger is an VM state logger and implements Tracer.
//
// StructLogger can capture state based on the given Log configuration and also keeps
// a track record of modified storage which is used in reporting snapshots of the
// contract their storage.
type StructLogger struct {
	cfg LogConfig

	logs          []StructLog
	changedValues map[string]Storage
	output        []byte
	err           error
}

// NewStructLogger returns a new logger
func NewStructLogger(cfg *LogConfig) *StructLogger {
	logger := &StructLogger{
		changedValues: make(map[string]Storage),
	}

	if cfg != nil {
		logger.cfg = *cfg
	}

	return logger
}

// CaptureStart implements the Tracer interface to initialize the tracing operation.
func (l *StructLogger) CaptureStart(from sdk.AccAddress, to sdk.AccAddress, create bool, input []byte, gas uint64, value *big.Int) error {
	return nil
}

// CaptureState logs a new structured log message and pushes it out to the environment
//
// CaptureState also tracks SSTORE ops to track dirty values.
func (l *StructLogger) CaptureState(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, contract *Contract, depth uint64, err error) error {
	// check if already accumulated the specified number of logs
	if l.cfg.Limit != 0 && l.cfg.Limit <= len(l.logs) {
		return ErrTraceLimitReached
	}

	// initialise new changed values storage container for this contract if not presend
	fmt.Println(contract.Address().String())
	if l.changedValues[contract.Address().String()] == nil {
		l.changedValues[contract.Address().String()] = make(Storage)
	}

	// capture SSTORE opcodes and determine the changed value and store it in the local storage container
	if op == SSTORE && stack.len() >= 2 {
		var (
			value   = sdk.BigToHash(stack.data[stack.len()-2])
			address = sdk.BigToHash(stack.data[stack.len()-1])
		)

		l.changedValues[contract.Address().String()][address] = value
	}

	// copy a snapshot of the current memory state to a new buffer
	var mem []byte
	if !l.cfg.DisableMemory {
		mem = make([]byte, len(memory.Data()))
		copy(mem, memory.Data())
	}

	// copy a snapshot of the current stack state to a new buffer
	var stck []*big.Int
	if !l.cfg.DisableStack {
		stck = make([]*big.Int, len(stack.Data()))
		for i, item := range stack.Data() {
			stck[i] = new(big.Int).Set(item)
		}
	}

	// copy a snapshot of the current storage to a new buffer
	var storage Storage
	if !l.cfg.DisableStorage {
		storage = l.changedValues[contract.Address().String()].Copy()
	}

	// create a new snapshot of the EVM
	log := StructLog{pc, op, gas, cost, mem, memory.Len(), stck, storage, depth, env.StateDB.GetRefund(), err}

	l.logs = append(l.logs, log)
	return nil
}

// CaptureFault implements the Tracer interface to trace an execution fault
// while running an opcode.
func (l *StructLogger) CaptureFault(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, contract *Contract, depth uint64, err error) error {
	return nil
}

// CaptureEnd is called after the call finishes to finalize the tracing
func (l *StructLogger) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) error {
	l.output = output
	l.err = err
	if l.cfg.Debug {
		fmt.Printf("0x%x\n", output)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}
	return nil
}

// StructLogs returns the captured log entries
func (l *StructLogger) StructLogs() []StructLog {
	return l.logs
}

// Error returns the VM error captured by the trace
func (l *StructLogger) Error() error {
	return l.err
}

// Output returns the VM return value captured by the trace
func (l *StructLogger) Output() []byte {
	return l.output
}

// WriteTrace writes a formatted trace to the given writer
func WriteTrace(writer io.Writer, logs []StructLog) {
	for _, log := range logs {
		fmt.Fprintf(writer, "%-16spc=%08d gas=%v cost=%v", log.Op, log.Pc, log.Gas, log.GasCost)
		if log.Err != nil {
			fmt.Fprintf(writer, " ERROR: %v", log.Err)
		}
		fmt.Fprintln(writer)

		if len(log.Stack) > 0 {
			fmt.Fprintln(writer, "Stack:")
			for i := len(log.Stack) - 1; i >= 0; i-- {
				fmt.Fprintf(writer, "%08d  %x\n", len(log.Stack)-i-1, math.PaddedBigBytes(log.Stack[i], 32))
			}
		}
		if len(log.Memory) > 0 {
			fmt.Fprintln(writer, "Memory:")
			fmt.Fprint(writer, hex.Dump(log.Memory))
		}
		if len(log.Storage) > 0 {
			fmt.Fprintln(writer, "Storage:")
			for h, item := range log.Storage {
				fmt.Fprintf(writer, "%x: %x\n", h, item)
			}
		}
		fmt.Fprintln(writer)
	}
}

// WriteLogs writes vm logs in a readable format to the given writer
func WriteLogs(writer io.Writer, logs []*types.Log) {
	for _, log := range logs {
		fmt.Fprintf(writer, "LOG%d: %x bn=%d txi=%x\n", len(log.Topics), log.Address, log.BlockNumber, log.TxIndex)

		for i, topic := range log.Topics {
			fmt.Fprintf(writer, "%08d  %x\n", i, topic)
		}

		fmt.Fprint(writer, hex.Dump(log.Data))
		fmt.Fprintln(writer)
	}
}
