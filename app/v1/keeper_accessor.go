package v1

import "github.com/netcloth/netcloth-chain/app/v0/gov"

func (p *ProtocolV1) GovKeeper() gov.Keeper {
	return p.govKeeper
}

func (p *ProtocolV1) SetGovKeeper(gk gov.Keeper) {
	p.govKeeper = gk
}
