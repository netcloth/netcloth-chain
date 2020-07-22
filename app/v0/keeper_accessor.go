package v0

import (
	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"github.com/netcloth/netcloth-chain/app/v0/gov"
	"github.com/netcloth/netcloth-chain/app/v0/staking"
	"github.com/netcloth/netcloth-chain/app/v0/supply"
)

// GovKeeper return govKeeper
func (p *ProtocolV0) GovKeeper() gov.Keeper {
	return p.govKeeper
}

// SetGovKeeper set govKeeper
func (p *ProtocolV0) SetGovKeeper(gk gov.Keeper) {
	p.govKeeper = gk
}

// StakingKeeper return stakingKeeper
func (p *ProtocolV0) StakingKeeper() staking.Keeper {
	return p.stakingKeeper
}

// AccountKeeper return accountKeeper
func (p *ProtocolV0) AccountKeeper() auth.AccountKeeper {
	return p.accountKeeper
}

// SupplyKeeper return supplyKeeper
func (p *ProtocolV0) SupplyKeeper() supply.Keeper {
	return p.supplyKeeper
}
