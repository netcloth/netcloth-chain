package types

// GenesisState - all slashing state that must be provided at genesis
type GenesisState struct {
	Params       Params                          `json:"params" yaml:"params"`
	SigningInfos map[string]ValidatorSigningInfo `json:"signing_infos" yaml:"signing_infos"`
	MissedBlocks map[string][]MissedBlock        `json:"missed_blocks" yaml:"missed_blocks"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params, signingInfos map[string]ValidatorSigningInfo, missedBlocks map[string][]MissedBlock,
) GenesisState {

	return GenesisState{
		Params:       params,
		SigningInfos: signingInfos,
		MissedBlocks: missedBlocks,
	}
}

// MissedBlock
type MissedBlock struct {
	Index  int64 `json:"index" yaml:"index"`
	Missed bool  `json:"missed" yaml:"missed"`
}

// DefaultGenesisState - default GenesisState
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:       DefaultParams(),
		SigningInfos: make(map[string]ValidatorSigningInfo),
		MissedBlocks: make(map[string][]MissedBlock),
	}
}

// ValidateGenesis validates the slashing genesis parameters
func ValidateGenesis(data GenesisState) error {
	err := data.Params.Validate()
	if err != nil {
		return err
	}

	return nil
}
