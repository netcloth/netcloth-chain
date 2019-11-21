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
