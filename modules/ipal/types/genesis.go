package types

type GenesisState struct {
	Params       Params       `json:"params" yaml:"params"`
	ServiceNodes ServiceNodes `json:"service_nodes" yaml:"service_nodes"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

func NewGenesisState(params Params, serviceNodes ServiceNodes) GenesisState {
	return GenesisState{
		Params:       params,
		ServiceNodes: serviceNodes,
	}
}
