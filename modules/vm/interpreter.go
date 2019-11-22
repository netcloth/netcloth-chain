package vm

type Config struct {
	Debug bool // Enables debugging
	//Tracer                  Tracer // Opcode logger
	NoRecursion             bool // Disables call, callcode, delegate call and create
	EnablePreimageRecording bool // Enables recording of SHA3/keccak preimages

	//JumpTable [256]operation // EVM instruction table, automatically populated if unset

	EWASMInterpreter string // External EWASM interpreter options
	EVMInterpreter   string // External EVM interpreter options
}

// EVMInterpreter represents an EVM interpreter
type EVMInterpreter struct {
	intPool *intPool

	readOnly bool // whether to throw on stateful modifications
	return Data []byte // last CALL's return data for subsequent reuse
}
