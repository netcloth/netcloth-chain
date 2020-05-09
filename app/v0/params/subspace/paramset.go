package subspace

type (
	ValueValidatorFn func(value interface{}) error

	// Used for associating paramsubspace key and field of param structs
	ParamSetPair struct {
		Key         []byte
		Value       interface{}
		ValidatorFn ValueValidatorFn
	}
)

func NewParamSetPair(key []byte, value interface{}, vfn ValueValidatorFn) ParamSetPair {
	return ParamSetPair{key, value, vfn}
}

// Slice of KeyFieldPair
type ParamSetPairs []ParamSetPair

// ParamSet defines an interface for structs containing parameters for a module
type ParamSet interface {
	ParamSetPairs() ParamSetPairs
}
