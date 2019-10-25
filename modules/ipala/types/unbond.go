package types

import (
    sdk "github.com/NetCloth/netcloth-chain/types"
    "time"
)

type UnBondings []UnBonding

type UnBonding struct {
    AccountAddress sdk.AccAddress `json:"account_address" yaml:"account_address"`
    Amount sdk.Coin `json:"amount" yaml:"amount"`
    EndTime time.Time `json:"end_time" yaml:"end_time"`
}

func NewUnBonding(aa sdk.AccAddress, amt sdk.Coin, endTime time.Time) UnBonding {
    return UnBonding {
        AccountAddress: aa,
        Amount:         amt,
        EndTime:        endTime,
    }
}

func (ub UnBonding) IsMature(now time.Time) bool {
    return !ub.EndTime.After(now)
}

