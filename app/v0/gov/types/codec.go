package types

import (
	"github.com/netcloth/netcloth-chain/codec"
)

// ModuleCdc - generic sealed codec to be used throughout this module
var ModuleCdc = codec.New()

// RegisterCodec registers all the necessary types and interfaces for
// governance.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Content)(nil), nil)

	cdc.RegisterConcrete(MsgSubmitProposal{}, "nch/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "nch/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgVote{}, "nch/MsgVote", nil)

	cdc.RegisterConcrete(TextProposal{}, "nch/TextProposal", nil)
	cdc.RegisterConcrete(SoftwareUpgradeProposal{}, "nch/SoftwareUpgradeProposal", nil)
}

// RegisterProposalTypeCodec registers an external proposal content type defined
// in another module for the internal ModuleCdc. This allows the MsgSubmitProposal
// to be correctly Amino encoded and decoded.
func RegisterProposalTypeCodec(o interface{}, name string) {
	ModuleCdc.RegisterConcrete(o, name, nil)
}

// TODO determine a good place to seal this codec
func init() {
	RegisterCodec(ModuleCdc)
}
