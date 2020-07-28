package types

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

// DONTCOVER

// QuerySigningInfoParams defines the params for the following queries:
// - 'custom/slashing/signingInfo'
type QuerySigningInfoParams struct {
	ConsAddress sdk.ConsAddress
}

// NewQuerySigningInfoParams creates a new QuerySigningInfoParams instance
func NewQuerySigningInfoParams(consAddr sdk.ConsAddress) QuerySigningInfoParams {
	return QuerySigningInfoParams{consAddr}
}

// QuerySigningInfosParams defines the params for the following queries:
// - 'custom/slashing/signingInfos'
type QuerySigningInfosParams struct {
	Page, Limit int
}

// NewQuerySigningInfosParams creates a new QuerySigningInfosParams instance
func NewQuerySigningInfosParams(page, limit int) QuerySigningInfosParams {
	return QuerySigningInfosParams{page, limit}
}
