// nolint

package v0

import (
	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"github.com/netcloth/netcloth-chain/app/v0/gov"
	"github.com/netcloth/netcloth-chain/app/v0/staking"
	"github.com/netcloth/netcloth-chain/app/v0/supply"
)

func (p *ProtocolV0) GovKeeper() gov.Keeper {
	return p.govKeeper
}

func (p *ProtocolV0) SetGovKeeper(gk gov.Keeper) {
	p.govKeeper = gk
}

func (p *ProtocolV0) StakingKeeper() staking.Keeper {
	return p.stakingKeeper
}

func (p *ProtocolV0) AccountKeeper() auth.AccountKeeper {
	return p.accountKeeper
}

func (p *ProtocolV0) SupplyKeeper() supply.Keeper {
	return p.supplyKeeper
}
