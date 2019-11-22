package vm

import (
	sdk "github.com/netcloth/netcloth-chain/types"
	"hash"
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

// EVMInterpreter represents an EVM interpreter
type EVMInterpreter struct {
	evm     *EVM
	intPool *intPool

	hasher     hashState
	hasherBuf  sdk.Hash
	readOnly   bool   // whether to throw on stateful modifications
	returnData []byte // last CALL's return data for subsequent reuse
}
