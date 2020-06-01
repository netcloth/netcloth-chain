package types

type GenesisState struct {
	Params Params `json:"params" yaml:"params"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

func NewGenesisState(params Params) GenesisState {
	return GenesisState{Params: params}
}

func ValidateGenesis(data GenesisState) error {
	maxCodeSize := data.Params.MaxCodeSize
	if err := validateMaxCodeSize(maxCodeSize); err != nil {
		return err
	}

	vmOpGasParams := data.Params.VMOpGasParams
	if err := validateVMOpGasParams(vmOpGasParams); err != nil {
		return err
	}

	vmCommonGasParams := data.Params.VMCommonGasParams
	return validateVMCommonGasParams(vmCommonGasParams)
}
