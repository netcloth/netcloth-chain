package types

import (
	"github.com/netcloth/netcloth-chain/baseapp/protocol"
)

const (
	// ModuleName is the module name constant used in many places
	ModuleName = protocol.SupplyStoreKey

	// StoreKey is the store key string for supply
	StoreKey = ModuleName

	// RouterKey is the message route for supply
	RouterKey = ModuleName

	// QuerierRoute is the querier route for supply
	QuerierRoute = ModuleName
)
