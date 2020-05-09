package types

type GenesisState struct {
	Params    Params    `json:"params" yaml:"params"`
	IPALNodes IPALNodes `json:"ipal_nodes" yaml:"ipal_nodes"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

func NewGenesisState(params Params, ipalNodes IPALNodes) GenesisState {
	return GenesisState{
		Params:    params,
		IPALNodes: ipalNodes,
	}
}
