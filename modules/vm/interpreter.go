package vm

import (
	"hash"

	"github.com/netcloth/netcloth-chain/modules/vm/common"
)

type Config struct {
	Debug bool // Enables debugging
	//Tracer                  Tracer // Opcode logger
	NoRecursion             bool // Disables call, callcode, delegate call and create
	EnablePreimageRecording bool // Enables recording of SHA3/keccak preimages

	//JumpTable [256]operation // EVM instruction table, automatically populated if unset

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
	Run(contract *Contract, input []byte, static bool) ([]byte, error)
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
	evm     *EVM
	intPool *intPool

	hasher     hashState
	hasherBuf  common.Hash
	readOnly   bool   // whether to throw on stateful modifications
	returnData []byte // last CALL's return data for subsequent reuse
}
