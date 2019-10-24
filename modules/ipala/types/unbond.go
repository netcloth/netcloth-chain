package types

import (
    "fmt"
    "github.com/NetCloth/netcloth-chain/codec"
    sdk "github.com/NetCloth/netcloth-chain/types"
    "time"
)

type UnBond struct {
    EndTime     time.Time   `json:"end_time" yaml:"end_time"`
    Amount      sdk.Coin    `json:"amount" yaml:"amount"`
}

func NewUnBond(endTime time.Time, amt sdk.Coin) UnBond {
    return UnBond {
        EndTime: endTime,
        Amount:  amt,
    }
}

func (ub UnBond) IsMature(now time.Time) bool {
    return !ub.EndTime.After(now)
}

type UnBonds struct {
    AccountAddress      sdk.AccAddress `json:"account_address" yaml:"account_address"`
    Entries             []UnBond       `json:"entries" yaml:"entries"`
}

func NewUnBonds(aa sdk.AccAddress, endTime time.Time, amt sdk.Coin) UnBonds  {
    entry := NewUnBond(endTime, amt)
    return UnBonds {
        AccountAddress: aa,
        Entries:        []UnBond{entry},
    }
}

func (p *UnBonds) AddEntry(endTime time.Time, amt sdk.Coin) {
    entry := NewUnBond(endTime, amt)
    p.Entries = append(p.Entries, entry)
}

func (p *UnBonds) RemoveEntry(i int64) {
    p.Entries = append(p.Entries[:i], p.Entries[i+1:]...)
}

func MustMarshalUnstaking(cdc *codec.Codec, unBonds UnBonds) []byte {
    return cdc.MustMarshalBinaryLengthPrefixed(unBonds)
}

func MustUnmarshalUnstaking(cdc *codec.Codec, value []byte) UnBonds {
    unBonds, err := UnmarshalUnstaking(cdc, value)
    if err != nil {
        panic(err)
    }
    return unBonds
}

func UnmarshalUnstaking(cdc *codec.Codec, value []byte) (unBonds UnBonds, err error) {
    err = cdc.UnmarshalBinaryLengthPrefixed(value, &unBonds)
    return unBonds, err
}

func (ubs UnBonds) String() string {
    out := fmt.Sprintf(`
UnBonds:
  AccountAddress: %s
  Entries:
`,
ubs.AccountAddress)

    for i, entry := range ubs.Entries {
        out += fmt.Sprintf(`    UnBond %d:
      EndTime: %v
      Amount: %s
`,
i, entry.EndTime.String(), entry.Amount)
    }

    return out
}

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

