package v0

import "github.com/netcloth/netcloth-chain/app/v0/gov"

func (p *ProtocolV0) GovKeeper() gov.Keeper {
	return p.govKeeper
}

func (p *ProtocolV0) SetGovKeeper(gk gov.Keeper) {
	p.govKeeper = gk
}
