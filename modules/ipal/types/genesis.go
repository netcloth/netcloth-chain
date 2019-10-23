package types

type GenesisState struct {
	Params Params `json:"params" yaml:"params"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState {
		Params:DefaultParams(),
	}
}

func NewGenesisState(params Params) GenesisState {
	return GenesisState{Params:params}
}

