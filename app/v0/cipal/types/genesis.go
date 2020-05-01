package types

// GenesisState is the supply state that must be provided at genesis.
type GenesisState struct {
	CIPALObjs CIPALObjects `json:"cipal_objects" yaml:"cipal_objects"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(CIPALObjs CIPALObjects) GenesisState {
	return GenesisState{
		CIPALObjs: CIPALObjs,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}
