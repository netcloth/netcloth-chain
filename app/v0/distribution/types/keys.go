package types

import (
	"github.com/netcloth/netcloth-chain/baseapp/protocol"
)

const (
	ModuleName = protocol.DistributionModuleName

	// StoreKey is the store key string for distribution
	StoreKey = ModuleName

	// RouterKey is the message route for distribution
	RouterKey = ModuleName

	// QuerierRoute is the querier route for distribution
	QuerierRoute = ModuleName
)
